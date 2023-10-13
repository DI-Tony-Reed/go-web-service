package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func databaseInit() *sql.DB {
	var path string

	// This is dictated via build flags
	if environment == "development" {
		path = ".env.development"
	} else {
		path = ".env"
	}

	err := godotenv.Load(path)

	if err != nil {
		panic("failed to load .env file")
	}

	// Setup DB connection
	config := mysql.Config{
		User:                 os.Getenv("DATABASE_USER"),
		Passwd:               os.Getenv("DATABASE_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DATABASE_ADDRESS") + ":" + os.Getenv("DATABASE_PORT"),
		DBName:               os.Getenv("DATABASE_NAME"),
		AllowNativePasswords: true,
	}

	// Get DB handle
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
