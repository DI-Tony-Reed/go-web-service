package rest

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"go-web-service/server/models"
	"log"
	"net/http"
	"strings"
)

type Albums struct {
	Db *sql.DB
}

func (e *Albums) GetAlbums(w http.ResponseWriter, req *http.Request) {
	albums, err := e.getAlbumsRows()

	if err != nil {
		http.Error(w, fmt.Sprintf("GetAlbums %v", err), http.StatusInternalServerError)
	}

	err = ServeJSON(w, albums, http.StatusOK)
	if err != nil {
		http.Error(w, fmt.Sprintf("GetAlbums %v", err), http.StatusInternalServerError)
	}
}

func (e *Albums) GetAlbumsByArtist(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimPrefix(req.URL.Path, "/albums/artist/")

	rows, err := e.Db.Query(`SELECT * FROM album WHERE artist LIKE CONCAT('%', ?, '%')`, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("GetAlbumsByArtist %v", err), http.StatusInternalServerError)
	}

	albums, err := handleAlbumRows(rows)
	if err != nil {
		http.Error(w, "failed to get albums by artist", http.StatusInternalServerError)
	}

	if len(albums) > 0 {
		err := ServeJSON(w, albums, http.StatusOK)
		if err != nil {
			http.Error(w, fmt.Sprintf("GetAlbumsByArtist %v", err), http.StatusInternalServerError)
		}

		return
	}

	err = ServeJSON(w, map[string]any{"errors": "failed to find an album with provided search: " + name}, http.StatusNotFound)
	if err != nil {
		http.Error(w, fmt.Sprintf("GetAlbumsByArtist %v", err), http.StatusInternalServerError)
	}
}

func (e *Albums) getAlbumsRows() ([]models.Album, error) {
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

func (e *Albums) AddAlbum(w http.ResponseWriter, req *http.Request) {
	//postParameters := c.Request.URL.Query()

	//requiredParametersKeys := []string{"title", "artist", "price"}
	//for _, value := range requiredParametersKeys {
	//	if _, ok := postParameters[value]; !ok {
	//c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": "must pass in a '" + value + "'"})
	//return
	//}
	//}

	//price, _ := strconv.ParseFloat(postParameters.Get("price"), 32)
	//
	//album := models.Album{
	//	Title:  postParameters.Get("title"),
	//	Artist: postParameters.Get("artist"),
	//	Price:  float32(price),
	//}
	//
	//result, err := e.Db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ? ,?)", album.Title, album.Artist, album.Price)
	//if err != nil {
	//	log.Fatalf("AddAlbum: %v", err)
	//}
	//
	//lastId, err := result.LastInsertId()
	//if err != nil {
	//	log.Fatalf("addAlbum: %v", err)
	//	c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": "failed to create album"})
	//} else {
	//	album.ID = lastId
	//	c.IndentedJSON(http.StatusOK, album)
	//}
}

func (e *Albums) GetAlbumByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/albums/")
	rows, err := e.Db.Query("SELECT * FROM album WHERE id = ?", id)

	if err != nil {
		log.Fatalf("getAlbumsByID %q: %v", id, err)
	}

	albums, err := handleAlbumRows(rows)

	if err != nil {
		log.Fatal("failure within handleAlbumRows")
	}

	if len(albums) == 0 {
		err := ServeJSON(w, map[string]any{"errors": "album not found"}, http.StatusNotFound)
		if err != nil {
			http.Error(w, fmt.Sprintf("GetAlbumByID %v", err), http.StatusInternalServerError)
		}
	} else {
		err := ServeJSON(w, albums, http.StatusOK)
		if err != nil {
			http.Error(w, fmt.Sprintf("GetAlbumByID %v", err), http.StatusInternalServerError)
		}
	}
}

func (e *Albums) DeleteAlbum(w http.ResponseWriter, req *http.Request) {
	//_, err := e.Db.Exec("DELETE FROM album WHERE id = ?", c.Param("id"))
	//if err != nil {
	//	log.Fatalf("failed to delete album")
	//}
	//
	//c.IndentedJSON(http.StatusOK, gin.H{"message": "album successfully removed"})
}

func (e *Albums) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	id := strings.TrimPrefix(r.URL.Path, "/albums/")

	var keys []string
	var values []any

	for key, value := range parameters {
		keys = append(keys, key)
		values = append(values, value[0])
	}

	dynamicSql := `UPDATE album SET `
	for key, value := range keys {
		dynamicSql += value + " = ?"

		if (key + 1) < len(keys) {
			dynamicSql += ", "
		} else {
			dynamicSql += " "
		}
	}
	dynamicSql = dynamicSql + "WHERE id = ?"
	values = append(values, id)

	_, err := e.Db.Exec(dynamicSql, values...)
	if err != nil {
		err := ServeJSON(w, map[string]any{"errors": "could not update album", "error": err}, http.StatusBadRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("UpdateAlbum %v", err), http.StatusInternalServerError)
		}
		return
	}

	err = ServeJSON(w, map[string]any{"message": "album successfully updated"}, http.StatusOK)
	if err != nil {
		http.Error(w, fmt.Sprintf("UpdateAlbum %v", err), http.StatusInternalServerError)
	}
}

func (e *Albums) AddRandom(w http.ResponseWriter, r *http.Request) {
	name := gofakeit.Name()
	title := gofakeit.Slogan()
	price := gofakeit.Float32Range(1, 100)

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
		err := ServeJSON(w, map[string]any{"errors": "failed to create album"}, http.StatusInternalServerError)
		if err != nil {
			http.Error(w, fmt.Sprintf("AddRandom %v", err), http.StatusInternalServerError)
		}
	} else {
		album.ID = lastId
		err := ServeJSON(w, album, http.StatusOK)
		if err != nil {
			http.Error(w, fmt.Sprintf("AddRandom %v", err), http.StatusInternalServerError)
		}
	}
}
