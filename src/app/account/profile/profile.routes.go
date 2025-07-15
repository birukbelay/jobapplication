package profile

import (
	"net/http"

	constant "github.com/birukbelay/gocmn/src/const"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

// HProfileHandler .
type HProfileHandler struct {
	CmnServ *providers.IProviderS
	Service *Service
}

// NewProfileHandler creates a profile handler from ProfileService & Generic Gorm Service
func NewProfileHandler(serv *providers.IProviderS, service *Service) *HProfileHandler {
	return &HProfileHandler{CmnServ: serv, Service: service}
}

const (
	GetMyProfile    = constant.OperationId("GetMyProfile")
	UpdateMyProfile = constant.OperationId("UpdateMyProfile")
	DeleteMyProfile = constant.OperationId("DeleteMyProfile")
	ChangeMyPwd     = constant.OperationId("ChangeMyPassword")
)

// var ProfileOpArr = []constant.OperationId{GetMyProfile, UpdateMyProfile, DeleteMyProfile}

var PermissionMap = map[constant.OperationId][]enums.Role{
	GetMyProfile:    {},
	UpdateMyProfile: {},
	DeleteMyProfile: {},
	ChangeMyPwd:     {},
}

func SetProfileRoutes(humaRouter huma.API, provServ *providers.IProviderS) {
	handler := NewProfileHandler(provServ, NewProfileServH(provServ))
	tags := []string{"account_profile"}
	path := constant.ApiV1 + "/profile"
	//pathId := system_const.ApiV1 + "/profile/{id}"

	huma.Register(humaRouter, huma.Operation{
		OperationID: GetMyProfile.Str(),
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(GetMyProfile, true, PermissionMap[GetMyProfile]...)},
	}, handler.GetMyProfile)

	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateMyProfile.Str(),
		Method:      http.MethodPatch,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(UpdateMyProfile, true)},
	}, handler.UpdateMyProfile)

	// huma.Register(humaRouter, huma.Operation{
	// 	OperationID: DeleteMyProfile.Str(),
	// 	Method:      http.MethodDelete,
	// 	Path:        path,
	// 	Tags:        tags,
	// 	Middlewares: huma.Middlewares{provServ.Authorization(DeleteMyProfile, true)},
	// }, handler.DeleteMyProfile)
	huma.Register(humaRouter, huma.Operation{
		OperationID: ChangeMyPwd.Str(),
		Method:      http.MethodPost,
		Path:        path + "/change_pwd",
		Tags:        tags,
		Middlewares: huma.Middlewares{provServ.Authorization(ChangeMyPwd, true)},
	}, handler.UpdateMyPassword)

}
