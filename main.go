package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-web-service/server/rest"
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
	env := &rest.Env{Db: db}

	// Setup gin router
	router := gin.Default()

	router.GET("/albums", env.GetAlbums)
	router.GET("/albums/:id", env.GetAlbumByID)
	router.GET("/albums/artist/:artist", env.GetAlbumsByArtist)

	router.PUT("/albums", env.AddAlbum)
	router.PUT("/albums/random", env.AddRandom)
	router.PATCH("/albums/:id", env.UpdateAlbum)
	router.DELETE("/albums/:id", env.DeleteAlbum)

	err := router.Run(":" + os.Getenv("APPLICATION_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
