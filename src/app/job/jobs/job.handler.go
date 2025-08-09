package jobs

import (
	"context"
	"net/http"

	"github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
)

func (jh *HumaHandler) OffsetPaginated(ctx context.Context, filter *struct {
	models.JobFilter
	models.JobQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Job]], error) {
	sort, selectedFields := filter.JobQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	v, ok := ctx.Value(consts.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	filter.JobFilter.CreatedBy = v.UserId
	resp, err := generic.DbFetchManyWithOffset[models.Job](jh.GHandler.GormConn, ctx, filter.JobFilter, filter.PaginationInput, nil)
	return dtos.PHumaReturn(resp, err)
}

func (jh *HumaHandler) GetOpenJobs(ctx context.Context, filter *struct {
	models.JobFilter
	models.JobQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Job]], error) {
	sort, selectedFields := filter.JobQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	filter.JobFilter.JobStatus = enums.STATUS_OPEN
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
