package auth

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
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

const (
	RegisterOwner = constant.OperationId("Ad_Au_1-RegisterOwner")
	VerifyOwner   = constant.OperationId("Ad_Au_2-VerifyOwner")
	Login         = constant.OperationId("Ad_Au_3-AdminLogin")
	RefreshToken  = constant.OperationId("Ad_Au_4-AdminRefreshToken")
	Logout        = constant.OperationId("Ad_Au_5-AdminLogout")
	ForgotPwd     = constant.OperationId("Ad_Au_6-AdminForgotPwd")
	ResetPwd      = constant.OperationId("Ad_Au_7-AdminResetPwd")
)

func SetupAdminAuthRoutes(humaRouter huma.API, providerS *providers.IProviderS, serv *Service[models.Admin]) {
	handler := NewAuthHandler(providerS, serv)
	tags := []string{"admin_auth"}
	path := constant.ApiV1 + "/admin_auth"

	huma.Register(humaRouter, huma.Operation{
		OperationID: RegisterOwner.Str(),
		Method:      http.MethodPost,
		Path:        path + "/signup",
		Tags:        tags}, handler.Register,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: VerifyOwner.Str(),
		Method:      http.MethodPost,
		Path:        path + "/verify_owner",
		Tags:        tags}, handler.VerifyAccount,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: Login.Str(),
		Method:      http.MethodPost,
		Path:        path + "/login",
		Tags:        tags}, handler.Login,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: RefreshToken.Str(),
		Method:      http.MethodPost,
		Path:        path + "/refresh",
		Tags:        tags}, handler.RefreshToken,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: Logout.Str(),
		Method:      http.MethodPost,
		Path:        path + "/logout",
		Tags:        tags}, handler.Logout,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: ForgotPwd.Str(),
		Method:      http.MethodPost,
		Path:        path + "/forgot_password",
		Tags:        tags}, handler.ForgotPwd,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: ResetPwd.Str(),
		Method:      http.MethodPost,
		Path:        path + "/reset_password",
		Tags:        tags}, handler.ResetPwd,
	)

}
