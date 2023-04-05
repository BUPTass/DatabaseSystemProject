package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

var (
	username string = ""
	password string = ""
	ip       string = "127.0.0.1"
	port     int    = 3306
	dbName   string = "test"
	charSet  string = "utf-8"
)
var e *echo.Echo
var db *sql.DB

func main() {
	e = echo.New()
	e.POST("/login", login)
	e.Logger.Fatal(e.Start(":1323"))
}

func login(c echo.Context) error {
	//"/login?username=uname&password=pwd&dbName=name"
	username = c.QueryParam("username")
	password = c.QueryParam("password")
	dbName = c.QueryParam("dbName")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, ip, port, dbName)
	_db, _ := sql.Open("mysql", dsn)
	err := _db.Ping()
	if err == nil {
		db = _db
		return c.String(http.StatusOK, username+":"+dbName)
	}
	return err
}
