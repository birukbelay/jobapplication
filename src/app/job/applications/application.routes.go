package applications

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

type HumaHandler struct {
	Service  *Service
	GHandler *generic.IGenericController[models.Application, models.ApplicationDto, models.ApplicationUpdateDto, models.ApplicationFilter, models.ApplicationQuery]
}

func NewHandler(serv *Service) *HumaHandler {
	return &HumaHandler{
		Service:  serv,
		GHandler: generic.NewGenericController[models.Application, models.ApplicationDto, models.ApplicationUpdateDto, models.ApplicationFilter, models.ApplicationQuery](serv.ProvServ.GormConn),
	}
}

const (
	OffsetPaginatedApplications = consts.OperationId("App-1-MYCompaniesApplications")
	GetOneApplicationByID       = consts.OperationId("App-2-GetApplicationBYID")

	UpdateApplication   = consts.OperationId("App-3-UpdateApplication")
	UpdateMyApplication = consts.OperationId("App-4-UpdateMyApplication")

	//applicants
	CreateApplication    = consts.OperationId("App-5-CreateApplication")
	GetMyApplications    = consts.OperationId("App-6-GetMyApplications")
	GetMyApplicationById = consts.OperationId("App-7-GetMyApplicationBYID")
	DeleteApplication    = consts.OperationId("App-8-DeleteApplication")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{
	OffsetPaginatedApplications: {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: "Get paginated list of applications (for job owners)"},
	UpdateApplication:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Update application (applicant can update details, owner can update status)"},

	//applicants
	CreateApplication:    {AllowedRoles: []string{enums.APPLICANT.S()}, Description: "Apply for a job"},
	GetMyApplicationById: {AllowedRoles: []string{enums.APPLICANT.S(), enums.PLATFORM_ADMIN.S()}, Description: "Get  ID"},
	GetMyApplications:    {AllowedRoles: []string{enums.APPLICANT.S()}, Description: "Get current user's applications"},
	DeleteApplication:    {AllowedRoles: []string{enums.APPLICANT.S()}, Description: "Withdraw application"},
}

func SetupApplicationRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewHandler(serv)

	tags := []string{"applications"}
	path := consts.ApiV1 + "/applications"
	pathId := consts.ApiV1 + "/applications/{id}"
	myApplicationsPath := consts.ApiV1 + "/my-applications"

	// Get paginated applications (for job owners/admins)
	huma.Register(humaRouter, huma.Operation{
		OperationID: OffsetPaginatedApplications.Str(),
		Description: "Get paginated list of applications with filtering and sorting (for job owners and admins)",
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(OffsetPaginatedApplications, OperationMap[OffsetPaginatedApplications].AllowedRoles, nil)},
	}, genericController.OffsetPaginated)

	// Update application
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateApplication.Str(),
		Description: "Update an application , job owner can update status)",
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateApplication, OperationMap[UpdateApplication].AllowedRoles, nil)},
	}, genericController.UpdateApplications)

	//==============  Aplicants ===============

	// Update  My application application
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateMyApplication.Str(),
		Description: "Update My application (applicant can update details",
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateMyApplication, OperationMap[UpdateMyApplication].AllowedRoles, nil)},
	}, genericController.GHandler.UpdateOneById)

	// Get application by ID
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetMyApplicationById.Str(),
		Description: "Get THe applicants own application by ID",
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetMyApplicationById, OperationMap[GetMyApplicationById].AllowedRoles, nil)},
	}, genericController.GHandler.GetOneById)

	// Delete application (withdraw)
	huma.Register(humaRouter, huma.Operation{
		OperationID: DeleteApplication.Str(),
		Description: "Withdraw/delete an application",
		Method:      http.MethodDelete,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(DeleteApplication, OperationMap[DeleteApplication].AllowedRoles, nil)},
	}, genericController.GHandler.DeleteOneByID)

	// Get current user's applications
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetMyApplications.Str(),
		Description: "Get current applicants's job applications",
		Method:      http.MethodGet,
		Path:        myApplicationsPath,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetMyApplications, OperationMap[GetMyApplications].AllowedRoles, nil)},
	}, genericController.GetMyApplications)
	// Create new application (apply for job)
	huma.Register(humaRouter, huma.Operation{
		OperationID: CreateApplication.Str(),
		Description: "Apply for a job by creating an application",
		Method:      http.MethodPost,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateApplication, OperationMap[CreateApplication].AllowedRoles, nil)},
	}, genericController.CreateApplication)
}
