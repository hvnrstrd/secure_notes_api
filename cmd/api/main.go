package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hvnrstrd/secure_notes_api/internal/handler"
	"github.com/hvnrstrd/secure_notes_api/internal/storage"
)

func main() {
	s := storage.New()
	h := handler.New(s)

	http.Handle("/notes", h)
	http.Handle("/notes/", h)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("server started on :8080")
	log.Fatal(srv.ListenAndServe())
}