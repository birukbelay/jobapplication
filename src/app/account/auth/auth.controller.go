package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/birukbelay/gocmn/src/dtos"

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

func (ah *GinAuthHandler) Login(ctx context.Context, inputs *dtos.HumaReqBody[LoginData]) (*dtos.HumaResponse[dtos.GResp[TokenResponse]], error) {
	tkn, err := ah.AuthServ.Login(ctx, inputs.Body)
	if err != nil {
		return dtos.HumaReturnG(tkn, err)
	}
	refreshCookie := CreateCookie(Icnst.RefreshToken, tkn.Body.AuthTokens.RefreshToken, ah.AuthServ.Config.JwtVar.RefreshExpireMin)
	accessCookie := CreateCookie(Icnst.AccessToken, tkn.Body.AuthTokens.AccessToken, ah.AuthServ.Config.JwtVar.AccessExpireMin)
	return dtos.HumaReturnGWithCookie(tkn, err, []http.Cookie{refreshCookie, accessCookie})
}

func (ah *GinAuthHandler) RefreshToken(ctx context.Context, inputs *dtos.HumaReqBody[RefreshTokenInput]) (*dtos.HumaResponse[dtos.GResp[TokenResponse]], error) {
	tkn, err := ah.AuthServ.ResetToken(ctx, inputs.Body.Token)
	refreshCookie := CreateCookie(Icnst.RefreshToken, tkn.Body.AuthTokens.RefreshToken, ah.AuthServ.Config.JwtVar.RefreshExpireMin)
	accessCookie := CreateCookie(Icnst.AccessToken, tkn.Body.AuthTokens.AccessToken, ah.AuthServ.Config.JwtVar.AccessExpireMin)
	return dtos.HumaReturnGWithCookie(tkn, err, []http.Cookie{accessCookie, refreshCookie})
}
func (ah *GinAuthHandler) RegisterOwner(ctx context.Context, inputs *dtos.HumaReqBody[RegisterClientInput]) (*dtos.HumaResponse[dtos.GResp[models.Admin]], error) {
	usr, err := ah.AuthServ.RegisterCompanyOwner(ctx, inputs.Body)
	return dtos.HumaReturnG(usr, err)
}
func (ah *GinAuthHandler) VerifyCompanyOwner(ctx context.Context, inputs *dtos.HumaReqBody[VerificationInput]) (*dtos.HumaResponse[dtos.GResp[bool]], error) {
	usr, err := ah.AuthServ.VerifyCompanyUser(ctx, inputs.Body)
	return dtos.HumaReturnG(usr, err)
}
