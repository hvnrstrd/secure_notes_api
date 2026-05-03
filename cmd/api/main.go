package main

import (
	"log"
	"net/http"

	"github.com/hvnrstrd/secure_notes_api/internal/handler"
	"github.com/hvnrstrd/secure_notes_api/internal/storage"
)

func main() {
	s := storage.New()
	h := handler.New(s)

	http.Handle("/notes", h)
	http.Handle("/notes/", h)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}