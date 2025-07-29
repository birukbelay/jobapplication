package company

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

type HumaCompanyHandler struct {
	// CmnServ *gen.IGenericGormServ
	Service  *Service
	GHandler *generic.IGenericController[models.Company, models.CompanyDto, models.CompanyUpdateDto, models.CompanyFilter, models.CompanyQuery]
	// DbService *gen.IGenericGormServT[models.Company, models.CompanyDto, models.LogTraceDto]
}

// NewLogTraceHandler creates a content handler from IContentService & Generic Gorm Service
func NewCompanyHandler(serv *Service) *HumaCompanyHandler {
	return &HumaCompanyHandler{Service: serv, GHandler: generic.NewGenericController[models.Company, models.CompanyDto, models.CompanyUpdateDto, models.CompanyFilter, models.CompanyQuery](serv.ProvServ.GormConn)}
}

const (
	CreateCompany   = consts.OperationId("Ow_Cp-1-CreateCompany")
	GetMyCompany    = consts.OperationId("Ow_Cp-2-GetMyCompany")
	UpdateMyCompany = consts.OperationId("Ow_Cp-3-UpdateMyCompany")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{

	CreateCompany:   {AllowedRoles: []string{enums.UNVERIFIED_USER.S()}, Description: ".."},
	GetMyCompany:    {AllowedRoles: []string{enums.OWNER.S(), enums.UNVERIFIED_USER.S()}, Description: ".."},
	UpdateMyCompany: {AllowedRoles: []string{enums.OWNER.S()}, Description: ".."},
}

func SetupOwnerCompanyRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewCompanyHandler(serv)

	tags := []string{"owner-manage_company"}
	path := consts.ApiV1 + "/my_company"
	// pathId := consts.ApiV1 + "/my_company/{id}"
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateMyCompany.Str(),
		Method:      http.MethodPost,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateCompany,  OperationMap[CreateCompany].AllowedRoles,nil)},
	}, genericController.CreateCompany,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetMyCompany.Str(),
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetMyCompany,  OperationMap[GetMyCompany].AllowedRoles,nil)},
	}, genericController.GetMyCompany,
	)
	//------------- used to approve companies
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateMyCompany.Str(),
		Method:      http.MethodPatch,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateMyCompany,  OperationMap[UpdateMyCompany].AllowedRoles,nil)},
	}, genericController.UpdateMyCompany,
	)

}
