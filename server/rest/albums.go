package rest

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"go-web-service/server/models"
	"log"
	"net/http"
	"strconv"
)

type Env struct {
	Db *sql.DB
}

func (e *Env) GetAlbums(c *gin.Context) {
	albums, err := e.getAlbumsRows()

	if err != nil {
		fmt.Printf("GetAlbums %v", err)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func (e *Env) GetAlbumsByArtist(c *gin.Context) {
	name := c.Param("artist")

	rows, err := e.Db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		log.Fatalf("GetAlbumsByArtist %q: %v", name, err)
	}

	albums, err := handleAlbumRows(rows)
	if err != nil {
		log.Fatalf("failed to get albums by artist")
	}

	c.IndentedJSON(http.StatusCreated, albums)
}

func (e *Env) getAlbumsRows() ([]models.Album, error) {
	rows, err := e.Db.Query("SELECT * FROM album")

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

func (e *Env) AddAlbum(c *gin.Context) {
	postParameters := c.Request.URL.Query()

	requiredParametersKeys := []string{"title", "artist", "price"}
	for _, value := range requiredParametersKeys {
		if _, ok := postParameters[value]; !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "must pass in a '" + value + "'"})
			return
		}
	}

	price, _ := strconv.ParseFloat(postParameters.Get("price"), 32)

	album := models.Album{
		Title:  postParameters.Get("title"),
		Artist: postParameters.Get("artist"),
		Price:  float32(price),
	}

	result, err := e.Db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ? ,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		log.Fatalf("AddAlbum: %v", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("addAlbum: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to create album"})
	} else {
		album.ID = lastId
		c.IndentedJSON(http.StatusOK, album)
	}
}

func (e *Env) GetAlbumByID(c *gin.Context) {
	id := c.Param("id")
	rows, err := e.Db.Query("SELECT * FROM album WHERE id = ?", id)

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

func (e *Env) DeleteAlbum(c *gin.Context) {
	_, err := e.Db.Exec("DELETE FROM album WHERE id = ?", c.Param("id"))
	if err != nil {
		log.Fatalf("failed to delete album")
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "album successfully removed"})
}

func (e *Env) UpdateAlbum(c *gin.Context) {
	parameters := c.Request.URL.Query()

	var keys []string
	var values []any

	for key, value := range parameters {
		keys = append(keys, key)
		values = append(values, value[0])
	}

	dynamicSql := "UPDATE album SET "
	for key, value := range keys {
		dynamicSql += value + " = ?"

		if (key + 1) < len(keys) {
			dynamicSql += ", "
		} else {
			dynamicSql += " "
		}
	}
	dynamicSql = dynamicSql + "WHERE id = ?"
	values = append(values, c.Param("id"))

	_, err := e.Db.Exec(dynamicSql, values...)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not update album"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "album successfully updated"})
}

func (e *Env) AddRandom(c *gin.Context) {
	name := gofakeit.Name()
	title := gofakeit.Slogan()
	price := gofakeit.Float32Range(1.00, 50.00)

	album := models.Album{
		Title:  title,
		Artist: name,
		Price:  price,
	}

	result, err := e.Db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ? ,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		log.Fatalf("AddAlbum: %v", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("addAlbum: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to create random album"})
	} else {
		album.ID = lastId
		c.IndentedJSON(http.StatusOK, album)
	}
}
