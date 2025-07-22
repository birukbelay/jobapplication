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
	GetMyProfile:    {AllowedRoles: []string{enums.USER.S()}, Description: "Creating A User"},
	UpdateMyProfile: {AllowedRoles: []string{enums.USER.S()}},
	ChangeMyPwd:     {AllowedRoles: []string{enums.USER.S()}},
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

	huma.Register(humaRouter, huma.Operation{
		OperationID: ChangeEmailReq.Str(),
		Method:      http.MethodPost,
		Path:        path + "/change_email",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(ChangeEmailReq, true, ProfilePermissionsMap[AdminChangeEmailReq].AllowedRoles...)},
	}, adminHandler.UpdateMyEmailReq)

	huma.Register(humaRouter, huma.Operation{
		OperationID: ChangeEmailVerify.Str(),
		Method:      http.MethodPost,
		Path:        path + "/verify_change_email",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(ChangeEmailVerify, true, ProfilePermissionsMap[ChangeEmailVerify].AllowedRoles...)},
	}, adminHandler.VerifyMyChangeEmailReq)

}
