package models

import "github.com/projTemplate/goauth/src/models/enums"

type Company struct {
	Base       `mapstructure:",squash" `
	CompanyDto `mapstructure:",squash" `

	Owner     *Admin  `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Employees []User `gorm:"foreignKey:CompanyID"`
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
type UpdateCompanyDto struct {
	Name     string `json:"name"`
	About    string `json:"about,omitempty"`
	Location string `json:"location,omitempty"`
}

type CompanyUpdateDto struct {
	Name     string `json:"name"`
	About    string `json:"about,omitempty"`
	Location string `json:"location,omitempty"`
}
type CompanyFilter struct {
	ID      string
	Name    string `query:"name" `
	OwnerID string `query:"owner_id,omitempty"`
	Handle  string `json:"handle" `

	ItemCount int `query:"item_count"`
}
type CompanyQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"name,about,location,owner_id, id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"name,created_at,updated_at"`
}
