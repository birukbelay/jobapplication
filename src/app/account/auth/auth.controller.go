package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/birukbelay/gocmn/src/dtos"
	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/mitchellh/mapstructure"

	Icnst "github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models"
)

func CreateCookie(Name, Value string, minutes int) http.Cookie {
	return http.Cookie{
		Name:     Name,
		Value:    Value,
		Expires:  time.Now().Add(time.Minute * time.Duration(minutes)),
		HttpOnly: true, // Set HttpOnly to true to make the cookie accessible only through HTTP requests, not JavaScript
		Path:     "/",
		Domain:   "",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func (ah *GinAuthHandler) Login(ctx context.Context, inputs *dtos.HumaReqWithBody[LoginData]) (*dtos.HumaResponse[dtos.GResp[TokenResponse]], error) {
	tkn, err := ah.AuthServ.Login(ctx, inputs.Body)
	if err != nil {
		return dtos.HumaGError(tkn, err)
	}
	refreshCookie := CreateCookie(Icnst.RefreshToken, tkn.Body.AuthTokens.RefreshToken, ah.AuthServ.Config.JwtVar.RefreshExpireMin)
	accessCookie := CreateCookie(Icnst.AccessToken, tkn.Body.AuthTokens.AccessToken, ah.AuthServ.Config.JwtVar.AccessExpireMin)
	return dtos.ReturnWithCookie(tkn, err, []http.Cookie{refreshCookie, accessCookie})
}

func (ah *GinAuthHandler) RefreshToken(ctx context.Context, inputs *dtos.HumaReqBody[RefreshTokenInput]) (*dtos.HumaResponse[dtos.GResp[TokenResponse]], error) {
	tkn, err := ah.AuthServ.ResetToken(ctx, inputs.Body.Token)
	refreshCookie := CreateCookie(Icnst.RefreshToken, tkn.Body.AuthTokens.RefreshToken, ah.AuthServ.Config.JwtVar.RefreshExpireMin)
	accessCookie := CreateCookie(Icnst.AccessToken, tkn.Body.AuthTokens.AccessToken, ah.AuthServ.Config.JwtVar.AccessExpireMin)
	return dtos.ReturnWithCookie(tkn, err, []http.Cookie{accessCookie, refreshCookie})
}

// RegisterCompanyOwner (acc-01)
func (aus Service) RegisterCompanyOwner(ctx context.Context, input RegisterClientInput) (*dtos.GResp[models.User], error) {
	//Check if the email already exists
	usr, err := sql_db.GetOne[models.User](aus.Provider, ctx, models.UserFilter{Email: input.Email}, nil)
	//if the user already exists
	if err == nil {
		/*if the user is active
		if usr.Body.Status==StatusActive
		*/
		//If the User is not active
		if usr.RowsAffected > 0 {
			//FIXME Send password or email wrong depending on the scenario
			return dtos.GBadReqWithCode[models.User](ICnst.UserExists), ICnst.UserExistError
		}

	}
	//TODO: verify the Email

	var userModel models.User
	if err := mapstructure.Decode(input, &userModel); err != nil {
		return dtos.BadRequestMsg[models.User]("Decoding Input Error"), err
	}

	hash, err := ICrypt.BcryptCreateHash(input.Password)
	if err != nil {
		return dtos.InternalErr[models.User]("Hashing Error"), err
	}
	userModel.Password = hash
	userModel.Role = enums.OWNER
	userModel.AccountStatus = enums.AccountPendingVerification
	userModel.Active = false

	user, err := sql_db.CreateOne[models.User](aus.Provider, ctx, &userModel, nil)
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErr[models.User]("creating Error"), err
	}
	return dtos.Success[models.User](user.Body, user.RowsAffected), nil
}
