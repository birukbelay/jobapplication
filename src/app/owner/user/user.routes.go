package user

import (
	"net/http"

	"github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/generic"
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

type HumaUserHandler struct {
	// CmnServ *gen.IGenericGormServ
	Service  *Service
	GHandler *generic.IGenericController[models.User, models.UserDto, models.UserUpdateDto, models.UserFilter, models.UserQuery]
	// DbService *gen.IGenericGormServT[models.User, models.UserDto, models.LogTraceDto]
}

// NewLogTraceHandler creates a content handler from IContentService & Generic Gorm Service
func NewUserHandler(serv *Service) *HumaUserHandler {
	return &HumaUserHandler{Service: serv, GHandler: generic.NewGenericController[models.User, models.UserDto, models.UserUpdateDto, models.UserFilter, models.UserQuery](serv.ProvServ.GormConn)}
}

const (
	CreateUser    = consts.OperationId("Ow_Usr_1-CreateUser")
	GetOneUser    = consts.OperationId("Ow_Usr_2-GetOneUser")
	UpdateOneUser = consts.OperationId("Ow_Usr_3-UpdateOneUser")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{

	CreateUser:    {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
	GetOneUser:    {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
	UpdateOneUser: {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
}

func SetupCompanyUserRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewUserHandler(serv)

	tags := []string{"owner-manage_users"}
	path := consts.ApiV1 + "/my_users"
	pathId := consts.ApiV1 + "/my_users/{id}"
	huma.Register(humaRouter, huma.Operation{
		OperationID: CreateUser.Str(),
		Method:      http.MethodPost,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateUser, true, OperationMap[CreateUser].AllowedRoles...)},
	}, genericController.CreateUser,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetOneUser.Str(),
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneUser, true, OperationMap[GetOneUser].AllowedRoles...)},
	}, genericController.GetOneUser,
	)
	//------------- used to approve companies
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateOneUser.Str(),
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		// Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateMyUser, true, OperationMap[UpdateMyUser].AllowedRoles...)},
	}, genericController.UpdateOneUser,
	)

}
