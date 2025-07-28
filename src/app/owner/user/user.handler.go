package user

import (
	"context"
	"net/http"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models"
)

func (uh *HumaUserHandler) CreateUser(ctx context.Context, dto *dtos.HumaReqBody[UserCreateDto]) (*dtos.HumaResponse[dtos.GResp[bool]], error) {

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	usr, err := uh.Service.CreateUser(ctx, dto.Body, v.CompanyId)
	return dtos.HumaReturnG(usr, err)
}

func (uh *HumaUserHandler) GetOneUser(ctx context.Context, input *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[models.User]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := generic.DbGetOne[models.User](uh.GHandler.GormConn, ctx, models.UserFilter{CompanyID: v.CompanyId, ID: input.ID}, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}

func (uh *HumaUserHandler) UpdateOneUser(ctx context.Context, dto *dtos.HumaReqBodyId[models.UserUpdateDto]) (*dtos.HumaResponse[dtos.GResp[models.User]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}

	resp, err := generic.DbUpdateByFilter[models.User](uh.GHandler.GormConn, ctx, models.UserFilter{CompanyID: v.CompanyId, ID: dto.ID}, dto.Body, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}
