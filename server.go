package main

import (
	"DatabaseSystemProject/Auth"
	"DatabaseSystemProject/Export"
	"DatabaseSystemProject/Import"
	"DatabaseSystemProject/Query"
	"database/sql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
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
	// Connect to the MySQL database.
	db, err := sql.Open("mysql", "root:1taNWY1vXdTc4_-j@tcp(127.0.0.1:3306)/LTE")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()
	e.Static("/download", "static")
	e.POST("/import/tbCell", func(c echo.Context) error {
		path := c.FormValue("path")

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err := Import.AddtbCell(db, path)
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
		err := Import.AddtbKPI(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbKPI Added")
		}
	})
	e.POST("/import/tbPRB", func(c echo.Context) error {
		path := c.FormValue("path")

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err := Import.AddtbPRB(db, path)
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
		err := Import.AddtbMROData(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbMROData Added")
		}
	})
	e.GET("/export", func(c echo.Context) error {
		path := c.FormValue("path")
		table := c.FormValue("table")

		if len(table) == 0 {
			return c.String(http.StatusBadRequest, "Missing table")
		}
		ret, err := Export.TableAsCSV(db, path, table)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		} else {
			return c.String(http.StatusOK, ret)
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

	e.GET("/query/sector_name", func(c echo.Context) error {
		jsonByte, err := Query.GetAllSectorNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})
	e.GET("/query/tbCell", func(c echo.Context) error {
		query := c.FormValue("sector")
		jsonByte, err := Query.GetCellInfo(db, query)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})

	e.GET("/query/enodeb_name", func(c echo.Context) error {
		jsonByte, err := Query.GetAllEnodebNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})
	e.GET("/query/enodeb", func(c echo.Context) error {
		query := c.FormValue("enodeb")
		jsonByte, err := Query.GetEnodeb(db, query)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})

	e.GET("/query/kpi/sector_name", func(c echo.Context) error {
		jsonByte, err := Query.GetKPISectorNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})
	e.GET("/query/kpi", func(c echo.Context) error {
		query := c.FormValue("sector")
		jsonByte, err := Query.GetKPIInfoBySectorName(db, query)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})

	e.GET("/query/tbPRBNew/gen", func(c echo.Context) error {
		path := c.FormValue("path")
		ret, err := Query.GeneratePRBNewTable(db, path)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.String(http.StatusOK, ret)
		}
	})

	e.POST("/auth/signup", Auth.RegisterHandler)
	e.POST("/auth/login", Auth.LoginHandler)
	e.GET("/auth/logout", Auth.LogoutHandler)
	e.Logger.Fatal(e.Start(":1333"))
}
