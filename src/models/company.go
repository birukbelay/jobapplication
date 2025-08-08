package models

import (
	"time"

	"github.com/projTemplate/goauth/src/models/enums"
)

type Company struct {
	Base       `mapstructure:",squash" `
	CompanyDto `mapstructure:",squash" `
	Employees  []User `gorm:"foreignKey:CompanyID"`
	Owner      Admin  `gorm:"foreignKey:OwnerID"`
}
type CompanyDto struct {
	Name          string              `json:"name"`
	Handle        string              `json:"handle" gorm:"unique"`
	About         string              `json:"about,omitempty"`
	Location      string              `json:"location,omitempty"`
	CompanyStatus enums.CompanyStatus `json:"company_status,omitempty" enum:"approved,pending_approval,deleted"`

	//Relationships
	OwnerID string `json:"owner_id,omitempty"`
}
type CompanyUpdateDto struct {
	Handle   string `json:"handle" `
	Name     string `json:"name"`
	About    string `json:"about,omitempty"`
	Location string `json:"location,omitempty"`
}
type CompanyFilter struct {
	ID      string `query:"id"`
	Name    string `query:"name" `
	OwnerID string `query:"owner_id,omitempty"`
	Handle  string `query:"handle" `
}
type CompanyQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"name,about,location,owner_id, id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"name,created_at,updated_at"`
}

func (q CompanyQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}

// ==============. Invite Codes

type InviteCode struct {
	Base          `mapstructure:",squash" `
	InviteCodeDto `mapstructure:",squash" `

	Company Company `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	// Users   []User  `gorm:"foreignKey:CompanyID"`
	UsageCount int `json:"usate_count,omitempty" `
}
type InviteCodeDto struct {
	Name      string     `json:"name,omitempty" `
	CompanyID string     `json:"company_id,omitempty" gorm:"not null;index"`
	Code      string     `json:"code,omitempty" gorm:"uniqueindex"`
	UserRole  string     `json:"user_role,omitempty" `
	ExpiresAt *time.Time `json:"expires_at,omitempty" `

	MaxUsage *int    `json:"max_usage,omitempty" `
	UserInfo *string `json:"user_info,omitempty" `
	Active   bool    `json:"active,omitempty" `
}

func (d InviteCodeDto) SetOnCreate(key string) {
	d.CompanyID = key
}

type InviteCodeUpdateDto struct {
	Name      string     `json:"name,omitempty" `
	ExpiresAt *time.Time `json:"expires_at,omitempty" `

	MaxUsage *int `json:"max_usage,omitempty" `
	Active   bool `json:"active,omitempty" `
}
type InviteCodeFilter struct {
	ID         string `query:"id"`
	Name       string `query:"name" `
	CompanyID  string `query:"company_id" `
	Code       string `query:"code" `
	UserRole   string `query:"code" `
	UsageCount string `query:"code" `
	MaxUsage   int    `query:"max_usage" `
	UserInfo   string `query:"user_info" `
	Active     bool   `query:"active" `
}
type InviteCodeQuery struct {
	ExpiresAt      time.Time `query:"expires_at" `
	SelectedFields []string  `query:"selected_fields" enum:"company_id,expires_at,user_role,usage_count,max_usage,user_info, id,created_at,updated_at"`
	Sort           string    `query:"sort" enum:"company_id,usage_count,max_usage,expires_at,created_at,updated_at,id"`
}

func (q InviteCodeQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}

type JoinViaCode struct {
	Code string `json:"code" gorm:"uniqueindex"`
}
