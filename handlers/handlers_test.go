package handlers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/omegabytes/services-api/store"
	"github.com/stretchr/testify/assert"
)

// Only a few example tests. I didn't implement mocks for the database in Store,
// so tests on the happy path will absolutely break :)

func TestHandlers(t *testing.T) {
	h := Handler{
		Store: store.Store{},
	}

	t.Run("ListService rejects when offset is not a number", func(t *testing.T) {
		w := httptest.NewRecorder()
		h.ListServiceHandler(w, httptest.NewRequest("GET", "/services?offset=bacon", nil))

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, []byte("Invalid offset\n"), body)
	})

	t.Run("SearchService Should return 400 when search term is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		h.SearchServiceHandler(w, httptest.NewRequest("GET", "/services?search=", nil))

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, []byte("Invalid search\n"), body)
	})

	t.Run("SearchService Should return 400 when search term is too long", func(t *testing.T) {
		w := httptest.NewRecorder()
		bigString := "/services?search=thequickbrownfoxjumpsoverthelazydogthequickbrownfoxjumpsoverthelazydogthequickbrownfoxjumpsoverthelazydogthequickbrownfoxjumpsoverthelazydog"
		h.SearchServiceHandler(w, httptest.NewRequest("GET", bigString, nil))

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, 400, resp.StatusCode)
		assert.Equal(t, []byte("Search term to long\n"), body)
	})

	t.Run("GetService should return a 400 when id is missing", func(t *testing.T) {
		t.Skip()
	})
}
