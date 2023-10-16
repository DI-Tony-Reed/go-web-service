package utils

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var environment = "development"

func DatabaseInit() *sql.DB {
	var path string

	// This variable is updated via build flags for prod builds
	if environment == "production" {
		path = ".env.production"
	} else {
		path = ".env.development"
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
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}
