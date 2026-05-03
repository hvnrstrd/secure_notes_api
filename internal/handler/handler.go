package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hvnrstrd/secure_notes_api/internal/storage"
)

type Handler struct {
	storage *storage.Storage
}

func New(s *storage.Storage) *Handler {
	return &Handler{storage: s}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/notes")
	path = strings.TrimPrefix(path, "/")

	switch {
	case r.Method == http.MethodGet && path == "":
		h.getAll(w, r)
	case r.Method == http.MethodPost && path == "":
		h.create(w, r)
	case r.Method == http.MethodGet && path != "":
		h.getByID(w, r, path)
	case r.Method == http.MethodDelete && path != "":
		h.delete(w, r, path)
	default:
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	}
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	notes := h.storage.GetAll()
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}
	note := h.storage.Create(body.Title, body.Body)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	note, err := h.storage.GetByID(id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.storage.Delete(id); err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}