package main

import (
	"DBMS/controllers"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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
var db *sql.DB

func main() {
	e = echo.New()
	e.POST("/login", controllers.Login) //user login

	e.Logger.Fatal(e.Start(":1323"))
}
