package inviteCode

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

type HumaInviteCodeHandler struct {
	// CmnServ *gen.IGenericGormServ
	Service  *Service
	GHandler *generic.IGenericController[models.InviteCode, models.InviteCodeDto, models.InviteCodeUpdateDto, models.InviteCodeFilter, models.InviteCodeQuery]
	// DbService *gen.IGenericGormServT[models.InviteCode, models.InviteCodeDto, models.LogTraceDto]
}

// NewLogTraceHandler creates a content handler from IContentService & Generic Gorm Service
func NewInviteCodeHandler(serv *Service) *HumaInviteCodeHandler {
	return &HumaInviteCodeHandler{Service: serv, GHandler: generic.NewGenericController[models.InviteCode, models.InviteCodeDto, models.InviteCodeUpdateDto, models.InviteCodeFilter, models.InviteCodeQuery](serv.ProvServ.GormConn)}
}

const (
	CreateInviteCode        = consts.OperationId("Ow_Invt-1-CreateInviteCode")
	GetOneInviteCode        = consts.OperationId("Ow_Invt-2-GetOneInviteCode")
	UpdateOneInviteCode     = consts.OperationId("Ow_Invt-3-UpdateOneInviteCode")
	CreateAndSendInviteCode = consts.OperationId("Ow_Invt-4-UpdateOneInviteCode")
	JoinViaInviteCode       = consts.OperationId("Ow_Invt-5-UpdateOneInviteCode")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{

	CreateInviteCode:        {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
	GetOneInviteCode:        {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
	UpdateOneInviteCode:     {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
	CreateAndSendInviteCode: {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
	JoinViaInviteCode:       {AllowedRoles: []string{enums.USER.S()}, Description: ".."},
}

func SetupInviteCodeRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewInviteCodeHandler(serv)

	tags := []string{"owner-invite_codes"}
	path := consts.ApiV1 + "/invite_code"
	pathId := consts.ApiV1 + "/invite_code/{id}"
	huma.Register(humaRouter, huma.Operation{
		OperationID: CreateInviteCode.Str(),
		Method:      http.MethodPost,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateInviteCode, true, OperationMap[CreateInviteCode].AllowedRoles...)},
	}, genericController.CreateInviteCode,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetOneInviteCode.Str(),
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneInviteCode, true, OperationMap[GetOneInviteCode].AllowedRoles...)},
	}, genericController.GetOneInviteCode,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateOneInviteCode.Str(),
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateOneInviteCode, true, OperationMap[UpdateOneInviteCode].AllowedRoles...)},
	}, genericController.UpdateOneInviteCode,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: CreateAndSendInviteCode.Str(),
		Method:      http.MethodPost,
		Path:        path + "/send_one",
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateAndSendInviteCode, true, OperationMap[CreateAndSendInviteCode].AllowedRoles...)},
	}, genericController.CreateAndSendInviteCode,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: JoinViaInviteCode.Str(),
		Method:      http.MethodPost,
		Path:        path + "/join",
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(JoinViaInviteCode, true, OperationMap[JoinViaInviteCode].AllowedRoles...)},
	}, genericController.JoinViaInviteCode,
	)

}
