package models

import (
	Imdl "github.com/birukbelay/gocmn/src/dtos"

	"github.com/projTemplate/goauth/src/models/enums"
)

type User struct {
	Base    `mapstructure:",squash" `
	UserDto `mapstructure:",squash" `
	// Company *Job `json:"member_company,omitempty" gorm:"foreignKey:CompanyID"`
}

type UserDto struct {
	FirstName string  `json:"first_name,omitempty" minLength:"1"`
	LastName  string  `json:"last_name,omitempty" `
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

func (d UserDto) SetOnCreate(key string) {
	d.CompanyID = &key
}

// cant update the email, accountstatus
type UserUpdateDto struct {
	FirstName string `json:"firsName,omitempty" minLength:"1"`
	LastName  string `json:"lastName,omitempty" `
	Username  string `json:"username,omitempty"`
	Avatar    string `json:"avatar,omitempty" `
	Active    bool   `json:"active,omitempty"`
	//
	Role          enums.Role          ` json:"role,omitempty" `                                 //TODO: what are roles related to users
	AccountStatus enums.AccountStatus `json:"account_status,omitempty" enum:"active, disabled"` //will only be allowed to disable and enable
}
type UserFilter struct {
	ID            string     `query:"id,omitempty"`
	FName         string     `query:"fName,omitempty"`
	LName         string     `query:"lName,omitempty" `
	Email         string     `query:"email,omitempty"`
	Role          enums.Role `query:"role,omitempty"`
	Username      string     `query:"username,omitempty"`
	AccountStatus enums.AccountStatus
	CompanyID     string `json:"company_id"`
}

type UserQuery struct {
	Imdl.PaginationInput `mapstructure:",squash"`
	SelectedFields       []string `query:"selected_fields" enum:"first_name,last_name,email,avatar,company_id,account_status,id,created_at,updated_at"`
	Sort                 string   `query:"sort" enum:"first_name,last_name,email,avatar,company_id,account_status,created_at,updated_at"`
}

func (q UserQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}
