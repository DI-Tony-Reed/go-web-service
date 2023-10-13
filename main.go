package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-web-service/models"
)

func main() {
	db = databaseInit()

	// Setup gin router
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.GET("/albums/artist/:artist", getAlbumsByArtist)
	router.POST("/albums", addAlbum)
	//router.POST("/albums/:id", updateAlbum)

	err := router.Run(":8081")
	if err != nil {
		log.Fatal(err)
	}
}

var db *sql.DB

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	albums, err := getAlbumsRows()

	if err != nil {
		fmt.Printf("getAlbumsByArtist %v", err)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumsByArtist(c *gin.Context) {
	name := c.Param("artist")

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		log.Fatalf("getAlbumsByArtist %q: %v", name, err)
	}

	albums, err := handleAlbumRows(rows)
	if err != nil {
		log.Fatalf("failed to get albums by artist")
	}

	c.IndentedJSON(http.StatusCreated, albums)
}

func getAlbumsRows() ([]models.Album, error) {
	rows, err := db.Query("SELECT * FROM album")

	if err != nil {
		return nil, fmt.Errorf("handleAlbumRows %v", err)
	}

	return handleAlbumRows(rows)
}

func handleAlbumRows(rows *sql.Rows) ([]models.Album, error) {
	// Albums slice to hold db rows
	var albums []models.Album

	defer rows.Close()

	// Loop rows using Scan to assign to struct fields
	for rows.Next() {
		var album models.Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("handleAlbumRows %v", err)
		}

		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("handleAlbumRows %v", err)
	}

	return albums, nil
}

// addAlbum adds an album to the database
// returning the new album in the response
func addAlbum(c *gin.Context) {
	postParameters := c.Request.URL.Query()

	if _, ok := postParameters["title"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "must pass in a 'title'"})
		return
	}

	if _, ok := postParameters["artist"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "must pass in an 'artist'"})
		return
	}

	if _, ok := postParameters["price"]; !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "must pass in a 'price'"})
		return
	}

	price, _ := strconv.ParseFloat(postParameters.Get("price"), 32)

	album := models.Album{
		Title:  postParameters.Get("title"),
		Artist: postParameters.Get("artist"),
		Price:  float32(price),
	}

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ? ,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		log.Fatalf("addAlbum: %v", err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		log.Fatalf("addAlbum: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to create album"})
	} else {
		c.IndentedJSON(http.StatusOK, album)
	}
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	rows, err := db.Query("SELECT * FROM album WHERE id = ?", id)

	if err != nil {
		log.Fatalf("getAlbumsByID %q: %v", id, err)
	}

	albums, err := handleAlbumRows(rows)

	if err != nil {
		log.Fatal("failure within handleAlbumRows")
	}

	if len(albums) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	} else {
		c.IndentedJSON(http.StatusOK, albums)
	}
}

// TODO add update method
