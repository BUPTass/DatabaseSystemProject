package main

import (
	"DatabaseSystemProject/Export"
	"DatabaseSystemProject/Import"
	"DatabaseSystemProject/Query"
	"DatabaseSystemProject/controllers"
	"database/sql"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
)

var dbHost, dbPort, dbUser, dbPassword, dbName string

func main() {
	// Read the database connection details from environmental variables
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                        // Allow all origins
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}, // Allow specified methods
	}))

	// Construct the database connection string
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dbURI)
	// Connect to the MySQL database.
	//db, err := sql.Open("mysql", "root:1taNWY1vXdTc4_-j@tcp(127.0.0.1:3306)/LTE")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Database System Project API backend")
	})

	//route
	//user management
	e.POST("/login", controllers.Login)                               //user login
	e.GET("/logout", controllers.Logout)                              //user logout
	e.GET("/show/users", controllers.GetUsers)                        //show all users
	e.GET("/show/users/unconfirmed", controllers.GetUnconfirmedUsers) //show unconfirmed users
	e.POST("/add/user", controllers.AddUser)                          //add user
	e.DELETE("/delete/user", controllers.DeleteUser)                  //delete user
	e.POST("/register", controllers.Register)                         //user register

	//database management
	e.GET("/manage/databaseInfo", controllers.DatabaseInfo) //check database info
	e.GET("/manage/databaseConnection", controllers.DatabaseConnection)
	e.POST("/manage/database", controllers.SetDatabase)

	//3.6 3.7
	e.POST("/calc/C2Inew", controllers.C2InewCalc)
	e.POST("/calc/C2I3", controllers.C2I3Calc)
	//3.8
	e.POST("/calc/MRO", controllers.MROMREcalc)

	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "hello") })

	e.Static("/download", "./download")
	e.POST("/import/tbCell", func(c echo.Context) error {
		path := c.FormValue("path")

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err = Import.AddtbCell(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbCell Added")
		}
	})
	e.POST("/import/tbKPI", func(c echo.Context) error {
		path := c.FormValue("path")

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err = Import.AddtbKPI(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbKPI Added")
		}
	})
	e.POST("/import/tbPRB", func(c echo.Context) error {
		path := c.FormValue("path")

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err = Import.AddtbPRB(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbRPB Added")
		}
	})
	e.POST("/import/tbMROData", func(c echo.Context) error {
		path := c.FormValue("path")

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err = Import.AddtbMROData(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbMROData Added")
		}
	})
	e.POST("/import/tbC2I", func(c echo.Context) error {
		path := c.FormValue("path")

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		if len(path) == 0 {
			return c.String(http.StatusBadRequest, "No path provided")
		}
		err = Import.AddtbC2I(db, path)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		} else {
			return c.String(http.StatusOK, "New tbC2I Added")
		}
	})

	e.GET("/export", func(c echo.Context) error {
		path := c.FormValue("path")
		table := c.FormValue("table")

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

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

		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

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
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		jsonByte, err := Query.GetAllSectorNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})
	e.GET("/query/tbCell", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		query := c.FormValue("sector")
		jsonByte, err := Query.GetCellInfo(db, query)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})

	e.GET("/query/enodeb_name", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		jsonByte, err := Query.GetAllEnodebNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})
	e.GET("/query/enodeb", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		query := c.FormValue("enodeb")
		jsonByte, err := Query.GetEnodeb(db, query)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})

	e.GET("/query/kpi/sector_name", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		jsonByte, err := Query.GetKPISectorNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})
	e.GET("/query/kpi", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		query := c.FormValue("sector")
		jsonByte, err := Query.GetKPIInfoBySectorName(db, query)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, jsonByte)
		}
	})

	e.GET("/query/tbPRBNew/gen", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		path := c.FormValue("path")
		ret, err := Query.GeneratePRBNewTable(db, path)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.String(http.StatusOK, ret)
		}
	})
	e.GET("/query/tbPRB/sector_name", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		ret, err := Query.GetPRBSectorNames(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, ret)
		}
	})
	e.GET("/query/tbPRB", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		sector := c.FormValue("sector")
		ret, err := Query.GetPRBBySectorName(db, sector)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, ret)
		}
	})
	e.GET("/query/tbPRBNew", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		sector := c.FormValue("sector")
		ret, err := Query.GetPRBNewBySectorName(db, sector)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSON(http.StatusOK, ret)
		}
	})

	e.GET("/query/community", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || !(sess.Values["level"] == 0 || sess.Values["level"] == 1) {
			return c.NoContent(http.StatusUnauthorized)
		}

		ret, err := Query.GetCommunity(db)
		if err != nil {
			return c.String(http.StatusBadGateway, err.Error())
		} else {
			return c.JSONBlob(http.StatusOK, ret)
		}
	})

	e.Logger.Fatal(e.Start(":1333"))
}
