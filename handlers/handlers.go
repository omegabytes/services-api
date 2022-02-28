package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/omegabytes/services-api/store"
)

type Handler struct {
	Store store.Store
}

// GetServiceHandler fetches a single service using a given service ID.
func (h *Handler) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.GetService")
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	results, err := h.Store.GetService(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(encode)
}

// ListServiceHandler returns a list of services. The maximum returned services is configured using the global
// config.Limit value that is set at runtime. A user-provided offset value is used to fetch subsquent results.
func (h *Handler) ListServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.ListService")
	offset := r.URL.Query().Get("offset") // assume offset = last record shown + 1, handled by the front end
	if offset == "" {
		offset = "0"
	}

	o, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	results, err := h.Store.ListServices(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(encode)
}

func (h *Handler) SearchServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.SearchService")
	searchTerm := r.URL.Query().Get("search")

	if searchTerm == "" {
		http.Error(w, "Invalid search", http.StatusBadRequest)
		return
	}

	if len(searchTerm) > 100 {
		http.Error(w, "Search term to long", http.StatusBadRequest)
		return
	}

	// todo: additional validation to prevent SQL injection etc
	results, err := h.Store.SearchServices(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encode, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(encode)
}
