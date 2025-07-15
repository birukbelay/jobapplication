package auth

import (
	ICrypt "github.com/birukbelay/gocmn/src/crypto"
)

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
