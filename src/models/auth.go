package models

import "time"

type RegisterClientInput struct {
	FirstName string `json:"fName" binding:"required,min=2" `
	LastName  string `json:"lName" `
	Email     string `json:"email" binding:"required,email" gorm:"unique"`
	Password  string `json:"password" binding:"required,min=6"`
	Avatar    string `json:"avatar" `
}

// ProfileUpdateDto This is Used for creating and updating the user
type ProfileUpdateDto struct {
	FirstName string `json:"first_name,omitempty" minLength:"1"`
	LastName  string `json:"last_name,omitempty" `
	// Email     string `json:"email,omitempty" format:"email"`
	Avatar string `json:"avatar,omitempty" `
}

// PasswordUpdateDto This is Used for creating and updating the user
type PasswordUpdateDto struct {
	OldPassword string `json:"old_password,omitempty" minLength:"6"`
	NewPassword string `json:"new_password,omitempty" minLength:"6"`
}

// session will be put on redis, make it polymorphism for admin and users
type Session struct {
	Base          `mapstructure:",squash" `
	SessionId     string `gorm:"uniqueIndex;not null"`
	UserId        string `gorm:"not null"`
	HashedRefresh string `gorm:"not null"`
	DeviceInfo    string
	// Blacklisted      bool
	// BlacklistedOn    *time.Time
}

type VerificationCode struct {
	Base      `mapstructure:",squash" `
	UserId    string `gorm:"uniqueIndex;not null"`
	Email     string
	CodeHash  string      `gorm:"not null"`
	Purpose   CodePurpose `gorm:"not null"`
	ExpiresAt *time.Time  `json:"-" `
}

//=================================   !  CompanyStatus  ============================

type CodePurpose string

const (
	SignupVerification = CodePurpose("SIGNUP_VERIFICATION")
	PasswordReset      = CodePurpose("PASSWORD_RESET")
	ChangeEmail        = CodePurpose("CHANGE_EMAIL")
	TWOFA              = CodePurpose("2FA")
)

type ChangeEmailReqDto struct {
	Password string `json:"password" minLength:"6"`
	NewEmail string `json:"new_email,omitempty" minLength:"6"`
}

type VerifyEmailDto struct {
	Code string `json:"code" `
}
