package utils

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabaseInit(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Replace the sqlOpen function with a mock
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	// Expect a ping to the database
	mock.ExpectPing()

	// Call the DatabaseInit function
	db, err = DatabaseInit()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Ensure the database connection is not nil
	if db == nil {
		t.Fatal("Expected non-nil database connection")
	}

	// Ensure the database can be pinged
	err = db.Ping()
	if err != nil {
		t.Fatalf("Expected successful ping, got error: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %v", err)
	}
}

func TestDatabaseInit_SQLOpenError(t *testing.T) {
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("sql open error")
	}

	db, err := DatabaseInit()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if db != nil {
		t.Fatal("Expected nil database connection")
	}
}

func TestDatabaseInit_DBPingError(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Replace the sqlOpen function with a mock
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	// Expect a ping to the database to return an error
	mock.ExpectPing().WillReturnError(errors.New("ping error"))

	// Call the DatabaseInit function
	db, err = DatabaseInit()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Ensure the database connection is nil
	if db != nil {
		t.Fatal("Expected nil database connection")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %v", err)
	}
}
