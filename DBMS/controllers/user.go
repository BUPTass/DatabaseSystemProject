package controllers

import (
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
	Username          string = "root"
	Password          string = "1594568520h"
	UserDbName        string = "userinfo"
	UserinfoTableName string = "info"
	Ip                string = "127.0.0.1"
	Port              int    = 3306
)

var user_info *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Username, Password, Ip, Port, UserDbName)
	user_info, _ = sqlx.Open("mysql", dsn)
	if err := user_info.Ping(); err != nil {
		panic(err)
	}
}
func check_password(username, password string) bool {
	var ans []UserInfo
	user_info.Select(&ans, fmt.Sprintf("select * from %s where username like '%s' and password like '%s' and conformed = 1", UserinfoTableName, username, password))
	return len(ans) == 1
}
func get_level(username, password string) int {
	var ans []UserInfo
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where username like '%s' and password like '%s' and conformed = 1", UserinfoTableName, username, password))
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
		sess, _ := session.Get(username, c)
		sess.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 7,
		}
		//record session data
		sess.Values["id"] = username
		sess.Values["level"] = get_level(username, password)
		sess.Values["isLogin"] = true

		//saving data
		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			panic(err)
		}
		return c.String(http.StatusOK, fmt.Sprintf("login success %s %d", username, sess.Values["level"]))
	} else {
		return c.String(http.StatusOK, "please check your username and password")
	}
}
func GetUsers(c echo.Context) error {
	//"/show/users?adminname=name"
	name := c.QueryParam("adminname")
	//check session
	sess, err := session.Get(name, c)
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//select users
	var ans []UserInfo
	err = user_info.Select(&ans, fmt.Sprintf("select * from %s", UserinfoTableName))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ans)
}
func AddUser(c echo.Context) error {
	//"/add/user?adminname=name&username=name"
	name := c.QueryParam("adminname")
	uname := c.QueryParam("username")

	//check session
	sess, err := session.Get(name, c)
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//change user conformed to 1
	_, err = user_info.Exec(fmt.Sprintf("update %s set conformed = 1 where username like '%s'", UserinfoTableName, uname))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
func DeleteUser(c echo.Context) error {
	//"/delete/user?adminname=name&username=name"
	name := c.QueryParam("adminname")
	uname := c.QueryParam("username")
	//check session
	sess, err := session.Get(name, c)
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//delete normal user:uname
	_, err = user_info.Exec(fmt.Sprintf("delete from %s where username like '%s' and level = 1", UserinfoTableName, uname))
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
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where username like '%s'", UserinfoTableName, username))
	if err != nil || len(ans) != 0 {
		return c.String(http.StatusOK, "username exist")
	}
	_, err = user_info.Exec(fmt.Sprintf("insert into %s values('%s','%s',%s,%s)", UserinfoTableName, username, password, level, "0"))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
func Logout(c echo.Context) error {
	//"/logout?username=name"
	username := c.QueryParam("username")
	sess, _ := session.Get(username, c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
	sess.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, "logout success")
}
