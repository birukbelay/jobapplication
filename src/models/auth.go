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
	Base             `mapstructure:",squash" `
	RefreshTokenHash string
	DeviceName       string
	Blacklisted      bool
	BlacklistedOn    *time.Time
	UserId           string
}

type VerificationCodes struct {
	Base     `mapstructure:",squash" `
	SendAt   time.Time
	CodeHash string
}
