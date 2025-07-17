package auth

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

type GinAuthHandler struct {
	AuthServ *Service
	CmnServ  *providers.IProviderS
}

func NewAuthHandler(cmnServ *providers.IProviderS, serv *Service) *GinAuthHandler {
	return &GinAuthHandler{
		AuthServ: serv,
		CmnServ:  cmnServ,
	}
}

const (
	RegisterOwner = constant.OperationId("Au-1_RegisterOwner")
	VerifyOwner   = constant.OperationId("Au-2_VerifyOwner")
	Login         = constant.OperationId("Au-3_Login")
	RefreshToken  = constant.OperationId("Au-4_RefreshToken")
	Logout        = constant.OperationId("Au-5_Logout")
	ForgotPwd     = constant.OperationId("Au-6_ForgotPwd")
	ResetPwd      = constant.OperationId("Au-7_ResetPwd")

	DeactivateMyCompany = constant.OperationId("Ow-4_DeactivateMyCompany")
	UpdateMyCompany     = constant.OperationId("Ow-3_UpdateMyCompany")
	ApproveCompany      = constant.OperationId("Ad-1_ApproveCompany")
)

var PermissionMap = map[constant.OperationId][]enums.Role{
	RegisterOwner:  {},
	VerifyOwner:    {},
	Login:          {},
	RefreshToken:   {},
	ApproveCompany: {enums.ADMIN},
}

func SetupAuthRoutes(humaRouter huma.API, providerS *providers.IProviderS, serv *Service) {
	handler := NewAuthHandler(providerS, serv)
	tags := []string{"admin_auth"}
	path := constant.ApiV1 + "/auth"

	huma.Register(humaRouter, huma.Operation{
		OperationID: RegisterOwner.Str(),
		Method:      http.MethodPost,
		Path:        path + "/signup_owner",
		Tags:        tags}, handler.RegisterOwner,
	)

	huma.Register(humaRouter, huma.Operation{
		OperationID: VerifyOwner.Str(),
		Method:      http.MethodPost,
		Path:        path + "/verify_owner",
		Tags:        tags}, handler.VerifyCompanyOwner,
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

	//===============  Depricated

	// huma.Register(humaRouter, huma.Operation{
	// 	OperationID: VerifyOwner.Str(),
	// 	Method:      http.MethodPost,
	// 	Path:        path + "/verify_owner",
	// 	Middlewares: huma.Middlewares{cmnServ.Authorization(VerifyOwner, true, PermissionMap[VerifyOwner]...)},
	// 	Tags:        tags}, handler.VerifyOwner,
	// )
	// huma.Register(humaRouter, huma.Operation{
	// 	OperationID: ApproveCompany.Str(),
	// 	Method:      http.MethodPost,
	// 	Path:        path + "/approve_company",
	// 	Middlewares: huma.Middlewares{cmnServ.Authorization(ApproveCompany, true, PermissionMap[ApproveCompany]...)},
	// 	Tags:        tags}, handler.ApproveCompany,
	// )

}
