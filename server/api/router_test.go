package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockRouterAlbums struct{}

func (m *MockRouterAlbums) GetAlbums(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, []string{"Album1", "Album2"}, http.StatusOK)
}

func (m *MockRouterAlbums) AddAlbum(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, "Album added", http.StatusCreated)
}

func (m *MockRouterAlbums) GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, "Album1", http.StatusOK)
}

func (m *MockRouterAlbums) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, "Album updated", http.StatusOK)
}

func (m *MockRouterAlbums) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, "Album deleted", http.StatusOK)
}

func (m *MockRouterAlbums) AddRandom(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, "Random album added", http.StatusCreated)
}

func (m *MockRouterAlbums) GetAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	ServeJSON(w, []string{"Album1", "Album2"}, http.StatusOK)
}

func TestServeJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	ServeJSONError(rr, "error message", http.StatusInternalServerError)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, rr.Code)
	}

	expected := `{"errors":"error message"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Expected body %v, got %v", expected, actual)
		t.Logf("Actual body: %q", actual)
	}
}

func TestSetupRouter(t *testing.T) {
	albums := &MockRouterAlbums{}

	// Get the configured router
	router := SetupRouter(albums)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{method: http.MethodGet, url: "/albums", expectedCode: http.StatusOK},
		{method: http.MethodPut, url: "/albums", expectedCode: http.StatusCreated},
		{method: http.MethodGet, url: "/albums/1", expectedCode: http.StatusOK},
		{method: http.MethodPatch, url: "/albums/1", expectedCode: http.StatusOK},
		{method: http.MethodDelete, url: "/albums/1", expectedCode: http.StatusOK},
		{method: http.MethodPut, url: "/albums/random", expectedCode: http.StatusCreated},
		{method: http.MethodGet, url: "/albums/artist/1", expectedCode: http.StatusOK},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.expectedCode {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedCode)
		}
	}
}

func TestSetupRouter__InvalidRoutes(t *testing.T) {
	albums := &MockRouterAlbums{}

	// Get the configured router
	router := SetupRouter(albums)

	tests := []struct {
		method       string
		url          string
		expectedCode int
	}{
		{method: http.MethodPost, url: "/albums", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, url: "/albums/1", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPost, url: "/albums/random", expectedCode: http.StatusMethodNotAllowed},
		{method: http.MethodPut, url: "/albums/artist/1", expectedCode: http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.expectedCode {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedCode)
		}
	}
}
