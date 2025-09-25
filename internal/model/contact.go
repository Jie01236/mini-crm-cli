package model

type Contact struct {
	ID    uint   `json:"id" mapstructure:"id" gorm:"primaryKey"`
	Name  string `json:"name" mapstructure:"name"`
	Email string `json:"email" mapstructure:"email"`
}
