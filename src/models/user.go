package models

import (
	"time"

	Imdl "github.com/birukbelay/gocmn/src/dtos"

	"github.com/projTemplate/goauth/src/models/enums"
)

type Admin struct {
	Base `mapstructure:",squash" `

	UserDto       `mapstructure:",squash" `
	Password      string     `json:"-" `
	Role          enums.Role `gorm:"default:OWNER" json:"role,omitempty" `
	AccountStatus enums.AccountStatus
	HashedRefresh string `json:"-" `
	//below Fields are Used For Verification of users by code
	VerificationCodeHash   string     `json:"-" `
	VerificationCodeExpire *time.Time `json:"-" `
	Company                Company    `gorm:"foreignKey:OwnerID"`
}

type User struct {
	Base `mapstructure:",squash" `

	UserDto       `mapstructure:",squash" `
	Password      string     `json:"-" `
	Role          enums.Role `gorm:"default:OWNER" json:"role,omitempty" `
	HashedRefresh string     `json:"-" `
	//below Fields are Used For Verification of users by code
	VerificationCodeHash   string     `json:"-" `
	VerificationCodeExpire *time.Time `json:"-" `
	Company                Company    `json:"member_company,omitempty" gorm:"foreignKey:CompanyID"`
}

type UserDto struct {
	FirstName string  `json:"firsName,omitempty" minLength:"1"`
	LastName  string  `json:"lastName,omitempty" `
	Email     *string `json:"email,omitempty" format:"email"  gorm:"uniqueIndex" `
	Username  string  `json:"username,omitempty"`
	Avatar    string  `json:"avatar,omitempty" `
	CompanyID string  `json:"company_id"`
	Active    bool    `json:"active,omitempty"`
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
	SelectedFields       []string `query:"selected_fields" enum:"f_name,l_name,email,avatar,company_id,account_status,id,created_at,updated_at"`
	Sort                 string   `query:"sort" enum:"f_name,l_name,email,avatar,company_id,account_status,created_at,updated_at"`
}
