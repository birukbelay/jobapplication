package auth

import (
	"context"
	"errors"
	"time"

	ICrypt "github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/resp_const"
	"github.com/birukbelay/gocmn/src/util"
	"github.com/mitchellh/mapstructure"
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
func (aus Service) RegisterCompanyOwner(ctx context.Context, input RegisterClientInput) (dtos.GResp[models.Admin], error) {
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
		return dtos.BadReqM[models.Admin]("Decoding Input Error"), err
	}

	hash, err := ICrypt.BcryptCreateHash(input.Password)
	if err != nil {
		return dtos.InternalErrMS[models.Admin]("Hashing Error"), err
	}
	userModel.Password = hash
	userModel.Role = enums.OWNER
	userModel.AccountStatus = enums.AccountPendingVerification
	userModel.Active = false

	verificationCode := "0000"

	codeHash, err := ICrypt.BcryptCreateHash(verificationCode)
	if err != nil {
		return dtos.InternalErrMS[models.Admin]("Hashing Error"), err
	}
	userModel.VerificationCodeHash = codeHash
	//TODO change ptr with cmm/util
	userModel.VerificationCodeExpire = models.Ptr(time.Now().Add(time.Minute * 3))

	emailerr := aus.Provider.VerificationCodeSender.SendVerificationCode(input.Email, verificationCode)
	if emailerr != nil {
		return dtos.InternalErrMS[models.Admin]("Sending Email error"), emailerr
	}

	user, err := generic.DbUpsertOne[models.Admin](aus.Provider.GormConn, ctx, &userModel, []clause.Column{{Name: "email"}}, []string{"Password", "VerificationCodeHash", "VerificationCodeExpire", "AccountStatus", "Role", "Active"}, &generic.Opt{Debug: true})
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[models.Admin]("creating Error"), err
	}
	return dtos.SuccessS(user.Body, user.RowsAffected), nil
}

// RegisterCompanyOwner (acc-01)
func (aus Service) VerifyCompanyUser(ctx context.Context, input VerificationInput) (dtos.GResp[bool], error) {
	//Check if the email already exists
	usr, err := generic.DbGetOne[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info, AccountStatus: enums.AccountPendingVerification}, nil)
	//if the user already exists
	if err != nil {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	if usr.Body.VerificationCodeExpire.Before(time.Now()) {
		return dtos.BadReqC[bool](resp_const.InfoOrCode), resp_const.InfoOrCodeErr
	}
	valid := ICrypt.BcryptPasswordsMatch(input.Code, usr.Body.VerificationCodeHash)
	if !valid {
		return dtos.BadReqC[bool](resp_const.EmailOrPassword), resp_const.EmailOrPasswordErr
	}

	user, err := generic.DbUpdateByFilter[models.Admin](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Info}, models.Admin{AccountStatus: enums.AccountVerified, UserDto: models.UserDto{Active: true}}, nil)
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
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

	//	3. Generate auth Token of password
	tokens, err := aus.GenerateTokens(&ICrypt.CustomClaims{Role: string(usr.Body.Role), UserId: usr.Body.ID, CompanyId: userCompany})
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}

	//TODO 4. replace this with Sessions Table: that have info on date, device, ip and etc
	usr, err = generic.DbUpdateOneById[models.Admin](aus.Provider.GormConn, ctx, usr.Body.ID, &models.Admin{HashedRefresh: tokens.RefreshToken}, nil)
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}
	if usr.RowsAffected < 1 {
		return dtos.InternalErrMS[TokenResponse]("token not updated"), errors.New("token Not Updated")
	}

	return dtos.SuccessS(TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, usr.RowsAffected), nil

}

// ResetToken (acc-04)
func (aus Service) ResetToken(ctx context.Context, refreshToken string) (dtos.GResp[TokenResponse], error) {
	claims, ok, err := ICrypt.Valid(refreshToken, aus.Config.RefreshSecret)
	if !ok || (err != nil) {
		return dtos.BadReqC[TokenResponse](resp_const.InvalidToken), resp_const.InvalidTokenError
	}
	cmn.LogTrace("claim", claims.UserId)
	usr, err := generic.DbGetOneByID[models.Admin](aus.Provider.GormConn, ctx, claims.UserId, nil)
	if err != nil {
		return dtos.BadReqC[TokenResponse](resp_const.DataNotFound), resp_const.UserNotFoundError
	}

	// cmn.LogTrace("user is", usr.Body.HashedRefresh)
	if refreshToken != usr.Body.HashedRefresh {
		return dtos.BadReqC[TokenResponse](resp_const.TokenDontMatch), resp_const.TokenDontMatchError
	}

	tokens, err := aus.GenerateTokens(&ICrypt.CustomClaims{Role: string(usr.Body.Role), UserId: usr.Body.ID})
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}
	//TODO: we can put the refresh token on redis
	resp, eror := generic.DbUpdateOneById[models.Admin](aus.Provider.GormConn, ctx, usr.Body.ID, &models.Admin{HashedRefresh: tokens.RefreshToken}, nil)
	if eror != nil {
		return dtos.InternalErrMS[TokenResponse](eror.Error()), eror
	}
	if resp.RowsAffected < 1 {
		return dtos.InternalErrMS[TokenResponse]("token not updated"), errors.New("token Not Updated")
	}
	return dtos.SuccessS(TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, resp.RowsAffected), nil

}

func (aus Service) Logout(ctx context.Context, input RefreshTokenInput) (dtos.GResp[bool], error) {

	panic("not implemented")
}

func (aus Service) ForgotPwd(ctx context.Context, input VerificationInput) (dtos.GResp[bool], error) {
	panic("not implemented")
}

func (aus Service) ResetPwd(ctx context.Context, input VerificationInput) (dtos.GResp[bool], error) {
	panic("not implemented")
}
