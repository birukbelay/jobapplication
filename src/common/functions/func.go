package functions

import (
	"context"
	"time"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
)

func UTIL_SendVerification(provServ *providers.IProviderS, tx *gorm.DB, ctx context.Context, email, userId string, purpose models.CodePurpose) (dtos.GResp[bool], error) {
	//TODO genereate code here
	verificationCode := "0000"

	codeHash, err := crypto.BcryptCreateHash(verificationCode)
	if err != nil {
		return dtos.InternalErrMS[bool]("Hashing Error"), err
	}
	emailerr := provServ.VerificationCodeSender.SendVerificationCode(email, verificationCode)
	if emailerr != nil {
		return dtos.InternalErrMS[bool]("Sending Email error"), emailerr
	}
	verificationResp, err := generic.DbUpsertOneAllFields[models.VerificationCode](tx, ctx, models.VerificationCode{
		ExpiresAt: util.Ptr(time.Now().Add(time.Minute * 3)),
		CodeHash:  codeHash,
		Purpose:   purpose,
		UserId:    userId,
	}, []clause.Column{{Name: "user_id"}}, nil)
	if err != nil {
		// logger.LogTrace("error crating", err)
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	return dtos.SuccessS(true, verificationResp.RowsAffected), nil
}
