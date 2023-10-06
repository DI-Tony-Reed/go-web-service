package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	// Setup DB connection
	// TODO use os.Getenv and move the `export` statements to a DockerFile so it's automatic
	config := mysql.Config{
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:33307",
		DBName:               "recordings",
		AllowNativePasswords: true,
	}

	// Get DB handle
	var err error
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// Setup gin router
	router := gin.Default()
	//router.GET("/albums", getAlbums)
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

// Album represents data about a record album.
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var albums []Album

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
func getAlbumsByArtist(name string) ([]Album, error) {
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

func getAlbumsRows() ([]Album, error) {
	rows, err := db.Query("SELECT * FROM album")

	if err != nil {
		return nil, fmt.Errorf("handleAlbumRows %v", err)
	}

	return handleAlbumRows(rows)
}

func handleAlbumRows(rows *sql.Rows) ([]Album, error) {
	// Albums slice to hold db rows
	var albums []Album

	defer rows.Close()

	// Loop rows using Scan to assign to struct fields
	for rows.Next() {
		var album Album
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
	var newAlbum Album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
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
