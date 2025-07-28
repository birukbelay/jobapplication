package companies

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
	GHandler *generic.IGenericController[models.Company, models.CompanyDto, models.CompanyDto, models.CompanyFilter, models.CompanyQuery]
	// DbService *gen.IGenericGormServT[models.Company, models.CompanyDto, models.LogTraceDto]
}

// NewLogTraceHandler creates a content handler from IContentService & Generic Gorm Service
func NewCompanyHandler(serv *Service) *HumaCompanyHandler {
	return &HumaCompanyHandler{Service: serv, GHandler: generic.NewGenericController[models.Company, models.CompanyDto, models.CompanyDto, models.CompanyFilter, models.CompanyQuery](serv.ProvServ.GormConn)}
}

const (
	OffsetPaginatedCompanies = consts.OperationId("Ad_Cp-1-OffsetPaginatedCompanies")
	GetOneCompanyById        = consts.OperationId("Ad_Cp-2-GetOneCompanyById")
	UpdateCompany            = consts.OperationId("Ad_Cp-3-UpdateCompany")
	// DisableCompany = consts.OperationId("Pr_1-DisableCompany")
	// OwnerCreateCompany       = consts.OperationId("Pr_1-GetMyProfile")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{
	OffsetPaginatedCompanies: {AllowedRoles: []string{enums.PLATFORM_ADMIN.S()}, Description: ""},
	GetOneCompanyById:        {AllowedRoles: []string{enums.PLATFORM_ADMIN.S()}, Description: ".."},
	UpdateCompany:            {AllowedRoles: []string{enums.PLATFORM_ADMIN.S()}, Description: ".."},
}

func SetupManageCompaniesRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewCompanyHandler(serv)

	tags := []string{"platform_admin-manage_companies"}
	path := consts.ApiV1 + "/company"
	pathId := consts.ApiV1 + "/company/{id}"
	huma.Register(humaRouter, huma.Operation{
		OperationID: OffsetPaginatedCompanies.Str(),
		Description: "",
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		// Middlewares: huma.Middlewares{cmnServ.Authorization(OffsetPaginatedCompanies, true, OperationMap[OffsetPaginatedCompanies].AllowedRoles...)},
	}, genericController.OffsetPaginated,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetOneCompanyById.Str(),
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		// Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneCompanyById, true, OperationMap[GetOneCompanyById].AllowedRoles...)},
	}, genericController.GHandler.GetOneById,
	)

	//------------- used to approve companies
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateCompany.Str(),
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		// Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateCompany, true, OperationMap[UpdateCompany].AllowedRoles...)},
	}, genericController.GHandler.UpdateOneById,
	)

	//------------- used to approve companies
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateCompany.Str(),
		Method:      http.MethodPatch,
		Path:        pathId + "/approve",
		Tags:        tags,
		// Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateCompany, true, OperationMap[UpdateCompany].AllowedRoles...)},
	}, genericController.ApproveCompany,
	)
}
