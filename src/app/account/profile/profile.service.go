package profile

import (
	"context"

	ICrypt "github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	sql_db "github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/provider/db/redis"
	ICnst "github.com/birukbelay/gocmn/src/resp_const"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
)

type Service[T models.IntUsr] struct {
	ProvServ *providers.IProviderS
}

func NewProfileServH[T models.IntUsr](genServ *providers.IProviderS) *Service[T] {
	return &Service[T]{
		ProvServ: genServ,
	}
}

// ChangePassword .
func (aus *Service[T]) ChangePassword(ctx context.Context, userId, sessionId string, input models.PasswordUpdateDto) (dtos.GResp[T], error) {
	resp, err := sql_db.DbGetOne[T](aus.ProvServ.GormConn, ctx, models.UserFilter{ID: userId}, nil)
	if err != nil {
		return resp, err
	}
	valid := ICrypt.BcryptPasswordsMatch(input.OldPassword, resp.Body.GetPwd())
	if !valid {
		return dtos.BadReqM[T](ICnst.PasswordDontMatch.Msg()), ICnst.PwdDontMatch
	}
	hash, err := ICrypt.BcryptCreateHash(input.NewPassword)
	if err != nil {
		return dtos.InternalErrMS[T]("Hashing Error"), err
	}
	updateResp, err := sql_db.DbUpdateOneById[T](aus.ProvServ.GormConn, ctx, userId, models.UserDto{Password: hash}, nil)
	if err != nil {
		return dtos.InternalErrMS[T]("Update Error"), err
	}

	_, err = sql_db.DbDeleteByFilter[models.Session](aus.ProvServ.GormConn, ctx, models.Session{UserId: userId, SessionId: sessionId}, nil)
	if err != nil {
		return dtos.InternalErrMS[T]("Session Removing error"), err
	}
	//blacklist all the session
	sessions, err := sql_db.DbFetchManyWithOffset[models.Session](aus.ProvServ.GormConn, ctx, models.Session{UserId: userId}, dtos.PaginationInput{Limit: 10000}, nil)
	if err != nil {

	}
	for _, val := range sessions.Body {
		err = redis.BlacklistSession(aus.ProvServ.KeyValServ, ctx, val.SessionId)
	}
	return updateResp, err
}

// SendChangeEmail .params{userId: from token}
func (aus *Service[T]) SendChangeEmail(ctx context.Context, userId string, input models.ChangeEmailReqDto) (dtos.GResp[bool], error) {
	//1.get ther user
	resp, err := sql_db.DbGetOne[T](aus.ProvServ.GormConn, ctx, models.UserFilter{ID: userId}, nil)
	if err != nil {
		return dtos.BadReqC[bool](ICnst.InfoOrCode), ICnst.InfoOrCodeErr
	}
	//2.check if his passowrd is correct
	valid := ICrypt.BcryptPasswordsMatch(input.Password, resp.Body.GetPwd())
	if !valid {
		return dtos.BadReqM[bool](ICnst.InfoOrCode.Msg()), ICnst.InfoOrCodeErr
	}
	//3. check if the email already exists
	_, err = sql_db.DbGetOne[T](aus.ProvServ.GormConn, ctx, models.UserFilter{Email: input.NewEmail}, nil)
	if err == nil {
		return dtos.BadReqC[bool](ICnst.EmailExists), ICnst.EmailExistsErr
	}
	return aus.UTIL_SendVerification(ctx, input.NewEmail, models.GetID(resp.Body), models.ChangeEmail)
}

// VerifyChangeEmail .change the users email
func (aus *Service[T]) VerifyChangeEmail(ctx context.Context, userId string, input models.VerifyEmailDto) (dtos.GResp[bool], error) {
	//1. Get the user
	resp, err := sql_db.DbGetOne[T](aus.ProvServ.GormConn, ctx, models.UserFilter{ID: userId}, nil)
	if err != nil {
		return dtos.BadReqC[bool](ICnst.InfoOrCode), ICnst.InfoOrCodeErr
	}
	//2. Validate the code
	codeValid, email := aus.VerifyCode(ctx, models.GetID(resp.Body), input.Code)
	if !codeValid {
		return dtos.BadReqC[bool](ICnst.InfoOrCode), ICnst.InfoOrCodeErr
	}
	//3. update the email
	updateResp, err := sql_db.DbUpdateOneById[T](aus.ProvServ.GormConn, ctx, userId, models.UserDto{Email: &email}, nil)
	if err != nil {
		return dtos.BadReqC[bool](ICnst.InfoOrCode), ICnst.InfoOrCodeErr
	}
	//4. delete all the users session
	_, err = sql_db.DbDeleteMany[models.Session](aus.ProvServ.GormConn, ctx, models.Session{UserId: userId}, nil)
	if err != nil {
		return dtos.InternalErrMS[bool](err.Error()), err
	}
	//TODO: black list all the sessions here
	return dtos.SuccessS(true, updateResp.RowsAffected), nil
}
