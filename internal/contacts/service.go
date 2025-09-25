package contacts

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"mini-crm-cli/internal/model"
	"mini-crm-cli/internal/storage"
)

var (
	emailRegexp = regexp.MustCompile(`^[\w.+-]+@[\w.-]+\.[A-Za-z]{2,}$`)

	ErrInvalidName = errors.New("name cannot be empty")

	ErrInvalidEmail = errors.New("invalid email format")
)

type Service struct {
	store storage.Store
}

func NewService(store storage.Store) *Service {
	return &Service{store: store}
}

func (s *Service) AddContact(ctx context.Context, name, email string) (*model.Contact, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" {
		return nil, ErrInvalidName
	}
	if !emailRegexp.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	contact := &model.Contact{Name: name, Email: email}
	if err := s.store.Create(ctx, contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *Service) ListContacts(ctx context.Context) ([]model.Contact, error) {
	return s.store.List(ctx)
}

func (s *Service) UpdateContact(ctx context.Context, id uint, name, email string) (*model.Contact, error) {
	contact, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if trimmed := strings.TrimSpace(name); trimmed != "" {
		contact.Name = trimmed
	}
	if trimmed := strings.TrimSpace(email); trimmed != "" {
		if !emailRegexp.MatchString(trimmed) {
			return nil, ErrInvalidEmail
		}
		contact.Email = trimmed
	}

	if contact.Name == "" {
		return nil, ErrInvalidName
	}

	if err := s.store.Update(ctx, contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *Service) DeleteContact(ctx context.Context, id uint) error {
	return s.store.Delete(ctx, id)
}
