package storage

import (
	"database/sql"
	"time"

	"github.com/hvnrstrd/secure_notes_api/internal/model"
	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgres(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id UUID PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (s *PostgresStorage) Create(title, body string) (model.Note, error) {
	note := model.Note{
		ID:        uuid.New().String(),
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}
	_, err := s.db.Exec(
		`INSERT INTO notes (id, title, body, created_at) VALUES ($1, $2, $3, $4)`,
		note.ID, note.Title, note.Body, note.CreatedAt,
	)
	return note, err
}

func (s *PostgresStorage) GetAll() ([]model.Note, error) {
	rows, err := s.db.Query(`SELECT id, title, body, created_at FROM notes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		var n model.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (s *PostgresStorage) GetByID(id string) (model.Note, error) {
	var n model.Note
	err := s.db.QueryRow(
		`SELECT id, title, body, created_at FROM notes WHERE id = $1`, id,
	).Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt)
	if err == sql.ErrNoRows {
		return model.Note{}, err
	}
	return n, err
}

func (s *PostgresStorage) Delete(id string) error {
	result, err := s.db.Exec(`DELETE FROM notes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}