package main

import (
	"github.com/joho/godotenv"
	"go-web-service/server/rest"
	"go-web-service/server/utils"
	"net/http"
	"os"
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
	albums := &rest.Albums{Db: db}

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

	err := http.ListenAndServe(":"+os.Getenv("APPLICATION_PORT"), handler)
	if err != nil {
		return
	}
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
