package profile

import (
	"context"
	"net/http"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	sql_db "github.com/birukbelay/gocmn/src/generic"
	"github.com/danielgtaylor/huma/v2"

	IConst "github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models"
)

func (uh *HProfileHandler) GetMyProfile(ctx context.Context, _ *dtos.AuthParam) (*dtos.HumaResponse[dtos.GResp[models.User]], error) {
	v, ok := ctx.Value(string(IConst.CtxClaims)).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := sql_db.DbGetOneByID[models.User](uh.CmnServ.GormConn, ctx, v.UserId, nil)
	return dtos.HumaReturnG[models.User](resp, err)
}

func (uh *HProfileHandler) UpdateMyProfile(ctx context.Context, filter *dtos.HumaReqBody[models.ProfileUpdateDto]) (*dtos.HumaResponse[dtos.GResp[models.User]], error) {
	v, ok := ctx.Value(string(IConst.CtxClaims)).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := sql_db.DbUpdateOneById[models.User](uh.CmnServ.GormConn, ctx, v.UserId, filter.Body, nil)
	return dtos.HumaReturnG[models.User](resp, err)
}

func (uh *HProfileHandler) UpdateMyPassword(ctx context.Context, input *dtos.HumaReqBody[models.PasswordUpdateDto]) (*dtos.HumaResponse[dtos.GResp[models.User]], error) {
	v, ok := ctx.Value(string(IConst.CtxClaims)).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := uh.Service.ChangePassword(ctx, v.UserId, input.Body)
	return dtos.HumaReturnG(resp, err)
}

// func (uh *HProfileHandler) DeleteMyProfile(ctx context.Context, _ *dtos.AuthParam) (*dtos.HumaResponse[dtos.GResp[models.User]], error) {
// 	v, ok := ctx.Value(string(IConst.CtxClaims)).(crypto.CustomClaims)
// 	if !ok {
// 		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
// 	}
// 	resp, err := sql_db.DbUpdateOneById[models.User](uh.CmnServ.GormConn, ctx, v.UserId, models.UserDto{Active: false, AccountStatus: enums.AccountDeleted}, nil)
// 	return dtos.HumaReturn[models.User](resp, err)
// }
