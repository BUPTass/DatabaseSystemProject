package main

import (
	"DBMS/controllers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const ip string = "127.0.0.1"
const port int = 3306

var dbinfo struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
	DbName   string `json:"dbName" form:"dbName" query:"dbName"`
}
var e *echo.Echo
var db *sqlx.DB
var sessionKey string

func main() {
	e = echo.New()
	//session init
	sessionPath := "./statics/session_data"
	sessionKey = "anything"
	e.Use(session.Middleware(sessions.NewFilesystemStore(sessionPath, []byte(sessionKey))))
	//route
	//user management
	e.POST("/login", controllers.Login)              //user login
	e.POST("/logout", controllers.Logout)            //user logout
	e.GET("/show/users", controllers.GetUsers)       //show all users
	e.POST("/add/user", controllers.AddUser)         //add user
	e.DELETE("/delete/user", controllers.DeleteUser) //delete user
	e.POST("/regist", controllers.Regist)            //user regist

	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "hello") })

	e.Logger.Fatal(e.Start(":1323"))
}
