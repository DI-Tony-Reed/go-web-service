package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func databaseInit() *sql.DB {
	// Setup DB connection
	config := mysql.Config{
		User:                 "user",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "mysql:3306", // Correlates to the mysql service name and port
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	// Get DB handle
	var err error
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}
