package storage

import (
	"context"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mini-crm-cli/internal/model"
)

type gormStore struct {
	db *gorm.DB
}

func NewGORMStore(path string) (Store, error) {
	if path == "" {
		path = "contacts.db"
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Contact{}); err != nil {
		return nil, err
	}
	return &gormStore{db: db}, nil
}

func (g *gormStore) Create(ctx context.Context, contact *model.Contact) error {
	contact.ID = 0
	return g.db.WithContext(ctx).Create(contact).Error
}

func (g *gormStore) List(ctx context.Context) ([]model.Contact, error) {
	var contacts []model.Contact
	if err := g.db.WithContext(ctx).Order("id asc").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

func (g *gormStore) Get(ctx context.Context, id uint) (*model.Contact, error) {
	var contact model.Contact
	err := g.db.WithContext(ctx).First(&contact, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &contact, nil
}

func (g *gormStore) Update(ctx context.Context, contact *model.Contact) error {
	if contact.ID == 0 {
		return ErrNotFound
	}
	tx := g.db.WithContext(ctx).Model(&model.Contact{}).Where("id = ?", contact.ID).Updates(map[string]any{
		"name":  contact.Name,
		"email": contact.Email,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (g *gormStore) Delete(ctx context.Context, id uint) error {
	tx := g.db.WithContext(ctx).Delete(&model.Contact{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
