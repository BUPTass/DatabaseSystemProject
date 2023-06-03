package controllers

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserName  string `db:"username"`
	Password  string `db:"password"`
	Level     int    `db:"level"`     //0 for admin,1 for normal user
	Conformed int    `db:"conformed"` //false for unavailable
}

const (
	username_          string = "root"
	password_          string = "1594568520h"
	userDbName_        string = "userinfo"
	userinfoTableName_ string = "info"
	ip_                string = "127.0.0.1"
	port_              int    = 3306
)

var user_info *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username_, password_, ip_, port_, userDbName_)
	user_info, _ = sqlx.Open("mysql", dsn)
	if err := user_info.Ping(); err != nil {
		panic(err)
	}
	database_init()
}
func get_hash(text string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return hash
}
func check_password(username, password string) bool {
	var ans []UserInfo
	err := user_info.Select(&ans, "select * from info where username like ? and conformed = 1", username)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if len(ans) != 1 {
		return false
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(ans[0].Password), []byte(password))
		if err != nil {
			log.Println(err.Error())
			return false
		}
		return true
	}
}
func get_level(username, password string) int {
	var ans []UserInfo
	err := user_info.Select(&ans, "select * from info where username like ? and conformed = 1", username)
	if err != nil && len(ans) != 1 {
		return -1
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(ans[0].Password), []byte(password))
		if err == nil {
			return ans[0].Level
		}
		return -1
	}
}
func Login(c echo.Context) error {
	//"/login?username=uname&password=pwd"
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	if check_password(username, password) {
		//create session

		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 7,
		}
		//record session data
		sess.Values["id"] = username
		sess.Values["level"] = get_level(username, password)

		//saving data
		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			log.Println(err.Error())
			return err
		}
		return c.String(http.StatusOK, fmt.Sprintf("login success %s %d", username, sess.Values["level"]))
	} else {
		return c.String(http.StatusForbidden, "please check your username and password")
	}
}
func GetUsers(c echo.Context) error {
	//"/show/users"
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusMethodNotAllowed, "insufficient permissions")
	}
	//select users
	var ans []UserInfo
	err = user_info.Select(&ans, fmt.Sprintf("select * from %s", userinfoTableName_))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ans)
}
func AddUser(c echo.Context) error {
	//"/add/user?username=name"
	uname := c.QueryParam("username")

	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusMethodNotAllowed, "insufficient permissions")
	}
	//change user conformed to 1
	stmt, _ := user_info.Prepare("update ? set conformed = 1 where username like ?")
	_, err = stmt.Exec(userinfoTableName_, uname)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
func DeleteUser(c echo.Context) error {
	//"/delete/user?username=name"
	uname := c.QueryParam("username")
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusMethodNotAllowed, "insufficient permissions")
	}
	//delete normal user:uname
	stmt, _ := user_info.Prepare("delete from info where username like ?")
	_, err = stmt.Exec(uname)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
func Regist(c echo.Context) error {
	//"/regist?username=name&password=password&level=num"
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	level := c.QueryParam("level")
	var ans []UserInfo
	err := user_info.Select(&ans, "select * from info where username like ?", username)
	if err != nil {
		log.Println(err.Error())
		return err
	} else if len(ans) != 0 {
		return c.String(http.StatusOK, "username exist")
	}
	hashText := get_hash(password)
	stmt, _ := user_info.Prepare("insert into info values(?,?,?,0)")
	_, err = stmt.Exec(username, string(hashText), level)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return c.String(http.StatusOK, "success")
}
func Logout(c echo.Context) error {
	//"/logout"
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
	sess.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, "logout success")
}
