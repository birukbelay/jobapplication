package inviteCode

import (
	"context"
	"net/http"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/util"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models"
)

func (uh *HumaInviteCodeHandler) CreateInviteCode(ctx context.Context, dto *dtos.HumaReqBody[models.InviteCodeDto]) (*dtos.HumaResponse[dtos.GResp[models.InviteCode]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}

	dto.Body.CompanyID = v.CompanyId
	dto.Body.Active = true
	//todo: make a unique generator with 10 digit(4-10)
	dto.Body.Code = util.GenerateRandomString(10)

	resp, err := generic.DbCreateOne[models.InviteCode](uh.GHandler.GormConn, ctx, dto.Body, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	// _, err = generic.DbUpdateOneById[models.Admin](uh.GHandler.GormConn, ctx, resp.Body.OwnerID, models.InviteCodeDto{InviteCodeID: &resp.Body.ID}, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *HumaInviteCodeHandler) GetOneInviteCode(ctx context.Context, input *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[models.InviteCode]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := generic.DbGetOne[models.InviteCode](uh.GHandler.GormConn, ctx, models.InviteCodeFilter{CompanyID: v.CompanyId, ID: input.ID}, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}

func (uh *HumaInviteCodeHandler) UpdateOneInviteCode(ctx context.Context, dto *dtos.HumaReqBodyId[models.InviteCodeUpdateDto]) (*dtos.HumaResponse[dtos.GResp[models.InviteCode]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}

	resp, err := generic.DbUpdateByFilter[models.InviteCode](uh.GHandler.GormConn, ctx, models.InviteCodeFilter{CompanyID: v.CompanyId, ID: dto.ID}, dto.Body, nil)
	if err != nil {
		return dtos.HumaReturnG(resp, err)
	}
	return dtos.HumaReturnG(resp, err)
}

func (uh *HumaInviteCodeHandler) CreateAndSendInviteCode(ctx context.Context, dto *dtos.HumaReqBody[models.InviteCodeDto]) (*dtos.HumaResponse[dtos.GResp[models.InviteCode]], error) {

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	if dto.Body.UserInfo == nil {
		return nil, huma.NewError(http.StatusBadRequest, "you need to add user info")
	}

	usr, err := uh.Service.CreateAndSendInviteCode(ctx, *dto.Body.UserInfo, v.CompanyId, dto.Body)
	return dtos.HumaReturnG(usr, err)
}

func (uh *HumaInviteCodeHandler) JoinViaInviteCode(ctx context.Context, dto *dtos.HumaReqBody[models.JoinViaCode]) (*dtos.HumaResponse[dtos.GResp[bool]], error) {

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}

	resp, err := uh.Service.JoinViaInviteCode(ctx, v.UserId, dto.Body.Code)
	return dtos.HumaReturnG(resp, err)
}
