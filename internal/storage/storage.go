package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/hvnrstrd/secure_notes_api/internal/model"
	"github.com/google/uuid"
)

type Storage struct {
	mu    sync.RWMutex
	notes map[string]model.Note
}

func New() *Storage {
	return &Storage{
		notes: make(map[string]model.Note),
	}
}

func (s *Storage) Create(title, body string) model.Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	note := model.Note{
		ID:        uuid.New().String(),
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}
	s.notes[note.ID] = note
	return note
}

func (s *Storage) GetAll() []model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notes := make([]model.Note, 0, len(s.notes))
	for _, n := range s.notes {
		notes = append(notes, n)
	}
	return notes
}

func (s *Storage) GetByID(id string) (model.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	note, ok := s.notes[id]
	if !ok {
		return model.Note{}, errors.New("note not found")
	}
	return note, nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.notes[id]; !ok {
		return errors.New("note not found")
	}
	delete(s.notes, id)
	return nil
}