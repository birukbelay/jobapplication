package company

import (
	"context"
	"net/http"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
)

func (uh *HumaCompanyHandler) CreateCompany(ctx context.Context, dto *dtos.HumaReqBody[models.CompanyDto]) (*dtos.HumaResponse[dtos.GResp[models.Company]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	dto.Body.CompanyStatus = enums.CompanyPendingApproval
	dto.Body.OwnerID = v.UserId
	//TODO: check if the user don't have previous company
	resp, err := generic.DbCreateOne[models.Company](uh.GHandler.GormConn, ctx, dto.Body, &generic.Opt{Debug: true})
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	_, err = generic.DbUpdateOneById[models.Admin](uh.GHandler.GormConn, ctx, resp.Body.OwnerID, models.UserDto{CompanyID: &resp.Body.ID}, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *HumaCompanyHandler) GetMyCompany(ctx context.Context, _ *dtos.AuthParam) (*dtos.HumaResponse[dtos.GResp[models.Company]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := generic.DbGetOne[models.Company](uh.GHandler.GormConn, ctx, models.CompanyFilter{OwnerID: v.UserId}, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}

func (uh *HumaCompanyHandler) UpdateMyCompany(ctx context.Context, dto *dtos.HumaReqBodyId[models.CompanyUpdateDto]) (*dtos.HumaResponse[dtos.GResp[models.Company]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}

	resp, err := generic.DbUpdateByFilter[models.Company](uh.GHandler.GormConn, ctx, models.CompanyFilter{OwnerID: v.UserId, ID: dto.ID}, dto.Body, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}
