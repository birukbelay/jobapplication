package models

import "time"

// ProfileUpdateDto This is Used for creating and updating the user
type ProfileUpdateDto struct {
	FirstName string `json:"first_name,omitempty" minLength:"1"`
	LastName  string `json:"last_name,omitempty" `
	Email     string `json:"email,omitempty" format:"email"`
	Avatar    string `json:"avatar,omitempty" `
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
	UserId    string      `gorm:"uniqueIndex;not null"`
	CodeHash  string      `gorm:"not null"`
	Purpose   CodePurpose `gorm:"not null"`
	ExpiresAt *time.Time  `json:"-" `
}

//=================================   !  CompanyStatus  ============================

type CodePurpose string

const (
	Verification  = CodePurpose("verification")
	PasswordReset = CodePurpose("password_reset")
	TWOFA         = CodePurpose("2fa")
)
