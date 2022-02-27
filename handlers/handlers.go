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

func (h *Handler) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.GetService")
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "id is required", 400)
		return
	}

	results, err := h.Store.GetService(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	encode, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(encode)
}

func (h *Handler) ListServiceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call.ListService")
	offset := r.URL.Query().Get("offset") // assume offset = last record shown + 1, handled by the front end
	o, err := strconv.Atoi(offset)

	// We could also set offset = 0 no matter what so the user gets results.
	// I am chosing to reject here to show some easy-to-test error handling.
	if err != nil {
		http.Error(w, "Invalid offset", 400)
		return
	}

	results, err := h.Store.ListServices(o)
	if err != nil {
		http.Error(w, err.Error(), 500)
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
		http.Error(w, "Invalid search", 400)
		return
	}

	if len(searchTerm) > 100 {
		http.Error(w, "Search term to long", 400)
		return
	}

	results, err := h.Store.SearchServices(searchTerm)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	encode, _ := json.Marshal(results)
	w.WriteHeader(http.StatusOK)
	w.Write(encode)
}
