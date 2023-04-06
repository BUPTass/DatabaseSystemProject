package controllers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	DataDbName string = "tdlte"
)

var db *sqlx.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Username, Password, Ip, Port, DataDbName)
	db, _ = sqlx.Open("mysql", dsn)
	if err := user_info.Ping(); err != nil {
		panic(err)
	}
}
