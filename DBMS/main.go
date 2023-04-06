package main

import (
	"DBMS/controllers"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var e *echo.Echo

func main() {
	e = echo.New()
	//session init
	sessionPath := "./statics/session_data"
	sessionKey := "anything"
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
