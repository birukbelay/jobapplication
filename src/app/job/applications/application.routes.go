package applications

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
	OffsetPaginatedApplications = consts.OperationId("App-1-OffsetPaginatedApplications")
	GetOneApplicationById       = consts.OperationId("App-2-GetOneApplicationById")

	UpdateApplication = consts.OperationId("App-4-UpdateApplication")
	DeleteApplication = consts.OperationId("App-5-DeleteApplication")

	//applicants
	CreateApplication = consts.OperationId("App-3-CreateApplication")
	GetMyApplications = consts.OperationId("App-6-GetMyApplications")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{
	OffsetPaginatedApplications: {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: "Get paginated list of applications (for job owners)"},
	GetOneApplicationById:       {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: "Get a specific application by ID"},
	CreateApplication:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Apply for a job"},
	UpdateApplication:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Update application (applicant can update details, owner can update status)"},
	DeleteApplication:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Withdraw application"},
	GetMyApplications:           {AllowedRoles: []string{enums.APPLICANT.S()}, Description: "Get current user's applications"},
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

	// Get application by ID
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetOneApplicationById.Str(),
		Description: "Get a specific application by its ID",
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneApplicationById, OperationMap[GetOneApplicationById].AllowedRoles, nil)},
	}, genericController.GHandler.GetOneById)

	// Create new application (apply for job)
	huma.Register(humaRouter, huma.Operation{
		OperationID: CreateApplication.Str(),
		Description: "Apply for a job by creating an application",
		Method:      http.MethodPost,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateApplication, OperationMap[CreateApplication].AllowedRoles, nil)},
	}, genericController.CreateApplication)

	// Update application
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateApplication.Str(),
		Description: "Update an application (applicant can update details, job owner can update status)",
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateApplication, OperationMap[UpdateApplication].AllowedRoles, nil)},
	}, genericController.GHandler.UpdateOneById)

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
		Description: "Get current user's job applications",
		Method:      http.MethodGet,
		Path:        myApplicationsPath,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetMyApplications, OperationMap[GetMyApplications].AllowedRoles, nil)},
	}, genericController.GetMyApplications)
}
