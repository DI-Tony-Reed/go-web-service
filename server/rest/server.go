package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ServeJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return errors.New("failed to encode JSON")
	}
	return nil
}
