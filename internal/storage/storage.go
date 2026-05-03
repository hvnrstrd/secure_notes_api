package storage

import "github.com/hvnrstrd/secure_notes_api/internal/model"

type Storage interface {
	Create(title, body string) (model.Note, error)
	GetAll() ([]model.Note, error)
	GetByID(id string) (model.Note, error)
	Delete(id string) error
}