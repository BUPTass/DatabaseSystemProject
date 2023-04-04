package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/show", func(c echo.Context) error {
		a := c.QueryParam("a")
		b := c.QueryParam("b")
		return c.String(http.StatusOK, a+b)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
