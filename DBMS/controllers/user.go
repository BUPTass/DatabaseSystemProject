package controllers

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type UserInfo struct {
	UserName  string `db:"username"`
	Password  string `db:"password"`
	Level     int    `db:"level"`     //0 for admin,1 for normal user
	Conformed int    `db:"conformed"` //false for unavailable
}

const (
	username          string = "root"
	password          string = "1594568520h"
	userDbName        string = "userinfo"
	userinfoTableName string = "info"
	ip                string = "127.0.0.1"
	port              int    = 3306
)

var user_info *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, ip, port, userDbName)
	user_info, _ = sqlx.Open("mysql", dsn)
	if err := user_info.Ping(); err != nil {
		panic(err)
	}
	database_init()
}
func get_hash(text, key string) []byte {
	hash := hmac.New(sha1.New, []byte(username))
	hash.Write([]byte(password))
	return hash.Sum(nil)
}
func check_password(username, password string) bool {
	hashText := get_hash(password, username)
	var ans []UserInfo
	user_info.Select(&ans, fmt.Sprintf("select * from %s where username like '%s' and password like '%x' and conformed = 1", userinfoTableName, username, hashText))
	return len(ans) == 1
}
func get_level(username, password string) int {
	hashText := get_hash(password, username)
	var ans []UserInfo
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where username like '%s' and password like '%x' and conformed = 1", userinfoTableName, username, hashText))
	if err == nil && len(ans) == 1 {
		return ans[0].Level
	}
	return -1
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
			panic(err)
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
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//select users
	var ans []UserInfo
	err = user_info.Select(&ans, fmt.Sprintf("select * from %s", userinfoTableName))
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
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//change user conformed to 1
	_, err = user_info.Exec(fmt.Sprintf("update %s set conformed = 1 where username like '%s'", userinfoTableName, uname))
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
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//delete normal user:uname
	_, err = user_info.Exec(fmt.Sprintf("delete from %s where username like '%s' and level = 1", userinfoTableName, uname))
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
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where username like '%s'", userinfoTableName, username))
	if err != nil || len(ans) != 0 {
		return c.String(http.StatusOK, "username exist")
	}
	hashText := get_hash(password, username)
	_, err = user_info.Exec(fmt.Sprintf("insert into %s values('%s','%s',%s,%s)", userinfoTableName, username, hashText, level, "0"))
	if err != nil {
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
