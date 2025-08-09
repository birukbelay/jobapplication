package admins

import (
	"net/http"

	"github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

type Service struct {
	ProvServ *providers.IProviderS
}

func NewService(genServ *providers.IProviderS) *Service {
	return &Service{
		ProvServ: genServ,
	}
}

const (
	OffsetPaginatedAdmins = consts.OperationId("Ad-1-OffsetPaginatedUsers")
	GetOneAdminById       = consts.OperationId("Ad-2-GetOneUserById")
	UpdateAdmin           = consts.OperationId("Ad-3-UpdateUser")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{
	OffsetPaginatedAdmins: {AllowedRoles: []string{enums.PLATFORM_ADMIN.S()}, Description: ""},
	GetOneAdminById:       {AllowedRoles: []string{enums.PLATFORM_ADMIN.S()}, Description: ".."},
	UpdateAdmin:           {AllowedRoles: []string{enums.PLATFORM_ADMIN.S()}, Description: ".."},
}

func SetupManageAdminUsersRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewHandler(serv)

	tags := []string{"platform_admin-manage_admins"}
	path := consts.ApiV1 + "/admin"
	pathId := consts.ApiV1 + "/admin/{id}"
	huma.Register(humaRouter, huma.Operation{
		OperationID: OffsetPaginatedAdmins.Str(),
		Description: "admins are platform admins and company owners",
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(OffsetPaginatedAdmins, OperationMap[OffsetPaginatedAdmins].AllowedRoles, nil)},
	}, genericController.OffsetPaginated,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetOneAdminById.Str(),
		Description: "admins are platform admins and company owners",
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneAdminById, OperationMap[GetOneAdminById].AllowedRoles, nil)},
	}, genericController.GHandler.GetOneById,
	)

	//--------------
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateAdmin.Str(),
		Description: "admins are platform admins and company owners: use this route to approve admins",
		Method:      http.MethodPatch,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateAdmin, OperationMap[UpdateAdmin].AllowedRoles, nil)},
	}, genericController.GHandler.UpdateOneById,
	)
}
