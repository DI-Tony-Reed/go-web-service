package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func databaseInit() *sql.DB {
	// Setup DB connection
	config := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS") + ":" + os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_NAME"),
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
