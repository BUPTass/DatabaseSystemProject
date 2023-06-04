package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type databaseInfo struct {
	Name  string
	Value string
}

const (
	dataDbName_ string = "LTE"
)

var db *sqlx.DB

func database_init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username_, password_, ip_, port_, dataDbName_)
	db, _ = sqlx.Open("mysql", dsn)
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

func DatabaseInfo(c echo.Context) error {
	//"/manage/databaseInfo?item=itemname&condition=string"
	item := c.QueryParam("item")
	cond := c.QueryParam("condition")
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	//if ask basic database information
	if cond == "databaseinfo" {
		var ans struct {
			Username          string
			UserDbName        string
			UserinfoTableName string
			Ip                string
			DataDbName        string
			Port              int
		}
		ans.DataDbName = dataDbName_
		ans.UserDbName = userDbName_
		ans.Username = username_
		ans.UserinfoTableName = userinfoTableName_
		ans.Ip = ip_
		ans.Port = port_
		return c.JSON(http.StatusOK, ans)
	}
	var rows *sql.Rows
	if cond != "" {
		rows, err = db.Query(fmt.Sprintf("show %s like '%s'", item, cond))
	} else {
		rows, err = db.Query(fmt.Sprintf("show %s", item))
	}
	defer rows.Close()
	if err != nil {
		return err
	}
	var ans []databaseInfo
	var tmp databaseInfo
	for rows.Next() {
		rows.Scan(&tmp.Name, &tmp.Value)
		ans = append(ans, tmp)
	}
	return c.JSON(http.StatusOK, ans)
}
func DatabaseConnection(c echo.Context) error {
	//"/manage/databaseConnection?condition=string"
	cond := c.QueryParam("condition")
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	//query
	//ip:查看当前连接中各个IP的连接数
	//user:查看当前连接中各个用户的连接数
	//list:查看当前数据库的连接情况
	switch {
	case cond == "ip":
		rows, err := db.Query("select SUBSTRING_INDEX(host,':',1) as ip , count(*) from information_schema.processlist group by ip")
		defer rows.Close()
		if err != nil {
			return err
		}
		type info struct {
			Ipname string
			Cnt    int
		}
		var ans []info
		var tmp info
		for rows.Next() {
			rows.Scan(&tmp.Ipname, &tmp.Cnt)
			ans = append(ans, tmp)
		}
		return c.JSON(http.StatusOK, ans)
	case cond == "user":
		rows, err := db.Query("select USER , count(*) from information_schema.processlist group by USER")
		defer rows.Close()
		if err != nil {
			return err
		}
		type info struct {
			User string
			Cnt  int
		}
		var ans []info
		var tmp info
		for rows.Next() {
			rows.Scan(&tmp.User, &tmp.Cnt)
			ans = append(ans, tmp)
		}
		return c.JSON(http.StatusOK, ans)
	case cond == "list":
		rows, err := db.Query("show full processlist")
		defer rows.Close()
		if err != nil {
			return err
		}
		type allinfo struct {
			Id      int
			User    string
			Host    string
			Db      string
			Command string
			Time    int
			State   string
			Info    string
		}
		var ans []allinfo
		var tmp allinfo
		for rows.Next() {
			rows.Scan(&tmp.Id, &tmp.User, &tmp.Host, &tmp.Db, &tmp.Command, &tmp.Time, &tmp.State, &tmp.Info)
			ans = append(ans, tmp)
		}
		return c.JSON(http.StatusOK, ans)
	}
	return c.String(http.StatusBadRequest, "please check your condition")
}
func SetDatabase(c echo.Context) error {
	//"/manage/database?item=itemname&value=value"
	item := c.QueryParam("item")
	cond := c.QueryParam("value")
	//check session
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	if sess.Values["level"] != 0 {
		return c.String(http.StatusForbidden, "insufficient permissions")
	}
	_, err = db.Exec(fmt.Sprintf("set global %s=%s", item, cond))
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "set success")
}
