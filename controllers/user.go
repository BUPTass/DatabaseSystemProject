package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserName  string `db:"username" json:"username"`
	Password  string `db:"password" json:"password,omitempty"`
	Level     int    `db:"level" json:"level"`         //0 for admin,1 for normal user
	Confirmed bool   `db:"confirmed" json:"confirmed"` //false for unavailable
}

const (
	username_          string = "root"
	password_          string = "1taNWY1vXdTc4_-j"
	userDbName_        string = "LTE"
	ip_                string = "127.0.0.1"
	port_              int    = 3306
	userinfoTableName_ string = "info"
)

var user_info *sqlx.DB

func init() {
	// Read the database connection details from environmental variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username_, password_, ip_, port_, userDbName_)
	//user_info, _ = sqlx.Open("mysql", dsn)

	user_info, _ = sqlx.Open("mysql", dbURI)
	if err := user_info.Ping(); err != nil {
		panic(err)
	}
	database_init()
}
func get_hash(text string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil
	}
	return hash
}
func check_password(username, password string) bool {
	var ans []UserInfo
	err := user_info.Select(&ans, "select * from info where username = ? and confirmed = true", username)
	if err != nil {
		log.Println(err)
		return false
	}
	if len(ans) != 1 {
		return false
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(ans[0].Password), []byte(password))
		if err != nil {
			log.Println(username + ": login failed")
			return false
		}
		return true
	}
}
func get_level(username, password string) int {
	var ans []UserInfo
	err := user_info.Select(&ans, "select * from info where username = ? and confirmed = true", username)
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
			log.Println(err)
			return c.NoContent(http.StatusBadGateway)
		}
		return c.String(http.StatusOK, fmt.Sprintf("login success %s %d", username, sess.Values["level"]))
	} else {
		return c.String(http.StatusUnauthorized, "please check your username and password")
	}
}
func GetUsers(c echo.Context) error {
	//"/show/users"
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	//select users
	var ans []UserInfo
	err = user_info.Select(&ans, fmt.Sprintf("select username,level,confirmed from %s", userinfoTableName_))
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, ans)
}

func GetUnconfirmedUsers(c echo.Context) error {
	//"/show/users"
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	//select users
	var ans []string
	err = user_info.Select(&ans, fmt.Sprintf("select username from %s where confirmed = false", userinfoTableName_))
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.JSON(http.StatusOK, ans)
}

func AddUser(c echo.Context) error {
	//"/add/user?username=name"
	uname := c.QueryParam("username")

	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	//change user confirmed to 1
	stmt, _ := user_info.Prepare("update info set confirmed = true where username = ?")
	_, err = stmt.Exec(uname)
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.String(http.StatusOK, "success")
}
func DeleteUser(c echo.Context) error {
	//"/delete/user?username=name"
	uname := c.QueryParam("username")
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	//delete normal user:uname
	stmt, _ := user_info.Prepare("delete from info where username = ?")
	_, err = stmt.Exec(uname)
	if err != nil {
		return c.NoContent(http.StatusBadGateway)
	}
	return c.String(http.StatusOK, "success")
}
func Register(c echo.Context) error {
	//"/register?username=name&password=password&level=num"
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	level := c.QueryParam("level")
	var ans []UserInfo
	err := user_info.Select(&ans, "select * from info where username = ?", username)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadGateway)
	} else if len(ans) != 0 {
		return c.String(http.StatusConflict, "username already existed")
	}
	hashText := get_hash(password)
	stmt, _ := user_info.Prepare("insert into info values(?,?,?,0)")
	_, err = stmt.Exec(username, string(hashText), level)

	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadGateway)
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
