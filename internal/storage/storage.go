package storage

import (
	"context"
	"errors"
	"fmt"

	"mini-crm-cli/internal/model"
)

var ErrNotFound = errors.New("contact not found")

type Store interface {
	Create(ctx context.Context, contact *model.Contact) error
	List(ctx context.Context) ([]model.Contact, error)
	Get(ctx context.Context, id uint) (*model.Contact, error)
	Update(ctx context.Context, contact *model.Contact) error
	Delete(ctx context.Context, id uint) error
}

type Config struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

func New(cfg Config) (Store, error) {
	switch cfg.Type {
	case "memory":
		return NewMemoryStore(), nil
	case "json":
		if cfg.Path == "" {
			return nil, fmt.Errorf("json storage requires a path")
		}
		return NewJSONStore(cfg.Path)
	case "gorm":
		return NewGORMStore(cfg.Path)
	case "":
		return nil, fmt.Errorf("storage type must be provided")
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
