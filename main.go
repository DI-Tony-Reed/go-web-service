package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-web-service/src/rest"
	"go-web-service/src/utils"
)

func main() {
	db := utils.DatabaseInit()
	env := &rest.Env{Db: db}

	// Setup gin router
	router := gin.Default()

	router.GET("/albums", env.GetAlbums)
	router.GET("/albums/:id", env.GetAlbumByID)
	router.GET("/albums/artist/:artist", env.GetAlbumsByArtist)

	router.PUT("/albums", env.AddAlbum)
	router.PATCH("/albums/:id", env.UpdateAlbum)
	router.DELETE("/albums/:id", env.DeleteAlbum)

	err := router.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
