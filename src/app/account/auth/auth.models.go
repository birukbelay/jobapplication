package auth

import (
	IUsr "github.com/projTemplate/goauth/src/models"
)

// RegisterClientInput This is Used for creating and updating the user
type RegisterClientInput struct {
	FName    string `json:"fName" binding:"required,min=2" `
	LName    string `json:"lName" `
	Email    string `json:"email" binding:"required,email" gorm:"unique"`
	Password string `json:"password" binding:"required,min=6"`
	Avatar   string `json:"avatar" `
}

type VerifyUserInput struct {
	Email string `json:"email" binding:"required,email" gorm:"unique"`
}

// VerificationInput is for sending verification code
type VerificationInput struct {
	Info string `json:"info" validate:"required"`
	Code string ` json:"code" validate:"required" `
}

// LoginData defines user login form struct
type LoginData struct {
	LoginInfo string ` json:"info" binding:"required"`
	Password  string ` json:"password" binding:"required"`
	InfoType  string `json:"info_type,omitempty" `
}

type ChangePwdInput struct {
	OldPwd string `json:"old_pwd"`
	NewPwd string `json:"new_pwd"`
}
type RefreshTokenInput struct {
	Token string `json:"token"`
}

// PwdResetInput defines user password reset form struct
type PwdResetInput struct {
	NewPassword string `json:"new_password"  binding:"required"`
	Code        string `json:"code"  binding:"required"`
	Info        string `json:"info"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type TokenResponse struct {
	AuthTokens AuthTokens `json:"auth_tokens"`
	UserData   IUsr.User  `json:"user_data"`
}
