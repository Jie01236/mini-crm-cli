package storage

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sort"
	"sync"

	"mini-crm-cli/internal/model"
)

type jsonStore struct {
	path     string
	mu       sync.RWMutex
	contacts map[uint]*model.Contact
	nextID   uint
}

func NewJSONStore(path string) (Store, error) {
	store := &jsonStore{
		path:     path,
		contacts: make(map[uint]*model.Contact),
		nextID:   1,
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *jsonStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var contacts []model.Contact
	if err := json.Unmarshal(data, &contacts); err != nil {
		return err
	}

	var maxID uint
	for i := range contacts {
		s.copyIn(&contacts[i])
		if contacts[i].ID > maxID {
			maxID = contacts[i].ID
		}
	}
	s.nextID = maxID + 1
	if s.nextID == 0 {
		s.nextID = 1
	}
	return nil
}

func (s *jsonStore) copyIn(contact *model.Contact) {
	c := *contact
	s.contacts[c.ID] = &c
}

func (s *jsonStore) persistLocked() error {
	contacts := make([]model.Contact, 0, len(s.contacts))
	for _, contact := range s.contacts {
		contacts = append(contacts, *contact)
	}
	sort.Slice(contacts, func(i, j int) bool { return contacts[i].ID < contacts[j].ID })

	data, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func (s *jsonStore) Create(_ context.Context, contact *model.Contact) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := *contact
	copy.ID = s.nextID
	s.nextID++

	s.contacts[copy.ID] = &copy
	contact.ID = copy.ID

	return s.persistLocked()
}

func (s *jsonStore) List(_ context.Context) ([]model.Contact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	contacts := make([]model.Contact, 0, len(s.contacts))
	for _, contact := range s.contacts {
		contacts = append(contacts, *contact)
	}
	sort.Slice(contacts, func(i, j int) bool { return contacts[i].ID < contacts[j].ID })
	return contacts, nil
}

func (s *jsonStore) Get(_ context.Context, id uint) (*model.Contact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	contact, ok := s.contacts[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *contact
	return &copy, nil
}

func (s *jsonStore) Update(_ context.Context, contact *model.Contact) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.contacts[contact.ID]; !ok {
		return ErrNotFound
	}
	copy := *contact
	s.contacts[contact.ID] = &copy

	return s.persistLocked()
}

func (s *jsonStore) Delete(_ context.Context, id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.contacts[id]; !ok {
		return ErrNotFound
	}
	delete(s.contacts, id)

	return s.persistLocked()
}
