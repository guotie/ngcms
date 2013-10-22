package models

import (
	"database/sql"
	"github.com/robfig/revel"
)

var (
	db *sql.DB
)

func init_db() {
	var err error
	var dsn string

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		revel.ERROR.Panicf("Open database failed: %s\n", err.Error())
	}

	err = db.Ping()
	if err != nil {
		revel.ERROR.Panicf("Connect to database failed: %s\n", err.Error())
	}
}

// must success or else panic
func get_db() *sql.DB {
	err := db.Ping()
	if err != nil {
		revel.ERROR.Panicf("Connect to database failed: %s\n", err.Error())
	}

	return db
}
