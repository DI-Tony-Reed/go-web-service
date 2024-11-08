package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type AlbumsInterface interface {
	GetAlbums(w http.ResponseWriter, r *http.Request)
	AddAlbum(w http.ResponseWriter, r *http.Request)
	GetAlbumByID(w http.ResponseWriter, r *http.Request)
	UpdateAlbum(w http.ResponseWriter, r *http.Request)
	DeleteAlbum(w http.ResponseWriter, r *http.Request)
	AddRandom(w http.ResponseWriter, r *http.Request)
	GetAlbumsByArtist(w http.ResponseWriter, r *http.Request)
	GetHandleAlbumRows(rows *sql.Rows) ([]Album, error)
}

func ServeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func ServeJSONError(w http.ResponseWriter, message string, statusCode int) {
	ServeJSON(w, map[string]any{"errors": message}, statusCode)
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func SetupRouter(albums AlbumsInterface) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			albums.GetAlbums(w, r)
		case http.MethodPut:
			albums.AddAlbum(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/albums/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			albums.GetAlbumByID(w, r)
		case http.MethodPatch:
			albums.UpdateAlbum(w, r)
		case http.MethodDelete:
			albums.DeleteAlbum(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/albums/random", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			albums.AddRandom(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/albums/artist/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			albums.GetAlbumsByArtist(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler := corsMiddleware(mux)

	return handler
}
