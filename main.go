package main

import (
	"database/sql"
	"example/web-service-gin/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	db = databaseInit()

	// Setup gin router
	router := gin.Default()
	router.GET("/albums", getAlbums)
	//router.GET("/albums/:id", getAlbumByID)
	router.GET("/albums/artist/:artist", getAlbumsByArtistJSON)
	//router.POST("/albums", postAlbums)

	alb, err := getAlbumsRows()
	fmt.Printf("Album(s) found: %v\n", alb)

	err = router.Run("localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
}

var albums []models.Album

var db *sql.DB

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	albums, err := getAlbumsRows()

	if err != nil {
		fmt.Printf("getAlbumsByArtist %v", err)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// Query by artist name
func getAlbumsByArtist(name string) ([]models.Album, error) {
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("getAlbumsByArtist %q: %v", name, err)
	}

	return handleAlbumRows(rows)
}

func getAlbumsByArtistJSON(c *gin.Context) {
	albums, err := getAlbumsByArtist(c.Param("artist"))

	if err != nil {
		panic("failed to get by artist name")
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

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum models.Album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// addAlbum adds an album to the database
// returning the album ID of the new entry
func addAlbum(album models.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ? ,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
//func getAlbumByID(c *gin.Context) {
//	id := c.Param("id")
//
//	// Loop over the list of albums, looking for
//	// an album whose ID value matches the parameter.
//	for _, album := range albums {
//		if album.ID == id {
//			c.IndentedJSON(http.StatusOK, album)
//			return
//		}
//	}
//
//	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
//}
