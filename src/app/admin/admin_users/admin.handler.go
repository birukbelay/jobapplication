package admins

import (
	"context"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"

	"github.com/projTemplate/goauth/src/models"
)

type HumaHandler struct {
	Service  *Service
	GHandler *generic.IGenericController[models.User, models.UserDto, models.UserDto, models.UserFilter, models.UserQuery]
}

func NewHandler(serv *Service) *HumaHandler {
	return &HumaHandler{Service: serv, GHandler: generic.NewGenericController[models.User, models.UserDto, models.UserDto, models.UserFilter, models.UserQuery](serv.ProvServ.GormConn)}
}

func (uh *HumaHandler) OffsetPaginated(ctx context.Context, filter *struct {
	models.UserFilter
	models.UserQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.User]], error) {
	sort, selectedFields := filter.UserQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	resp, err := generic.DbFetchManyWithOffset[models.User](uh.GHandler.GormConn, ctx, filter.UserFilter, filter.PaginationInput, nil)
	return dtos.PHumaReturn(resp, err)
}
