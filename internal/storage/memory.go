package storage

import (
	"context"
	"sort"
	"sync"

	"mini-crm-cli/internal/model"
)

type memoryStore struct {
	mu       sync.RWMutex
	contacts map[uint]*model.Contact
	nextID   uint
}

func NewMemoryStore() Store {
	return &memoryStore{
		contacts: make(map[uint]*model.Contact),
		nextID:   1,
	}
}

func (m *memoryStore) Create(_ context.Context, contact *model.Contact) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	contactCopy := *contact
	contactCopy.ID = m.nextID
	m.nextID++

	m.contacts[contactCopy.ID] = &contactCopy
	contact.ID = contactCopy.ID
	return nil
}

func (m *memoryStore) List(_ context.Context) ([]model.Contact, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.contacts) == 0 {
		return []model.Contact{}, nil
	}

	ids := make([]uint, 0, len(m.contacts))
	for id := range m.contacts {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	result := make([]model.Contact, 0, len(ids))
	for _, id := range ids {
		c := m.contacts[id]
		result = append(result, *c)
	}

	return result, nil
}

func (m *memoryStore) Get(_ context.Context, id uint) (*model.Contact, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contact, ok := m.contacts[id]
	if !ok {
		return nil, ErrNotFound
	}
	copy := *contact
	return &copy, nil
}

func (m *memoryStore) Update(_ context.Context, contact *model.Contact) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.contacts[contact.ID]; !ok {
		return ErrNotFound
	}
	copy := *contact
	m.contacts[contact.ID] = &copy
	return nil
}

func (m *memoryStore) Delete(_ context.Context, id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.contacts[id]; !ok {
		return ErrNotFound
	}
	delete(m.contacts, id)
	return nil
}
