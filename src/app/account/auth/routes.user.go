package auth

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
)

const (
	RegisterUser     = constant.OperationId("Au_1-Register")
	VerifyUser       = constant.OperationId("Au_2-Verify")
	LoginUser        = constant.OperationId("Au_3-Login")
	RefreshTokenUser = constant.OperationId("Au_4-RefreshToken")
	LogoutUser       = constant.OperationId("Au_5-Logout")
	ForgotPwdUser    = constant.OperationId("Au_6-ForgotPwd")
	ResetPwdUser     = constant.OperationId("Au_7-ResetPwd")
)

func SetupUserAuthRoutes(humaRouter huma.API, providerS *providers.IProviderS, serv *Service[models.User]) {
	handler := NewAuthHandler(providerS, serv)
	tags := []string{"user_auth"}
	path := constant.ApiV1 + "/auth"

	huma.Register(humaRouter, huma.Operation{
		OperationID: RegisterUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/signup",
		Tags:        tags}, handler.Register,
	)

	huma.Register(humaRouter, huma.Operation{
		OperationID: VerifyUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/verify",
		Tags:        tags}, handler.VerifyAccount,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: LoginUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/login",
		Tags:        tags}, handler.Login,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: RefreshTokenUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/refresh",
		Tags:        tags}, handler.RefreshToken,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: LogoutUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/logout",
		Tags:        tags}, handler.Logout,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: ForgotPwdUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/forgot_password",
		Tags:        tags}, handler.ForgotPwd,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: ResetPwdUser.Str(),
		Method:      http.MethodPost,
		Path:        path + "/reset_password",
		Tags:        tags}, handler.ResetPwd,
	)

}
