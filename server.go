package main

import (
	"DatabaseSystemProject/Auth"
	"DatabaseSystemProject/Import"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Database System Project API backend")
	})
	e.POST("/import/tbCell", func(c echo.Context) error {
		file, err := c.FormFile("xlsx")

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		err = Import.AddtbCell(file)
		if err != nil {
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbCell Added")
		}
	})

	e.POST("/auth/signup", Auth.RegisterHandler)
	e.POST("/auth/login", Auth.LoginHandler)
	e.GET("/auth/logout", Auth.LogoutHandler)
	e.Logger.Fatal(e.Start(":1333"))
}
