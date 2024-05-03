package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Instance *sql.DB

func InitializeDatabase() {
	var err error
	Instance, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/shop_demo")
	if err != nil {
		panic(err)
	}
	Instance.SetConnMaxLifetime(3 * time.Minute)
	Instance.SetMaxOpenConns(10)
	Instance.SetMaxIdleConns(10)
}
