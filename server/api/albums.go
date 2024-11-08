package api

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"strconv"
	"strings"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type Albums struct {
	Db *sql.DB
}

func (e *Albums) GetAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := e.getAlbumsRows()

	if err != nil {
		ServeJSONError(w, fmt.Sprintf("GetAlbums %v", err), http.StatusInternalServerError)
		return
	}

	ServeJSON(w, albums, http.StatusOK)
}

func (e *Albums) GetAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/albums/artist/")
	name = "%" + name + "%"

	stmt, err := e.Db.Prepare(`SELECT * FROM album WHERE artist LIKE ?`)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("GetAlbumsByArtist prepare %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("GetAlbumsByArtist %v", err), http.StatusInternalServerError)
		return
	}

	albums, err := handleAlbumRows(rows)
	if err != nil {
		ServeJSONError(w, "failed to get albums by artist", http.StatusInternalServerError)
		return
	}

	if len(albums) > 0 {
		ServeJSON(w, albums, http.StatusOK)
		return
	}

	ServeJSONError(w, fmt.Sprintf("failed to find an album with provided search: %v", name), http.StatusNotFound)
}

func (e *Albums) getAlbumsRows() ([]Album, error) {
	rows, err := e.Db.Query("SELECT * FROM album")

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

func (e *Albums) AddAlbum(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()

	requiredKeys := []string{"title", "artist", "price"}
	for _, value := range requiredKeys {
		if _, ok := parameters[value]; !ok {
			ServeJSONError(w, fmt.Sprintf("must pass in a '%v'", value), http.StatusBadRequest)
			return
		}
	}

	price, _ := strconv.ParseFloat(parameters.Get("price"), 32)

	album := Album{
		Title:  parameters.Get("title"),
		Artist: parameters.Get("artist"),
		Price:  float32(price),
	}

	stmt, err := e.Db.Prepare(`INSERT INTO album (title, artist, price) VALUES (?, ?, ?)`)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("AddAlbum prepare %v", err), http.StatusInternalServerError)
		return
	}

	result, err := stmt.Exec(album.Title, album.Artist, album.Price)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("AddAlbum: %v", err), http.StatusInternalServerError)
		return
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		ServeJSONError(w, "failed to create album", http.StatusBadRequest)
	} else {
		album.ID = lastId
		ServeJSON(w, album, http.StatusOK)
	}
}

func (e *Albums) GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/albums/")

	stmt, err := e.Db.Prepare(`SELECT * FROM album WHERE id = ?`)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("GetAlbumsByArtist prepare %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("GetAlbumsByArtist %v", err), http.StatusInternalServerError)
		return
	}

	albums, err := handleAlbumRows(rows)
	if err != nil {
		ServeJSONError(w, "failed to get album by id", http.StatusInternalServerError)
		return
	}

	if len(albums) == 0 {
		ServeJSONError(w, "album not found", http.StatusNotFound)
	} else {
		ServeJSON(w, albums, http.StatusOK)
	}
}

func (e *Albums) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	_, err := e.Db.Exec("DELETE FROM album WHERE id = ? LIMIT 1", strings.TrimPrefix(r.URL.Path, "/albums/"))
	if err != nil {
		ServeJSONError(w, "could not delete album", http.StatusInternalServerError)
		return
	}

	ServeJSON(w, map[string]any{"message": "album successfully removed"}, http.StatusOK)
}

func (e *Albums) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query()
	id := strings.TrimPrefix(r.URL.Path, "/albums/")

	var keys []string
	var values []any

	if title, ok := parameters["title"]; ok {
		keys = append(keys, "title")
		values = append(values, title[0])
	}
	if artist, ok := parameters["artist"]; ok {
		keys = append(keys, "artist")
		values = append(values, artist[0])
	}
	if price, ok := parameters["price"]; ok {
		keys = append(keys, "price")
		values = append(values, price[0])
	}

	dynamicSql := `UPDATE album SET `
	for i, key := range keys {
		dynamicSql += key + ` = ?`
		if i < len(keys)-1 {
			dynamicSql += `, `
		}
	}
	dynamicSql += ` WHERE id = ?`
	values = append(values, id)

	stmt, err := e.Db.Prepare(dynamicSql)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("UpdateAlbum prepare %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		ServeJSONError(w, "could not update album", http.StatusBadRequest)
		return
	}

	ServeJSON(w, map[string]any{"message": "album successfully updated"}, http.StatusOK)
}

func (e *Albums) AddRandom(w http.ResponseWriter, r *http.Request) {
	name := gofakeit.Name()
	title := gofakeit.Slogan()
	price := gofakeit.Float32Range(1, 100)

	album := Album{
		Title:  title,
		Artist: name,
		Price:  price,
	}

	stmt, err := e.Db.Prepare(`INSERT INTO album (title, artist, price) VALUES (?, ?, ?)`)
	if err != nil {
		ServeJSONError(w, fmt.Sprintf("UpdateAlbum prepare %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(album.Title, album.Artist, album.Price)
	if err != nil {
		ServeJSONError(w, "failed to create random album", http.StatusInternalServerError)
		return
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		ServeJSON(w, map[string]any{"errors": "failed to create random album"}, http.StatusInternalServerError)
	} else {
		album.ID = lastId
		ServeJSON(w, album, http.StatusOK)
	}
}
