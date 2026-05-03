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
		log.Fatalf("failed to migrate notes: %v", err)
	}

	if err := db.MigrateUsers(); err != nil {
		log.Fatalf("failed to migrate users: %v", err)
	}

	h := handler.New(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", h.Register)
	mux.HandleFunc("/auth/login", h.Login)
	mux.Handle("/notes", h.AuthMiddleware(h.ServeHTTP))
	mux.Handle("/notes/", h.AuthMiddleware(h.ServeHTTP))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler.CORSMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("server started on :8080")
	log.Fatal(srv.ListenAndServe())
}