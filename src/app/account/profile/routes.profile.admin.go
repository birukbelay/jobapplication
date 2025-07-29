package profile

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

const (
	GetAdminProfile    = constant.OperationId("OwPr_1-GetAdminProfile")
	UpdateAdminProfile = constant.OperationId("OwPr_2-UpdateAdminProfile")
	// DeleteMyProfile = constant.OperationId("DeleteMyProfile")
	AdminChangePwd         = constant.OperationId("OwPr_3-AdminChangePwd")
	AdminChangeEmailReq    = constant.OperationId("OwPr_4-AdminChangeEmailReq")
	AdminChangeEmailVerify = constant.OperationId("OwPr_5-AdminChangeEmailVerify")

	// //=== company related
	// DeactivateMyCompany = constant.OperationId("Ow-4_DeactivateMyCompany")
	// UpdateMyCompany     = constant.OperationId("Ow-3_UpdateMyCompany")
	// ApproveCompany      = constant.OperationId("Ad-1_ApproveCompany")
)

var AdminProfilePermissionsMap = map[constant.OperationId]models.OperationAccess{
	GetAdminProfile:        {AllowedRoles: []string{enums.PLATFORM_ADMIN.S(), enums.OWNER.S()}, Description: "Creating A User"},
	UpdateAdminProfile:     {AllowedRoles: []string{enums.PLATFORM_ADMIN.S(), enums.OWNER.S()}},
	AdminChangePwd:         {AllowedRoles: []string{enums.PLATFORM_ADMIN.S(), enums.OWNER.S()}},
	AdminChangeEmailReq:    {AllowedRoles: []string{enums.PLATFORM_ADMIN.S(), enums.OWNER.S()}},
	AdminChangeEmailVerify: {AllowedRoles: []string{enums.PLATFORM_ADMIN.S(), enums.OWNER.S()}},
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
		Middlewares: huma.Middlewares{provServ.Authorization(GetAdminProfile, ProfilePermissionsMap[GetAdminProfile].AllowedRoles, nil)},
	}, adminHandler.GetMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateAdminProfile.Str(),
		Method:      http.MethodPatch,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(UpdateAdminProfile, ProfilePermissionsMap[UpdateAdminProfile].AllowedRoles, nil)},
	}, adminHandler.UpdateMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: AdminChangePwd.Str(),
		Method:      http.MethodPost,
		Path:        path + "/change_pwd",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(AdminChangePwd, ProfilePermissionsMap[AdminChangePwd].AllowedRoles, nil)},
	}, adminHandler.UpdateMyPassword)

	huma.Register(humaRouter, huma.Operation{
		OperationID: AdminChangeEmailReq.Str(),
		Method:      http.MethodPost,
		Path:        path + "/change_email",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(AdminChangeEmailReq, ProfilePermissionsMap[AdminChangeEmailReq].AllowedRoles, nil)},
	}, adminHandler.UpdateMyEmailReq)

	huma.Register(humaRouter, huma.Operation{
		OperationID: AdminChangeEmailVerify.Str(),
		Method:      http.MethodPost,
		Path:        path + "/verify_change_email",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(AdminChangeEmailVerify, ProfilePermissionsMap[AdminChangeEmailVerify].AllowedRoles, nil)},
	}, adminHandler.VerifyMyChangeEmailReq)

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
