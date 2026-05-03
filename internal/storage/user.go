package storage

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hvnrstrd/secure_notes_api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

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