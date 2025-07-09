package models

import (
	Imdl "github.com/birukbelay/gocmn/src/dtos"

	"github.com/projTemplate/goauth/src/models/enums"
)

type User struct {
	Base `mapstructure:",squash" `

	UserDto  `mapstructure:",squash" `
	Password string `json:"-" `

	Role          enums.Role `gorm:"default:OWNER" json:"role,omitempty" `
	HashedRefresh string     `json:"-" `
	//below Fields are Used For Verification of users by code
	VerificationCodeHash   string `json:"-" `
	VerificationCodeExpire int64  `json:"-" `

	//the Company A user owns, one user can own 1 company only

	//relation for employee of the companies

}

type UserDto struct {
	FirstName string `json:"firsName,omitempty" minLength:"1"`
	LastName  string `json:"lastName,omitempty" `
	Email     string `json:"email,omitempty" format:"email"`
	Username  string `json:"username,omitempty"`
	Avatar    string `json:"avatar,omitempty" `
	// AccountStatus enums.AccountStatus `json:"account_status,omitempty" enum:"PendingVerification,Verified, "`
	// CompanyID     *string             `json:"company_id"`
	Active bool `json:"active,omitempty"`
}
type UserFilter struct {
	ID       string     `query:"id,omitempty"`
	FName    string     `query:"fName,omitempty"`
	LName    string     `query:"lName,omitempty" `
	Email    string     `query:"email,omitempty"`
	Role     enums.Role `query:"role,omitempty"`
	Username string     `query:"username,omitempty"`
}

type UserQuery struct {
	Imdl.PaginationInput `mapstructure:",squash"`
	SelectedFields []string `query:"selected_fields" enum:"f_name,l_name,email,avatar,company_id,account_status,id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"f_name,l_name,email,avatar,company_id,account_status,created_at,updated_at"`
}
