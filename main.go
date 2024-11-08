package main

import (
	"github.com/joho/godotenv"
	"go-web-service/server/api"
	"go-web-service/server/utils"
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
	db := utils.DatabaseInit()
	endpoints := &api.Albums{Db: db}

	err := api.SetupRouter(endpoints)
	if err != nil {
		panic("failed to setup router")
	}
}
