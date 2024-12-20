package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}
	return db, mock
}

func sendMockHTTPRequest(t *testing.T, handler http.HandlerFunc, method, url string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	return rr
}

func TestGetAlbums(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
		AddRow(1, "Album1", "Artist1", 10.99).
		AddRow(2, "Album2", "Artist2", 12.99)
	mock.ExpectQuery("SELECT \\* FROM album").WillReturnRows(rows)

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbums, http.MethodGet, "/albums")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"title":"Album1","artist":"Artist1","price":10.99},{"id":2,"title":"Album2","artist":"Artist2","price":12.99}]`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbums_Error(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM album").WillReturnError(fmt.Errorf("query error"))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbums, http.MethodGet, "/albums")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"GetAlbums handleAlbumRows query error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumsByArtist(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE artist LIKE \?`).
		ExpectQuery().
		WithArgs("%Artist1%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow(1, "Album1", "Artist1", 10.99).
			AddRow(2, "Album2", "Artist2", 12.99))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbumsByArtist, http.MethodGet, "/albums/artist/Artist1")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"title":"Album1","artist":"Artist1","price":10.99},{"id":2,"title":"Album2","artist":"Artist2","price":12.99}]`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumsByArtist_HandleAlbumRowsError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE artist LIKE \?`).
		ExpectQuery().
		WithArgs("%Artist1%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow(1, "Album1", "Artist1", 10.99).
			AddRow(2, "Album2", "Artist2", 12.99))

	albums := &Albums{Db: db}

	// Mock the handleAlbumRows function to return an error
	originalHandleAlbumRows := handleAlbumRows
	handleAlbumRows = func(rows *sql.Rows) ([]Album, error) {
		return nil, fmt.Errorf("mocked error")
	}
	// Reset after test
	defer func() { handleAlbumRows = originalHandleAlbumRows }()

	rr := sendMockHTTPRequest(t, albums.GetAlbumsByArtist, http.MethodGet, "/albums/artist/Artist1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"failed to get albums by artist"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumsByArtist_PrepareError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE artist LIKE \?`).WillReturnError(fmt.Errorf("prepare error"))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbumsByArtist, http.MethodGet, "/albums/artist/Artist1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"GetAlbumsByArtist prepare prepare error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumsByArtist_QueryError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE artist LIKE \?`).
		ExpectQuery().
		WithArgs("%Artist1%").
		WillReturnError(fmt.Errorf("query error"))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbumsByArtist, http.MethodGet, "/albums/artist/Artist1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"GetAlbumsByArtist query error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumsByArtist_NoAlbumsFound(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE artist LIKE \?`).
		ExpectQuery().
		WithArgs("%NonExistentArtist%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbumsByArtist, http.MethodGet, "/albums/artist/NonExistentArtist")

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected := `{"errors":"failed to find an album with provided search: %NonExistentArtist%"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddAlbum(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs("Album1", "Artist1", 10.00).
		WillReturnResult(sqlmock.NewResult(1, 1))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.AddAlbum, http.MethodPut, "/albums?title=Album1&artist=Artist1&price=10")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":1,"title":"Album1","artist":"Artist1","price":10}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddAlbum_MissingParameters(t *testing.T) {
	db, _ := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.AddAlbum, http.MethodPut, "/albums?title=Album1&artist=Artist1")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"errors":"must pass in a 'price'"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestAddAlbum_PrepareError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		WillReturnError(fmt.Errorf("prepare error"))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.AddAlbum, http.MethodPut, "/albums?title=Album1&artist=Artist1&price=10.99")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"AddAlbum prepare prepare error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddAlbum_InsertError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs("Album1", "Artist1", 10.00).
		WillReturnError(fmt.Errorf("insert error"))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.AddAlbum, http.MethodPut, "/albums?title=Album1&artist=Artist1&price=10")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"AddAlbum: insert error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddAlbum_LastInsertIdError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs("Album1", "Artist1", 10.00).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("last insert id error")))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.AddAlbum, http.MethodPut, "/albums?title=Album1&artist=Artist1&price=10")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"errors":"failed to create album"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumByID(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE id = \?`).
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow(1, "Album1", "Artist1", 10.99))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.GetAlbumByID, http.MethodGet, "/albums/1")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"title":"Album1","artist":"Artist1","price":10.99}]`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumByID_HandleAlbumRowsError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare(`SELECT \* FROM album WHERE id = \?`).
		ExpectQuery().
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow(1, "Album1", "Artist1", 10.99))

	albums := &Albums{Db: db}

	// Mock the handleAlbumRows function to return an error
	originalHandleAlbumRows := handleAlbumRows
	handleAlbumRows = func(rows *sql.Rows) ([]Album, error) {
		return nil, fmt.Errorf("mocked error")
	}
	// Reset after test
	defer func() { handleAlbumRows = originalHandleAlbumRows }()

	rr := sendMockHTTPRequest(t, albums.GetAlbumByID, http.MethodGet, "/albums/1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"failed to get album by id"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestGetAlbumByID_Errors(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	// Prepare error case
	mock.ExpectPrepare(`SELECT \* FROM album WHERE id = \?`).WillReturnError(fmt.Errorf("prepare error"))

	rr := sendMockHTTPRequest(t, albums.GetAlbumByID, http.MethodGet, "/albums/1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"GetAlbumsByArtist prepare prepare error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}

	// Query error case
	mock.ExpectPrepare(`SELECT \* FROM album WHERE id = \?`).
		ExpectQuery().
		WithArgs("1").
		WillReturnError(fmt.Errorf("query error"))

	rr = sendMockHTTPRequest(t, albums.GetAlbumByID, http.MethodGet, "/albums/1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected = `{"errors":"GetAlbumsByArtist query error"}`
	actual = strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}

	// No albums found case
	mock.ExpectPrepare(`SELECT \* FROM album WHERE id = \?`).
		ExpectQuery().
		WithArgs("999").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}))

	rr = sendMockHTTPRequest(t, albums.GetAlbumByID, http.MethodGet, "/albums/999")

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expected = `{"errors":"album not found"}`
	actual = strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestDeleteAlbum(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM album WHERE id = \\? LIMIT 1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.DeleteAlbum, http.MethodDelete, "/albums/1")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"album successfully removed"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestDeleteAlbum_Errors(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	// Exec error case
	mock.ExpectExec("DELETE FROM album WHERE id = \\? LIMIT 1").
		WithArgs("1").
		WillReturnError(fmt.Errorf("exec error"))

	rr := sendMockHTTPRequest(t, albums.DeleteAlbum, http.MethodDelete, "/albums/1")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"could not delete album"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestUpdateAlbum(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare("UPDATE album SET title = \\?, artist = \\?, price = \\? WHERE id = \\?").
		ExpectExec().
		WithArgs("UpdatedTitle", "UpdatedArtist", "20", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.UpdateAlbum, http.MethodPatch, "/albums/1?title=UpdatedTitle&artist=UpdatedArtist&price=20")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"album successfully updated"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestUpdateAlbum_Errors(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	// Prepare error case
	mock.ExpectPrepare("UPDATE album SET title = \\?, artist = \\?, price = \\? WHERE id = \\?").
		WillReturnError(fmt.Errorf("prepare error"))

	rr := sendMockHTTPRequest(t, albums.UpdateAlbum, http.MethodPatch, "/albums/1?title=UpdatedTitle&artist=UpdatedArtist&price=20")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"UpdateAlbum prepare prepare error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}

	// Exec error case
	mock.ExpectPrepare("UPDATE album SET title = \\?, artist = \\?, price = \\? WHERE id = \\?").
		ExpectExec().
		WithArgs("UpdatedTitle", "UpdatedArtist", "20.99", "1").
		WillReturnError(fmt.Errorf("exec error"))

	rr = sendMockHTTPRequest(t, albums.UpdateAlbum, http.MethodPatch, "/albums/1?title=UpdatedTitle&artist=UpdatedArtist&price=20.99")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected = `{"errors":"could not update album"}`
	actual = strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddRandom(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	albums := &Albums{Db: db}

	rr := sendMockHTTPRequest(t, albums.AddRandom, http.MethodPost, "/albums/random")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddRandom_Errors(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	// Exec error case
	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("exec error"))

	rr := sendMockHTTPRequest(t, albums.AddRandom, http.MethodPost, "/albums/random")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"failed to create random album"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddRandom_PrepareError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	// Prepare error case
	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		WillReturnError(fmt.Errorf("prepare error"))

	rr := sendMockHTTPRequest(t, albums.AddRandom, http.MethodPost, "/albums/random")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"UpdateAlbum prepare prepare error"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddRandom_LastInsertIdError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	albums := &Albums{Db: db}

	mock.ExpectPrepare("INSERT INTO album \\(title, artist, price\\) VALUES \\(\\?, \\?, \\?\\)").
		ExpectExec().
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("last insert id error")))

	rr := sendMockHTTPRequest(t, albums.AddRandom, http.MethodPost, "/albums/random")

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := `{"errors":"failed to create random album"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}
