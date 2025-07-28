package user

import (
	"context"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/resp_const"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm/clause"

	"github.com/projTemplate/goauth/src/common/functions"
	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
)

type UserCreateDto struct {
	FirstName string `json:"firsName,omitempty" minLength:"1"`
	LastName  string `json:"lastName,omitempty" `
	Email     string `json:"email,omitempty"`
	Avatar    string `json:"avatar,omitempty" `
	Active    bool   `json:"active,omitempty"`
	Password  string
	Role      enums.Role ` json:"role,omitempty" enum:"responder, client"`
}

func (aus Service) CreateUser(ctx context.Context, input UserCreateDto, companyId string) (dtos.GResp[bool], error) {

	//Check if the email already exists
	usr, err := generic.DbGetOne[models.User](aus.ProvServ.GormConn, ctx, models.UserFilter{Email: input.Email}, nil)
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
	//TODO: verify the Email
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

	tx := aus.ProvServ.GormConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return dtos.InternalErrMS[bool](err.Error()), err
	}
	// []string{"Password", "VerificationCodeHash", "VerificationCodeExpire", "AccountStatus", "Role", "Active"}
	user, err := generic.DbUpsertOneAllFields[models.User](tx, ctx, &userModel, []clause.Column{{Name: "email"}}, &generic.Opt{Debug: true})
	if err != nil {
		tx.Rollback()
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	resp, err := functions.UTIL_SendVerification(aus.ProvServ, tx, ctx, input.Email, models.GetID(user.Body), models.SignupVerification)
	if err != nil {
		tx.Rollback()
	}
	return resp, err
}
