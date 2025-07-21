package profile

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

// HProfileHandler .
type HProfileHandler[T models.IntUsr] struct {
	CmnServ *providers.IProviderS
	Service *Service[T]
}

// NewProfileHandler creates a profile handler from ProfileService & Generic Gorm Service
func NewProfileHandler[T models.IntUsr](serv *providers.IProviderS, service *Service[T]) *HProfileHandler[T] {
	return &HProfileHandler[T]{CmnServ: serv, Service: service}
}

const (
	GetMyProfile      = constant.OperationId("Pr_1-GetMyProfile")
	UpdateMyProfile   = constant.OperationId("Pr_2-UpdateMyProfile")
	ChangeMyPwd       = constant.OperationId("Pr_3-ChangeMyPassword")
	ChangeEmailReq    = constant.OperationId("Pr_4-ChangeEmailReq")
	ChangeEmailVerify = constant.OperationId("Pr_5-ChangeEmailVerify")
)

var ProfilePermissionsMap = map[constant.OperationId]models.OperationAccess{
	GetMyProfile:    {AllowedRoles: []string{enums.ADMIN.S()}, Description: "Creating A User"},
	UpdateMyProfile: {AllowedRoles: []string{enums.ADMIN.S()}},
	ChangeMyPwd:     {AllowedRoles: []string{enums.ADMIN.S()}},
}

func SetUserProfileRoutes(humaRouter huma.API, provServ *providers.IProviderS) {
	adminHandler := NewProfileHandler(provServ, NewProfileServH[models.User](provServ))
	tags := []string{"user_profile"}
	path := constant.ApiV1 + "/profile"
	//pathId := system_const.ApiV1 + "/profile/{id}"

	huma.Register(humaRouter, huma.Operation{
		OperationID: GetMyProfile.Str(),
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(GetMyProfile, true, ProfilePermissionsMap[GetMyProfile].AllowedRoles...)},
	}, adminHandler.GetMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateMyProfile.Str(),
		Method:      http.MethodPatch,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(UpdateMyProfile, true)},
	}, adminHandler.UpdateMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: ChangeMyPwd.Str(),
		Method:      http.MethodPost,
		Path:        path + "/change_pwd",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(ChangeMyPwd, true)},
	}, adminHandler.UpdateMyPassword)

}

const (
	GetAdminProfile    = constant.OperationId("OwPr-1_GetAdminProfile")
	UpdateAdminProfile = constant.OperationId("OwPr-2_UpdateAdminProfile")
	// DeleteMyProfile = constant.OperationId("DeleteMyProfile")
	ChangeAdminPwd = constant.OperationId("OwPr-3_ChangeAdminPassword")

	// //=== company related
	// DeactivateMyCompany = constant.OperationId("Ow-4_DeactivateMyCompany")
	// UpdateMyCompany     = constant.OperationId("Ow-3_UpdateMyCompany")
	// ApproveCompany      = constant.OperationId("Ad-1_ApproveCompany")
)

var AdminProfilePermissionsMap = map[constant.OperationId]models.OperationAccess{
	GetAdminProfile:    {AllowedRoles: []string{enums.ADMIN.S(), enums.OWNER.S()}, Description: "Creating A User"},
	UpdateAdminProfile: {AllowedRoles: []string{enums.ADMIN.S(), enums.OWNER.S()}},
	ChangeAdminPwd:     {AllowedRoles: []string{enums.ADMIN.S(), enums.OWNER.S()}},
}

func SetAdminProfileRoutes(humaRouter huma.API, provServ *providers.IProviderS) {
	adminHandler := NewProfileHandler(provServ, NewProfileServH[models.Admin](provServ))
	tags := []string{"admin_profile"}
	path := constant.ApiV1 + "/admin_profile"
	//pathId := system_const.ApiV1 + "/profile/{id}"

	huma.Register(humaRouter, huma.Operation{
		OperationID: GetAdminProfile.Str(),
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(GetAdminProfile, true, ProfilePermissionsMap[GetAdminProfile].AllowedRoles...)},
	}, adminHandler.GetMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateAdminProfile.Str(),
		Method:      http.MethodPatch,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(UpdateAdminProfile, true, ProfilePermissionsMap[UpdateAdminProfile].AllowedRoles...)},
	}, adminHandler.UpdateMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: ChangeAdminPwd.Str(),
		Method:      http.MethodPost,
		Path:        path + "/change_pwd",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(ChangeAdminPwd, true, ProfilePermissionsMap[ChangeAdminPwd].AllowedRoles...)},
	}, adminHandler.UpdateMyPassword)

	// ===============  Depricated

	// huma.Register(humaRouter, huma.Operation{
	// 	OperationID: DeleteMyProfile.Str(),
	// 	Method:      http.MethodDelete,
	// 	Path:        path,
	// 	Tags:        tags,
	// 	Middlewares: huma.Middlewares{provServ.Authorization(DeleteMyProfile, true)},
	// }, handler.DeleteMyProfile)

	// huma.Register(humaRouter, huma.Operation{
	// 	OperationID: ApproveCompany.Str(),
	// 	Method:      http.MethodPost,
	// 	Path:        path + "/approve_company",
	// 	Middlewares: huma.Middlewares{cmnServ.Authorization(ApproveCompany, true, PermissionMap[ApproveCompany]...)},
	// 	Tags:        tags}, handler.ApproveCompany,
	// )

}
