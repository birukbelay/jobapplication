package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/models/enums"
)

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

type IntUsr interface {
	GetID() string
	GetRole() string
	GetPwd() string
	GetStatus() enums.AccountStatus
	GetCompanyId() string
	GetInfo() string
}

func (b Base) GetID() string {
	return b.ID
}
func GetID[T IntUsr](t T) string {
	return t.GetID()
}
func (u UserDto) GetRole() string {
	return string(u.Role)
}
func (u UserDto) GetPwd() string {
	return u.Password
}
func (u UserDto) GetInfo() string {
	return *u.Email
}
func (u UserDto) GetStatus() enums.AccountStatus {
	return u.AccountStatus
}

func (u UserDto) GetCompanyId() string {
	if u.CompanyID != nil {
		return *u.CompanyID
	}
	return ""
}
