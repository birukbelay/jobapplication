package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

func Ptr[T any](value T) *T {
	return &value
}

type Base struct {
	ID        string     `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}
	return nil
}