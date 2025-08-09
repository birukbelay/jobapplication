package auth

import (
	"context"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/provider/db/redis"
	"github.com/birukbelay/gocmn/src/resp_const"
	"github.com/birukbelay/gocmn/src/util"
	"github.com/mitchellh/mapstructure"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm/clause"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

type Service[T models.IntUsr] struct {
	Config   *config.EnvConfig
	Provider *providers.IProviderS
}

func NewAdminAuthServH[T models.IntUsr](conf *config.EnvConfig, genServ *providers.IProviderS) *Service[T] {
	return &Service[T]{
		Config:   conf,
		Provider: genServ,
	}
}

// RegisterUser (acc-01) , [AccountStatus], set(pwd,)
func (aus Service[T]) RegisterUser(ctx context.Context, input models.RegisterUserInput) (dtos.GResp[bool], error) {
	//Check if the email already exists
	usr, err := generic.DbGetOne[T](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Email}, nil)
	//if the user already exists
	logger.LogTrace("usr", usr)
	if err == nil {
		//TODO create a timeout
		//If the User is still pending verification, throw error
		if usr.RowsAffected > 0 && usr.Body.GetStatus() != enums.AccountPendingVerification {
			//FIXME Send password or email wrong depending on the scenario
			return dtos.BadReqC[bool](resp_const.UserExists), resp_const.UserExistError
		}
	}

	//user userDto here
	var userModel models.UserDto
	if err := mapstructure.Decode(input, &userModel); err != nil {
		return dtos.BadReqM[bool]("Decoding Input Error"), err
	}

	hash, err := crypto.BcryptCreateHash(input.Password)
	if err != nil {
		return dtos.InternalErrMS[bool]("Hashing Error"), err
	}
	userModel.Password = hash
	userModel.AccountStatus = enums.AccountPendingVerification
	userModel.Active = false

	user, err := generic.DbUpsertOneAllFields[T](aus.Provider.GormConn, ctx, &userModel, []clause.Column{{Name: "email"}}, &generic.Opt{Debug: true})
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	return aus.UTIL_SendVerification(ctx, input.Email, models.GetID(user.Body), models.SignupVerification)
}

// VerifyUser (acc-01), [id]x
func (aus Service[T]) VerifyUser(ctx context.Context, input VerificationInput) (dtos.GResp[bool], error) {
	//Check if the email already exists
	usr, err := generic.DbGetOne[T](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info, AccountStatus: enums.AccountPendingVerification}, nil)
	if err != nil {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//Validate the code
	codeValid := aus.VerifyCode(ctx, models.GetID(usr.Body), input.Code)
	if !codeValid {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//Update the users status
	user, err := generic.DbUpdateByFilter[T](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info}, models.UserDto{AccountStatus: enums.AccountVerified, Active: true}, nil)
	if err != nil {
		// cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	//TODO: invalidate the code

	return dtos.SuccessS(true, user.RowsAffected), nil
}

// Login (acc-03) [companyId, Role, Password]
func (aus Service[T]) Login(ctx context.Context, input LoginData) (dtos.GResp[TokenResponse], error) {
	//1. check the user exists, to login the user must be active: with status: verified, companySetup...
	usr, err := generic.DbGetOne[T](aus.Provider.GormConn, ctx, models.UserDto{Email: util.Ptr(input.LoginInfo), Active: true}, &generic.Opt{Debug: true})
	if err != nil || usr.RowsAffected < 1 {
		return dtos.BadReqC[TokenResponse](resp_const.EmailOrPassword), resp_const.EmailOrPasswordErr
	}

	//2. compare the password hash
	valid := crypto.BcryptPasswordsMatch(input.Password, usr.Body.GetPwd())
	if !valid {
		return dtos.BadReqC[TokenResponse](resp_const.EmailOrPassword), resp_const.EmailOrPasswordErr
	}

	sessionID := ulid.Make().String()
	tokens, err := aus.UTIL_MakeSession(ctx, sessionID, string(usr.Body.GetRole()), models.GetID(usr.Body), usr.Body.GetCompanyId())
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}

	return dtos.SuccessS(TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, usr.RowsAffected), nil

}

