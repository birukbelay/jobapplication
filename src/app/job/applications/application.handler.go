package applications

import (
	"context"
	"net/http"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/danielgtaylor/huma/v2"

	common "github.com/birukbelay/gocmn/src/consts"
	"github.com/projTemplate/goauth/src/models"
)

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

func (ah *HumaHandler) OffsetPaginated(ctx context.Context, filter *struct {
	models.ApplicationFilter
	models.ApplicationQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Application]], error) {
	sort, selectedFields := filter.ApplicationQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	resp, err := generic.DbFetchManyWithOffset[models.Application](ah.GHandler.GormConn, ctx, filter.ApplicationFilter, filter.PaginationInput, nil)
	return dtos.PHumaReturn(resp, err)
}

// TODO add file handler here
func (ah *HumaHandler) CreateApplication(ctx context.Context, dto *dtos.HumaReqBody[models.ApplicationDto]) (*dtos.HumaResponse[dtos.GResp[models.Application]], error) {

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	dto.Body.ApplicantID = v.UserId
	resp, err := generic.DbCreateOne[models.Application](ah.GHandler.GormConn, ctx, dto.Body, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}

func (ah *HumaHandler) GetMyApplications(ctx context.Context, filter *struct {
	models.ApplicationFilter
	models.ApplicationQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Application]], error) {
	// Filter applications by current user
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}

	sort, selectedFields := filter.ApplicationQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	filter.ApplicationFilter.ApplicantID = v.UserId

	resp, err := generic.DbFetchManyWithOffset[models.Application](ah.GHandler.GormConn, ctx, filter.ApplicationFilter, filter.PaginationInput, nil)
	return dtos.PHumaReturn(resp, err)
}
