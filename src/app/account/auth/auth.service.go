package auth

import (
	"context"

	ICrypt "github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	cmn "github.com/birukbelay/gocmn/src/logger"
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

type Service struct {
	Config   *config.EnvConfig
	Provider *providers.IProviderS
}

func NewAuthServH(conf *config.EnvConfig, genServ *providers.IProviderS) *Service {
	return &Service{
		Config:   conf,
		Provider: genServ,
	}
}

// RegisterCompanyOwner (acc-01)
func (aus Service) RegisterCompanyOwner(ctx context.Context, input RegisterClientInput) (dtos.GResp[bool], error) {
	//Check if the email already exists
	usr, err := generic.DbGetOne[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Email}, nil)
	//if the user already exists
	logger.LogTrace("usr", usr)
	if err == nil {
		//TODO create a timeout
		//If the User is still pending verification, throw error
		if usr.RowsAffected > 0 && usr.Body.AccountStatus != enums.AccountPendingVerification {
			//FIXME Send password or email wrong depending on the scenario
			// return dtos.BadReqC[models.Admin](resp_const.UserExists), resp_const.UserExistError
		}
	}
	//TODO: verify the Email

	var userModel models.Admin
	if err := mapstructure.Decode(input, &userModel); err != nil {
		return dtos.BadReqM[bool]("Decoding Input Error"), err
	}

	hash, err := ICrypt.BcryptCreateHash(input.Password)
	if err != nil {
		return dtos.InternalErrMS[bool]("Hashing Error"), err
	}
	userModel.Password = hash
	userModel.Role = enums.OWNER
	userModel.AccountStatus = enums.AccountPendingVerification
	userModel.Active = false

	// []string{"Password", "VerificationCodeHash", "VerificationCodeExpire", "AccountStatus", "Role", "Active"}
	user, err := generic.DbUpsertOneAllFields[models.Admin](aus.Provider.GormConn, ctx, &userModel, []clause.Column{{Name: "email"}}, &generic.Opt{Debug: true})
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	return aus.UTIL_SendVerification(ctx, input.Email, user.Body.ID, models.Verification)
}

// RegisterCompanyOwner (acc-01)
func (aus Service) VerifyCompanyUser(ctx context.Context, input VerificationInput) (dtos.GResp[bool], error) {
	//Check if the email already exists
	usr, err := generic.DbGetOne[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info, AccountStatus: enums.AccountPendingVerification}, nil)
	if err != nil {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//Validate the code
	codeValid := aus.VerifyCode(ctx, usr.Body.ID, input.Code)
	if !codeValid {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//Update the users status
	user, err := generic.DbUpdateByFilter[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info}, models.Admin{AccountStatus: enums.AccountVerified, UserDto: models.UserDto{Active: true}}, nil)
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	//TODO: invalidate the code

	return dtos.SuccessS(true, user.RowsAffected), nil
}

// Login (acc-03)
func (aus Service) Login(ctx context.Context, input LoginData) (dtos.GResp[TokenResponse], error) {
	//1. check the user exists
	usr, err := generic.DbGetOne[models.Admin](aus.Provider.GormConn, ctx, models.UserDto{Email: util.Ptr(input.LoginInfo)}, &generic.Opt{Debug: true})
	if err != nil || usr.RowsAffected < 1 {
		return dtos.BadReqC[TokenResponse](resp_const.EmailOrPassword), resp_const.EmailOrPasswordErr
	}
	//if the current user is not active return user not active error
	if !usr.Body.Active {
		return dtos.BadReqC[TokenResponse](resp_const.UserNotActive), resp_const.UserNotActiveErr
	}
	userCompany := ""
	//if usr.Body.Role == enums.OWNER {
	//	userCompany = usr.Body.OwnedCompany.ID
	//} else if usr.Body.Role == enums.USER {
	//	userCompany = usr.Body.EmployerCompany.ID
	//}
	//2. compare the password hash
	valid := ICrypt.BcryptPasswordsMatch(input.Password, usr.Body.Password)
	if !valid {
		return dtos.BadReqC[TokenResponse](resp_const.EmailOrPassword), resp_const.EmailOrPasswordErr
	}

	sessionID := ulid.Make().String()
	tokens, err := aus.UTIL_MakeSession(ctx, sessionID, string(usr.Body.Role), usr.Body.ID, userCompany)
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}

	return dtos.SuccessS(TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, usr.RowsAffected), nil

}

