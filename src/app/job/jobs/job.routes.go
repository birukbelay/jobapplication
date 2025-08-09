package jobs

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
	OffsetPaginatedJobs = consts.OperationId("Job-1-OffsetPaginatedJobs")
	GetOneJobById       = consts.OperationId("Job-2-GetOneJobById")
	CreateJob           = consts.OperationId("Job-3-CreateJob")
	UpdateJob           = consts.OperationId("Job-4-UpdateJob")
	DeleteJob           = consts.OperationId("Job-5-DeleteJob")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{
	OffsetPaginatedJobs: {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: "Get paginated list of jobs"},
	GetOneJobById:       {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: "Get a specific job by ID"},
	CreateJob:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Create a new job posting"},
	UpdateJob:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Update an existing job"},
	DeleteJob:           {AllowedRoles: []string{enums.COMPANY.S()}, Description: "Delete a job posting"},
}

func SetupJobRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *Service) {
	genericController := NewHandler(serv)

	tags := []string{"jobs"}
	path := consts.ApiV1 + "/jobs"
	pathId := consts.ApiV1 + "/jobs/{id}"

	// Get paginated jobs
	huma.Register(humaRouter, huma.Operation{
		OperationID: OffsetPaginatedJobs.Str(),
		Description: "Get paginated list of jobs with filtering and sorting",
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(OffsetPaginatedJobs, OperationMap[OffsetPaginatedJobs].AllowedRoles, nil)},
	}, genericController.OffsetPaginated)

	// Get job by ID
	huma.Register(humaRouter, huma.Operation{
		OperationID: GetOneJobById.Str(),
		Description: "Get a specific job by its ID",
		Method:      http.MethodGet,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneJobById, OperationMap[GetOneJobById].AllowedRoles, nil)},
	}, genericController.GHandler.GetOneById)

	// Create new job
	huma.Register(humaRouter, huma.Operation{
		OperationID: CreateJob.Str(),
		Description: "Create a new job posting",
		Method:      http.MethodPost,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(CreateJob, OperationMap[CreateJob].AllowedRoles, nil)},
	}, genericController.GHandler.CreateOne)

	// Update job
	huma.Register(humaRouter, huma.Operation{
		OperationID: UpdateJob.Str(),
		Description: "Update an existing job posting",
		Method:      http.MethodPatch,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(UpdateJob, OperationMap[UpdateJob].AllowedRoles, nil)},
	}, genericController.GHandler.UpdateOneById)

	// Delete job
	huma.Register(humaRouter, huma.Operation{
		OperationID: DeleteJob.Str(),
		Description: "Delete a job posting",
		Method:      http.MethodDelete,
		Path:        pathId,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(DeleteJob, OperationMap[DeleteJob].AllowedRoles, nil)},
	}, genericController.GHandler.DeleteOneByID)
}
