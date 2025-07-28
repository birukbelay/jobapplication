package companies

import (
	"context"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
)

func (uh *HumaCompanyHandler) ApproveCompany(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[models.Company]], error) {
	resp, err := generic.DbUpdateOneById[models.Company](uh.GHandler.GormConn, ctx, filter.ID, models.CompanyDto{CompanyStatus: enums.CompanyApproved}, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	//update the admin's status
	_, err = generic.DbUpdateOneById[models.Admin](uh.GHandler.GormConn, ctx, resp.Body.OwnerID, models.UserDto{Active: true, AccountStatus: enums.AccountActive, Role: enums.OWNER}, nil)

	return dtos.HumaReturnG(resp, err)
}
func (uh *HumaCompanyHandler) OffsetPaginated(ctx context.Context, filter *struct {
	models.CompanyFilter
	models.CompanyQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Company]], error) {
	sort, selectedFields := filter.CompanyQuery.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	resp, err := generic.DbFetchManyWithOffset[models.Company](uh.GHandler.GormConn, ctx, filter.CompanyFilter, filter.PaginationInput, nil)
	return dtos.PHumaReturn(resp, err)
}
