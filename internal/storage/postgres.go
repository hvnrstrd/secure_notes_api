package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hvnrstrd/secure_notes_api/internal/model"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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
			user_id UUID NOT NULL,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (s *PostgresStorage) MigrateUsers() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (s *PostgresStorage) Create(userID, title, body string) (model.Note, error) {
	note := model.Note{
		ID:        uuid.New().String(),
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}
	_, err := s.db.Exec(
		`INSERT INTO notes (id, user_id, title, body, created_at) VALUES ($1, $2, $3, $4, $5)`,
		note.ID, userID, note.Title, note.Body, note.CreatedAt,
	)
	return note, err
}

func (s *PostgresStorage) GetAll(userID string) ([]model.Note, error) {
	rows, err := s.db.Query(
		`SELECT id, title, body, created_at FROM notes WHERE user_id = $1`, userID,
	)
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

func (s *PostgresStorage) GetByID(userID, id string) (model.Note, error) {
	var n model.Note
	err := s.db.QueryRow(
		`SELECT id, title, body, created_at FROM notes WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt)
	if err == sql.ErrNoRows {
		return model.Note{}, err
	}
	return n, err
}

func (s *PostgresStorage) Delete(userID, id string) error {
	result, err := s.db.Exec(
		`DELETE FROM notes WHERE id = $1 AND user_id = $2`, id, userID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *PostgresStorage) CreateUser(email, password string) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user := model.User{
		ID:        uuid.New().String(),
		Email:     email,
		CreatedAt: time.Now(),
	}
	_, err = s.db.Exec(
		`INSERT INTO users (id, email, password, created_at) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Email, string(hash), user.CreatedAt,
	)
	return user, err
}

func (s *PostgresStorage) GetUserByEmail(email string) (model.User, string, error) {
	var user model.User
	var hash string
	err := s.db.QueryRow(
		`SELECT id, email, password, created_at FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &hash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return model.User{}, "", err
	}
	return user, hash, err
}