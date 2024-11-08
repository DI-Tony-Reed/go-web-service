package utils

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

var sqlOpen = sql.Open

func DatabaseInit() (*sql.DB, error) {
	config := mysql.Config{
		User:                 os.Getenv("DATABASE_USER"),
		Passwd:               os.Getenv("DATABASE_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DATABASE_ADDRESS") + ":" + os.Getenv("DATABASE_PORT"),
		DBName:               os.Getenv("DATABASE_NAME"),
		AllowNativePasswords: true,
	}

	db, err := sqlOpen("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
