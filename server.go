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
		path := c.FormValue("path")

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err := Import.AddtbCell(path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbCell Added")
		}
	})
	e.POST("/import/tbKPI", func(c echo.Context) error {
		path := c.FormValue("path")

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err := Import.AddtbKPI(path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbKPI Added")
		}
	})
	e.POST("/import/tbRPB", func(c echo.Context) error {
		path := c.FormValue("path")

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err := Import.AddtbRPB(path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbRPB Added")
		}
	})
	e.POST("/import/tbMROData", func(c echo.Context) error {
		path := c.FormValue("path")

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err := Import.AddtbMROData(path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbMROData Added")
		}
	})

	e.POST("/upload", func(c echo.Context) error {
		file, err := c.FormFile("file")

		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		filename, err := Import.UploadFile(file)
		if err != nil {
			return c.NoContent(http.StatusBadGateway)
		} else {
			return c.String(http.StatusOK, filename)
		}
	})

	e.POST("/auth/signup", Auth.RegisterHandler)
	e.POST("/auth/login", Auth.LoginHandler)
	e.GET("/auth/logout", Auth.LogoutHandler)
	e.Logger.Fatal(e.Start(":1333"))
}
