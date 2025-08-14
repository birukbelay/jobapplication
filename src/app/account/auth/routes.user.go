package auth

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
)

const (
	RegisterUser     = constant.OperationId("Au-1-Register")
	VerifyUser       = constant.OperationId("Au-2-Verify")
	LoginUser        = constant.OperationId("Au-3-Login")
	RefreshTokenUser = constant.OperationId("Au-4-RefreshToken")
	LogoutUser       = constant.OperationId("Au-5-Logout")
	ForgotPwdUser    = constant.OperationId("Au-6-ForgotPwd")
	ResetPwdUser     = constant.OperationId("Au-7-ResetPwd")
)

type GinAuthHandler[T models.IntUsr] struct {
	AdminAuthServ *Service[T]
	CmnServ       *providers.IProviderS
}

func NewAuthHandler[T models.IntUsr](cmnServ *providers.IProviderS, serv *Service[T]) *GinAuthHandler[T] {
	return &GinAuthHandler[T]{
		AdminAuthServ: serv,
		CmnServ:       cmnServ,
	}
}
func SetupUserAuthRoutes(humaRouter huma.API, providerS *providers.IProviderS, serv *Service[models.User]) {
	handler := NewAuthHandler(providerS, serv)
	tags := []string{"auth"}
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
