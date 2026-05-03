package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hvnrstrd/secure_notes_api/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func New(s storage.Storage) *Handler {
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
	notes, err := h.storage.GetAll()
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, `{"error":"encoding failed"}`, http.StatusInternalServerError)
	}
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
	note, err := h.storage.Create(body.Title, body.Body)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, `{"error":"encoding failed"}`, http.StatusInternalServerError)
	}
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request, id string) {
	note, err := h.storage.GetByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, `{"error":"encoding failed"}`, http.StatusInternalServerError)
	}
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.storage.Delete(id); err == sql.ErrNoRows {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}