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
	UserName  string `db:"userName"`
	Password  string `db:"password"`
	Level     int    `db:"level"`     //0 for admin,1 for normal user
	Conformed int    `db:"conformed"` //false for unavailable
}

var dbinfo struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

const (
	Username string = "root"
	Password string = "1594568520h"
	DbName   string = "userinfo"
	ip       string = "127.0.0.1"
	port     int    = 3306
)

var user_info *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Username, Password, ip, port, DbName)
	user_info, _ = sqlx.Open("mysql", dsn)
	if err := user_info.Ping(); err != nil {
		panic(err)
	}
}
func check_password(username, password string) bool {
	var ans []UserInfo
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where userName like %s and password like %s ans conformed = 1", DbName, username, password))
	if err == nil && len(ans) == 0 {
		return true
	}
	return false
}
func get_level(username, password string) int {
	var ans []UserInfo
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where userName like %s and password like %s ans conformed = 1", DbName, username, password))
	if err == nil && len(ans) == 0 {
		return ans[0].Level
	}
	return -1
}
func Login(c echo.Context) error {
	//"/login?username=uname&password=pwd"
	if err := c.Bind(dbinfo); err != nil {
		return err
	}
	if check_password(dbinfo.Username, dbinfo.Password) {
		//create session
		sess, _ := session.Get(dbinfo.Username, c)
		sess.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 7,
		}
		//record session data
		sess.Values["id"] = dbinfo.Username
		sess.Values["level"] = get_level(dbinfo.Username, dbinfo.Password)
		sess.Values["isLogin"] = true

		//saving data
		sess.Save(c.Request(), c.Response())

		return c.String(http.StatusOK, "login success")
	} else {
		return c.String(http.StatusOK, "please check your username and password")
	}
}
func GetUsers(c echo.Context) error {
	//"/show/users/adminname=name"
	name := c.QueryParam("adminname")
	//check session
	sess, err := session.Get(name, c)
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//select users
	var ans []UserInfo
	err = user_info.Select(&ans, fmt.Sprintf("select * from %s", DbName))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ans)
}
func AddUser(c echo.Context) error {
	//"/add/user/adminname=name&username=name"
	name := c.QueryParam("adminname")
	uname := c.QueryParam("username")
	//check session
	sess, err := session.Get(name, c)
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//change user conformed to 1
	_, err = user_info.Exec(fmt.Sprintf("update %s set conformed = 1 where userName = %s", DbName, uname))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
func DeleteUser(c echo.Context) error {
	//"/delete/user/adminname=name&username=name"
	name := c.QueryParam("adminname")
	uname := c.QueryParam("username")
	//check session
	sess, err := session.Get(name, c)
	if err != nil || sess.Values["level"] != 0 {
		return err
	}
	//delete normal user:uname
	_, err = user_info.Exec(fmt.Sprintf("delete from %s where username = %s and level = 1", DbName, uname))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
func Regist(c echo.Context) error {
	//"/regist/username=name&password=password&level=num"
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	level := c.QueryParam("level")
	var ans []UserInfo
	err := user_info.Select(&ans, fmt.Sprintf("select * from %s where userName = %s", DbName, username))
	if err != nil {
		return err
	} else if len(ans) != 0 {
		return c.String(http.StatusOK, "username exist")
	}
	_, err = user_info.Exec(fmt.Sprintf("insert into %s values(%s,%s,%s,%s)", DbName, username, password, level, "0"))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "success")
}
