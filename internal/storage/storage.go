package storage

import "github.com/hvnrstrd/secure_notes_api/internal/model"

type Storage interface {
	Create(userID, title, body string) (model.Note, error)
	GetAll(userID string) ([]model.Note, error)
	GetByID(userID, id string) (model.Note, error)
	Delete(userID, id string) error
	CreateUser(email, password string) (model.User, error)
	GetUserByEmail(email string) (model.User, string, error)
}