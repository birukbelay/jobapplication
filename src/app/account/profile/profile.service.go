package profile

import (
	"context"

	sql_db "github.com/birukbelay/gocmn/src/generic"
	ICrypt "github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	ICnst "github.com/birukbelay/gocmn/src/resp_const"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
)

type Service struct {
	ProvServ *providers.IProviderS
}

func NewProfileServH(genServ *providers.IProviderS) *Service {
	return &Service{
		ProvServ: genServ,
	}
}

func (aus *Service) ChangePassword(ctx context.Context, userId string, input models.PasswordUpdateDto) (dtos.GResp[models.User], error) {
	resp, err := sql_db.DbGetOne[models.User](aus.ProvServ.GormConn, ctx, models.UserFilter{ID: userId}, nil)
	if err != nil {
		return resp, err
	}
	valid := ICrypt.BcryptPasswordsMatch(input.OldPassword, resp.Body.Password)
	if !valid {
		return dtos.BadReqM[models.User](ICnst.PasswordDontMatch.Msg()), ICnst.PwdDontMatch
	}
	hash, err := ICrypt.BcryptCreateHash(input.NewPassword)
	if err != nil {
		return dtos.InternalErrMS[models.User]("Hashing Error"), err
	}
	updateResp, err := sql_db.DbUpdateOneById[models.User](aus.ProvServ.GormConn, ctx, userId, models.User{Password: hash}, nil)
	//TODO: reset Logged in things
	return updateResp, err
}
