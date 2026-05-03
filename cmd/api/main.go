package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hvnrstrd/secure_notes_api/internal/handler"
	"github.com/hvnrstrd/secure_notes_api/internal/storage"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := storage.NewPostgres(connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	h := handler.New(db)

	http.Handle("/notes", h)
	http.Handle("/notes/", h)

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("server started on :8080")
	log.Fatal(srv.ListenAndServe())
}