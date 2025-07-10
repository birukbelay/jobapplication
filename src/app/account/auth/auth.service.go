package auth

import (
	"context"
	"errors"

	ICrypt "github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/resp_const"
	"github.com/mitchellh/mapstructure"

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

// Login (acc-03)
func (aus Service) Login(ctx context.Context, input LoginData) (dtos.GResp[TokenResponse], error) {
	//1. check the user exists
	usr, err := generic.DbGetOne[models.User](aus.Provider.GormConn, ctx, models.UserDto{Username: input.LoginInfo}, &generic.Opt{Debug: true})
	if err != nil || usr.RowsAffected < 1 {
		return dtos.BadReqC[TokenResponse](resp_const.EmailOrPassword), resp_const.EmailOrPasswordErr
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

	//TODO 4. Update The Users Refresh hash
	usr, err = generic.DbUpdateOneById[models.User](aus.Provider.GormConn, ctx, usr.Body.ID, &models.User{HashedRefresh: tokens.RefreshToken}, nil)
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}
	if usr.RowsAffected < 1 {
		return dtos.InternalErrMS[TokenResponse]("token not updated"), errors.New("token Not Updated")
	}

	return dtos.SuccessS[TokenResponse](TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, usr.RowsAffected), nil

}

func (aus Service) GenerateTokens(user *ICrypt.CustomClaims) (*AuthTokens, error) {
	claims := &ICrypt.CustomClaims{
		Role:      user.Role,
		UserId:    user.UserId,
		CompanyId: user.CompanyId,
	}
	accessToken, err := ICrypt.SignAccessToken(aus.Config.AccessSecret, 30, claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := ICrypt.SignRefreshToken(aus.Config.AccessSecret, 60*24*7, claims)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		accessToken, refreshToken,
	}, nil

}

// ResetToken (acc-04)
func (aus Service) ResetToken(ctx context.Context, refreshToken string) (dtos.GResp[TokenResponse], error) {
	claims, ok, err := ICrypt.Valid(refreshToken, aus.Config.RefreshSecret)
	if !ok || (err != nil) {
		return dtos.BadReqC[TokenResponse](resp_const.InvalidToken), resp_const.InvalidTokenError
	}
	cmn.LogTrace("claim", claims.UserId)
	usr, err := generic.DbGetOneByID[models.User](aus.Provider.GormConn, ctx, claims.UserId, nil)
	if err != nil {
		return dtos.BadReqC[TokenResponse](resp_const.DataNotFound), resp_const.UserNotFoundError
	}

	cmn.LogTrace("user is", usr.Body.HashedRefresh)
	if refreshToken != usr.Body.HashedRefresh {
		return dtos.BadReqC[TokenResponse](resp_const.TokenDontMatch), resp_const.TokenDontMatchError

	}

	tokens, err := aus.GenerateTokens(&ICrypt.CustomClaims{Role: string(usr.Body.Role), UserId: usr.Body.ID})
	if err != nil {
		return dtos.InternalErrMS[TokenResponse](err.Error()), err
	}
	resp, eror := generic.DbUpdateOneById[models.User](aus.Provider.GormConn, ctx, usr.Body.ID, &models.User{HashedRefresh: tokens.RefreshToken}, nil)
	if eror != nil {
		return dtos.InternalErrMS[TokenResponse](eror.Error()), eror
	}
	if resp.RowsAffected < 1 {
		return dtos.InternalErrMS[TokenResponse]("token not updated"), errors.New("token Not Updated")
	}
	return dtos.SuccessS[TokenResponse](TokenResponse{
		AuthTokens: *tokens,
		UserData:   usr.Body,
	}, resp.RowsAffected), nil

}

// RegisterCompanyOwner (acc-01)
func (aus Service) RegisterCompanyOwner(ctx context.Context, input RegisterClientInput) (dtos.GResp[models.User], error) {
	//Check if the email already exists
	usr, err := generic.DbGetOne[models.User](aus.Provider.GormConn, ctx, models.UserFilter{Email: input.Email}, nil)
	//if the user already exists
	if err == nil {
		/*if the user is active
		if usr.Body.Status==StatusActive
		*/
		//If the User is not active
		if usr.RowsAffected > 0 {
			//FIXME Send password or email wrong depending on the scenario
			return dtos.BadReqC[models.User](resp_const.UserExists), resp_const.UserExistError
		}

	}
	//TODO: verify the Email

	var userModel models.User
	if err := mapstructure.Decode(input, &userModel); err != nil {
		return dtos.BadReqM[models.User]("Decoding Input Error"), err
	}

	hash, err := ICrypt.BcryptCreateHash(input.Password)
	if err != nil {
		return dtos.InternalErrMS[models.User]("Hashing Error"), err
	}
	userModel.Password = hash
	userModel.Role = enums.OWNER
	// userModel.AccountStatus = enums.AccountPendingVerification
	userModel.Active = false

	user, err := generic.DbCreateOne[models.User](aus.Provider.GormConn, ctx, &userModel, nil)
	if err != nil {
		cmn.LogTrace("error crating", err)
		return dtos.InternalErrMS[models.User]("creating Error"), err
	}
	return dtos.SuccessS(user.Body, user.RowsAffected), nil
}
