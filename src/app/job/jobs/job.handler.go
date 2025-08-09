package jobs

import (
	"context"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"

	"github.com/projTemplate/goauth/src/models"
)

type HumaHandler struct {
	Service  *Service
	GHandler *generic.IGenericController[models.Job, models.JobDto, models.JobUpdateDto, models.JobFilter, models.JobQuery]
}

func NewHandler(serv *Service) *HumaHandler {
	return &HumaHandler{
		Service:  serv,
		GHandler: generic.NewGenericController[models.Job, models.JobDto, models.JobUpdateDto, models.JobFilter, models.JobQuery](serv.ProvServ.GormConn),
	}
}

func (jh *HumaHandler) OffsetPaginated(ctx context.Context, filter *struct {
	models.JobFilter
	models.JobQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Job]], error) {
	sort, selectedFields := filter.JobQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	resp, err := generic.DbFetchManyWithOffset[models.Job](jh.GHandler.GormConn, ctx, filter.JobFilter, filter.PaginationInput, nil)
	return dtos.PHumaReturn(resp, err)
}

// func (jh *HumaHandler) CreateJob(ctx context.Context, input *struct {
// 	Body models.JobDto `json:"job"`
// }) (*dtos.HumaResponse[models.Job], error) {
// 	// Set the created_by from context (user ID)
// 	userID := ctx.Value("user_id").(string)
// 	input.Body.SetOnCreate(userID)

// 	resp, err := jh.GHandler.CreateOne(ctx, &struct {
// 		Body models.JobDto `json:"job"`
// 	}{Body: input.Body})
// 	return resp, err
// }