// ResetToken (acc-04): FIXME to be update with redis
func (aus Service) ResetToken(ctx context.Context, refreshToken string) (dtos.GResp[TokenResponse], error) {
	//1. validate the refresh token
	claims, ok, err := ICrypt.Valid(refreshToken, aus.Config.RefreshSecret)
	if !ok || (err != nil) {
		return dtos.BadReqC[TokenResponse](resp_const.InvalidToken), resp_const.InvalidTokenError
	}
	// 2. get the user
	usr, err := generic.DbGetOneByID[models.Admin](aus.Provider.GormConn, ctx, claims.UserId, nil)
	if err != nil {
		return dtos.BadReqC[TokenResponse](resp_const.DataNotFound), resp_const.UserNotFoundError
	}
	// 3. get the session
	session, err := generic.DbGetOne[models.Session](aus.Provider.GormConn, ctx, models.Session{UserId: claims.UserId, SessionId: claims.SessionId}, nil)
	if err != nil {
		return dtos.BadReqC[TokenResponse](resp_const.DataNotFound), resp_const.UserNotFoundError
	}

	// 4. validate the refresh token is the same
	valid := ICrypt.ArgonPasswordsMatch(refreshToken, session.Body.HashedRefresh)
	if !valid {
		return dtos.BadReqC[TokenResponse](resp_const.TokenDontMatch), resp_const.TokenDontMatchError
	}
	//5. update the sessions, incase the users role is changed we need to user the users new role
	tokens, err := aus.UTIL_MakeSession(ctx, claims.SessionId, string(usr.Body.Role), usr.Body.ID, claims.CompanyId)
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}

	return dtos.SuccessS(TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, 1), nil

}

func (aus Service) Logout(ctx context.Context, refreshToken string) (dtos.GResp[bool], error) {

	//1. the jwt token
	claims, ok, err := ICrypt.Valid(refreshToken, aus.Config.RefreshSecret)
	if !ok || (err != nil) {
		return dtos.BadReqC[bool](resp_const.InvalidToken), resp_const.InvalidTokenError
	}
	// 2. get the session from the database
	session, err := generic.DbGetOne[models.Session](aus.Provider.GormConn, ctx, models.Session{SessionId: claims.SessionId}, nil)
	if err != nil {
		return dtos.BadReqC[bool](resp_const.DataNotFound), resp_const.UserNotFoundError
	}
	//3. verify the refresh token is the same
	valid := ICrypt.BcryptPasswordsMatch(refreshToken, session.Body.HashedRefresh)
	if !valid {
		return dtos.BadReqC[bool](resp_const.TokenDontMatch), resp_const.TokenDontMatchError
	}
	//4. delete the session
	resp, eror := generic.DbDeleteOneById[models.Session](aus.Provider.GormConn, ctx, session.Body.ID, nil)
	if eror != nil {
		return dtos.InternalErrMS[bool](eror.Error()), eror
	}
	return dtos.SuccessS(true, resp.RowsAffected), nil
}

func (aus Service) ForgotPwd(ctx context.Context, input VerificationInput) (dtos.GResp[bool], error) {
	usr, err := generic.DbGetOne[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info}, nil)

	if err != nil {
		logger.LogTrace("Could not find User", usr.Body.Email)
		return dtos.SuccessS(true, 0), nil
	}
	return aus.UTIL_SendVerification(ctx, input.Info, usr.Body.ID, models.PasswordReset)
}

func (aus Service) ResetPwd(ctx context.Context, input PwdResetInput) (dtos.GResp[bool], error) {
	//1. Get the users
	usr, err := generic.DbGetOne[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info, AccountStatus: enums.AccountPendingVerification}, nil)
	if err != nil {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//2. Validate the code
	codeValid := aus.VerifyCode(ctx, usr.Body.ID, input.Code)
	if !codeValid {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	//3. hash the new Password
	hash, err := ICrypt.BcryptCreateHash(input.NewPassword)
	if err != nil {
		return dtos.InternalErrMS[bool]("Hashing Error"), err
	}
	//4. Update the users passowrd
	user, err := generic.DbUpdateByFilter[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info}, models.Admin{Password: hash}, nil)
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	//TODO: invalidate the code

	return dtos.SuccessS(true, user.RowsAffected), nil
}
