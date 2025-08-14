package auth

import (
	"context"
	"time"

	ICrypt "github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/util"
	"gorm.io/gorm/clause"

	"github.com/projTemplate/goauth/src/models"
)

func (aus Service[T]) GenerateTokens(user *ICrypt.CustomClaims) (*AuthTokens, error) {
	claims := &ICrypt.CustomClaims{
		Role:      user.Role,
		UserId:    user.UserId,
		CompanyId: user.CompanyId,
		SessionId: user.SessionId,
	}
	accessToken, err := ICrypt.SignAccessToken(aus.Config.AccessSecret, 30, claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := ICrypt.SignRefreshToken(aus.Config.RefreshSecret, 60*24*7, claims)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		accessToken, refreshToken,
	}, nil

}
func (aus Service[T]) UTIL_SendVerification(ctx context.Context, email, userId string, purpose models.CodePurpose) (dtos.GResp[bool], error) {
	//TODO genereate code here
	verificationCode := "0000"

	codeHash, err := ICrypt.BcryptCreateHash(verificationCode)
	if err != nil {
		return dtos.InternalErrMS[bool]("Hashing Error"), err
	}
	emailerr := aus.Provider.VerificationCodeSender.SendVerificationCode(email, verificationCode)
	if emailerr != nil {
		// return dtos.InternalErrMS[bool]("Sending Email error"), emailerr
	}
	verificationResp, err := generic.DbUpsertOneAllFields[models.VerificationCode](aus.Provider.GormConn, ctx, models.VerificationCode{
		ExpiresAt: util.Ptr(time.Now().Add(time.Minute * 60)),
		CodeHash:  codeHash,
		Purpose:   purpose,
		UserId:    userId,
		Email:     email,
	}, []clause.Column{{Name: "email"}}, &generic.Opt{Debug: true})
	if err != nil {
		logger.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	return dtos.SuccessS(true, verificationResp.RowsAffected), nil
}

// VerifyCode TODO: add reason of error, like code expires
func (aus Service[T]) VerifyCode(ctx context.Context, email string, code string) bool {
	codeModel, err := generic.DbGetOne[models.VerificationCode](aus.Provider.GormConn, ctx, models.VerificationCode{Email: email}, nil)
	if err != nil {
		return false
	}
	//if the expiry date is Passed
	if codeModel.Body.ExpiresAt.Before(time.Now()) {
		return false
	}
	valid := ICrypt.BcryptPasswordsMatch(code, codeModel.Body.CodeHash)
	if !valid {
		return false
	}
	return true
}
func (aus Service[T]) UTIL_MakeSession(ctx context.Context, sessionId, role, userId, companyId string) (*AuthTokens, error) {
	//	3. Generate auth Token of password
	tokens, err := aus.GenerateTokens(&ICrypt.CustomClaims{Role: role, UserId: userId, CompanyId: companyId, SessionId: sessionId})
	if err != nil {
		return nil, err
	}
	//4. hash the refresh token
	refreshHash, err := ICrypt.ArgonCreateHash(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	//4.Create a session or update previous's hashed_refresh
	_, err = generic.DbUpsertOneListedFields[models.Session](aus.Provider.GormConn, ctx, models.Session{UserId: userId, HashedRefresh: refreshHash, SessionId: sessionId}, []clause.Column{{Name: "session_id"}}, []string{"hashed_refresh"}, &generic.Opt{Debug: true})
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
