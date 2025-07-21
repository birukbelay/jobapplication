package models

import (
	Imdl "github.com/birukbelay/gocmn/src/dtos"

	"github.com/projTemplate/goauth/src/models/enums"
)

type Admin struct {
	Base `mapstructure:",squash" `

	UserDto       `mapstructure:",squash" `
	CompanyStatus string  `json:"company_status,omitempty" ` //approved, info filled,
	Company       Company `gorm:"foreignKey:CompanyID"`
}

type User struct {
	Base `mapstructure:",squash" `

	UserDto `mapstructure:",squash" `

	Company Company `json:"member_company,omitempty" gorm:"foreignKey:CompanyID"`
}

type UserDto struct {
	FirstName string  `json:"firsName,omitempty" minLength:"1"`
	LastName  string  `json:"lastName,omitempty" `
	Email     *string `json:"email,omitempty" format:"email"  gorm:"uniqueIndex" `
	Username  string  `json:"username,omitempty"`
	Avatar    string  `json:"avatar,omitempty" `
	CompanyID *string `json:"company_id"`
	Active    bool    `json:"active,omitempty"`
	//
	Password      string     `json:"-" `
	Role          enums.Role `gorm:"default:OWNER" json:"role,omitempty" `
	AccountStatus enums.AccountStatus
}
type UserFilter struct {
	ID            string     `query:"id,omitempty"`
	FName         string     `query:"fName,omitempty"`
	LName         string     `query:"lName,omitempty" `
	Email         string     `query:"email,omitempty"`
	Role          enums.Role `query:"role,omitempty"`
	Username      string     `query:"username,omitempty"`
	AccountStatus enums.AccountStatus
}

type UserQuery struct {
	Imdl.PaginationInput `mapstructure:",squash"`
	SelectedFields       []string `query:"selected_fields" enum:"first_name,last_name,email,avatar,company_id,account_status,id,created_at,updated_at"`
	Sort                 string   `query:"sort" enum:"first_name,last_name,email,avatar,company_id,account_status,created_at,updated_at"`
}