// ResetToken (acc-04): FIXME to be update with redis: [Role, id]
func (aus Service[T]) ResetToken(ctx context.Context, refreshToken string) (dtos.GResp[TokenResponse], error) {
	//1. validate the refresh token
	claims, ok, err := crypto.Valid(refreshToken, aus.Config.RefreshSecret)
	if !ok || (err != nil) {
		return dtos.BadReqC[TokenResponse](resp_const.InvalidToken), resp_const.InvalidTokenError
	}
	// 2. get the user
	usr, err := generic.DbGetOneByID[T](aus.Provider.GormConn, ctx, claims.UserId, nil)
	if err != nil {
		return dtos.BadReqC[TokenResponse](resp_const.DataNotFound), resp_const.UserNotFoundError
	}
	// 3. get the session
	session, err := generic.DbGetOne[models.Session](aus.Provider.GormConn, ctx, models.Session{UserId: claims.UserId, SessionId: claims.SessionId}, nil)
	if err != nil {
		return dtos.BadReqC[TokenResponse](resp_const.DataNotFound), resp_const.UserNotFoundError
	}

	// 4. validate the refresh token is the same
	valid := crypto.ArgonPasswordsMatch(refreshToken, session.Body.HashedRefresh)
	if !valid {
		return dtos.BadReqC[TokenResponse](resp_const.TokenDontMatch), resp_const.TokenDontMatchError
	}
	//5. update the sessions, incase the users role is changed we need to user the users new role
	tokens, err := aus.UTIL_MakeSession(ctx, claims.SessionId, string(usr.Body.GetRole()), models.GetID(usr.Body), claims.CompanyId)
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}

	return dtos.SuccessS(TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, 1), nil

}

// Logout [-]
func (aus Service[T]) Logout(ctx context.Context, refreshToken string) (dtos.GResp[bool], error) {

	//1. the jwt token
	claims, ok, err := crypto.Valid(refreshToken, aus.Config.RefreshSecret)
	if !ok || (err != nil) {
		return dtos.BadReqC[bool](resp_const.InvalidToken), resp_const.InvalidTokenError
	}
	// 2. get the session from the database
	session, err := generic.DbGetOne[models.Session](aus.Provider.GormConn, ctx, models.Session{SessionId: claims.SessionId}, nil)
	if err != nil {
		return dtos.BadReqC[bool](resp_const.DataNotFound), resp_const.UserNotFoundError
	}
	//3. verify the refresh token is the same
	valid := crypto.ArgonPasswordsMatch(refreshToken, session.Body.HashedRefresh)
	if !valid {
		return dtos.BadReqC[bool](resp_const.TokenDontMatch), resp_const.TokenDontMatchError
	}
	//4. delete the session
	resp, eror := generic.DbDeleteOneById[models.Session](aus.Provider.GormConn, ctx, session.Body.ID, nil)
	if eror != nil {
		return dtos.InternalErrMS[bool](eror.Error()), eror
	}
	err = redis.BlacklistSession(aus.Provider.KeyValServ, ctx, session.Body.ID)
	return dtos.SuccessS(true, resp.RowsAffected), nil
}

// ForgotPwd [ID]
func (aus Service[T]) ForgotPwd(ctx context.Context, input VerifyReqInput) (dtos.GResp[bool], error) {
	usr, err := generic.DbGetOne[T](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Email}, nil)

	if err != nil {
		return dtos.SuccessS(true, 0), nil
	}
	return aus.UTIL_SendVerification(ctx, input.Email, models.GetID(usr.Body), models.PasswordReset)
}

// ResetPwd [ID]
func (aus Service[T]) ResetPwd(ctx context.Context, input PwdResetInput) (dtos.GResp[bool], error) {
	tx := aus.Provider.GormConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return dtos.InternalErrMS[bool](err.Error()), err
	}

	//1. Get the user
	usr, err := generic.DbGetOne[T](tx, ctx, models.UserFilter{Email: input.Info}, nil)
	if err != nil {
		tx.Rollback()
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//2. Validate the code
	codeValid := aus.VerifyCode(ctx, models.GetID(usr.Body), input.Code)
	if !codeValid {
		tx.Rollback()
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//3. hash the new Password
	hash, err := crypto.BcryptCreateHash(input.NewPassword)
	if err != nil {
		tx.Rollback()
		return dtos.InternalErrMS[bool]("Hashing Error"), err
	}
	//4. Update the users passowrd
	user, err := generic.DbUpdateByFilter[T](tx, ctx, models.UserFilter{Email: input.Info}, models.UserDto{Password: hash}, nil)
	if err != nil {
		tx.Rollback()
		// cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	//5: invalidate the code
	_, err = generic.DbDeleteByFilter[models.VerificationCode](tx, ctx, models.VerificationCode{UserId: models.GetID(usr.Body)}, nil)
	if err != nil {
		tx.Rollback()
		return dtos.InternalErrMS[bool]("Deleting users Codes errors"), err
	}
	commit := tx.Commit()
	if commit.Error != nil {
		return dtos.InternalErrMS[bool](commit.Error.Error()), commit.Error
	}
	//6. black list all the users sessions
	sessions, err := generic.DbFetchManyWithOffset[models.Session](aus.Provider.GormConn, ctx, models.Session{UserId: models.GetID(usr.Body)}, dtos.PaginationInput{Limit: 10000}, nil)
	if err != nil {

	}
	for _, val := range sessions.Body {
		err = redis.BlacklistSession(aus.Provider.KeyValServ, ctx, val.SessionId)
	}

	return dtos.SuccessS(true, user.RowsAffected), nil
}
