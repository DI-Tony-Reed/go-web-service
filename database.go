package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func databaseInit() *sql.DB {
	// Setup DB connection
	// TODO use os.Getenv and move the `export` statements to a DockerFile so it's automatic
	config := mysql.Config{
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:33307",
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
