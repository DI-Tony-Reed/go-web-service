package main

import (
	"github.com/joho/godotenv"
	"go-web-service/api"
	"go-web-service/utils"
	"net/http"
	"os"
)

var environment = "development"

func init() {
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
}

func main() {
	db, err := utils.DatabaseInit()
	if err != nil {
		panic("failed to connect to database")
	}

	endpoints := &api.Albums{Db: db}

	router := api.SetupRouter(endpoints)
	err = http.ListenAndServe(":"+os.Getenv("APPLICATION_PORT"), router)
	if err != nil {
		panic(err)
	}
}
