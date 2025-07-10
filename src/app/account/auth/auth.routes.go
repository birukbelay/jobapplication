package auth

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	constant "github.com/projTemplate/goauth/src/common"
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
	Login        = constant.OperationId("Au-1_Login")
	RefreshToken = constant.OperationId("Au-2_RefreshToken")

	RegisterOwner       = constant.OperationId("Ow-1_RegisterOwner")
	VerifyOwner         = constant.OperationId("Ow-2_VerifyOwner")
	ApproveCompany      = constant.OperationId("Ad-1_ApproveCompany")
	DeactivateMyCompany = constant.OperationId("Ow-4_DeactivateMyCompany")
	UpdateMyCompany     = constant.OperationId("Ow-3_UpdateMyCompany")
)

var PermissionMap = map[constant.OperationId][]enums.Role{
	RegisterOwner:  {},
	VerifyOwner:    {enums.ADMIN},
	Login:          {},
	RefreshToken:   {},
	ApproveCompany: {enums.ADMIN},
}

func SetupAuthRoutes(humaRouter huma.API, providerS *providers.IProviderS, serv *Service) {
	handler := NewAuthHandler(providerS, serv)
	tags := []string{"account_auth"}
	path := constant.ApiV1 + "/auth"

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

	//===============  Depricated

	huma.Register(humaRouter, huma.Operation{
		OperationID: RegisterOwner.Str(),
		Method:      http.MethodPost,
		Path:        path + "/register_owner",
		Tags:        tags}, handler.RegisterOwner,
	)
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
