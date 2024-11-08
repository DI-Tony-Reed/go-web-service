package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ServeJSON(w http.ResponseWriter, data any, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return errors.New("failed to encode JSON")
	}
	return nil
}
