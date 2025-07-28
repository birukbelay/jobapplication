package inviteCode

import (
	"context"
	"errors"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/resp_const"
	"github.com/birukbelay/gocmn/src/util"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers/email"
)

//Tasks
/*
- [x] just create the invitation code(this version)
- create and send one invitation code, to multiple users
- create and send unique invitation codes to multiple inviteCodes
- join via invitation code
*/

func (aus Service) CreateAndSendInviteCode(ctx context.Context, to, companyId string, invitation models.InviteCodeDto) (dtos.GResp[models.InviteCode], error) {
	//TODO genereate code here

	invitation.CompanyID = companyId
	invitation.Active = true
	//todo: make a unique generator with 10 digit(4-10)
	invitationCode := util.GenerateRandomString(10)
	//TODO: hash or encrypt, the invitation
	// codeHash, err := crypto.BcryptCreateHash(invitationCode)
	// if err != nil {
	// 	return dtos.InternalErrMS[bool]("Hashing Error"), err
	// }
	invitation.Code = invitationCode

	emailStruct := models.InviteCodeDto{
		Code: invitationCode,
	}

	emailerr := aus.ProvServ.EmailSender.SendEmail(to, "Invitation to Join Company", email.UserInvitationTemplate, emailStruct)
	if emailerr != nil {
		return dtos.InternalErrMS[models.InviteCode]("Sending Email error"), emailerr
	}
	invitationResp, err := generic.DbCreateOne[models.InviteCode](aus.ProvServ.GormConn, ctx, invitation, nil)
	if err != nil {
		logger.LogTrace("error crating", err)
		return dtos.InternalErrMS[models.InviteCode]("creating Error"), err
	}
	return dtos.SuccessS(invitationResp.Body, invitationResp.RowsAffected), nil
}

// JoinViaInviteCode when a logged in user joins a company invite code
func (aus Service) JoinViaInviteCode(ctx context.Context, userId, code string) (dtos.GResp[bool], error) {
	tx := aus.ProvServ.GormConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return dtos.InternalErrMS[bool](err.Error()), err
	}
	// 1. getch the code from the database
	inviteCode, err := generic.DbGetOne[models.InviteCode](tx, ctx, models.InviteCodeFilter{Code: code, Active: true}, nil)
	if err != nil {
		tx.Rollback()
		return dtos.NotFoundErrS[bool](err.Error()), err
	}
	//check different conditions
	if inviteCode.Body.MaxUsage != nil {
		if inviteCode.Body.UsageCount >= *inviteCode.Body.MaxUsage {
			tx.Rollback()
			return dtos.BadReqM[bool]("this invite code has reached max usage"), errors.New("this invite code has reached max usage")
		}
	}
	userResp, err := generic.DbGetOne[models.User](tx, ctx, models.UserFilter{ID: userId}, nil)
	if err != nil {
		tx.Rollback()
		return dtos.NotFoundErrS[bool](err.Error()), err
	}
	if inviteCode.Body.UserInfo != nil {
		if userResp.Body.Email != inviteCode.Body.UserInfo {
			return dtos.NotFoundErrS[bool](resp_const.DataNotFound.Msg()), resp_const.DataNotFoundError
		}
	}
	updateResp, err := generic.DbUpdateOneById[models.User](tx, ctx, userId, models.UserDto{CompanyID: &inviteCode.Body.CompanyID, Role: enums.Role(inviteCode.Body.UserRole)}, nil)
	if err != nil {
		tx.Rollback()
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	_, err = generic.DbUpdateOneById[models.InviteCode](tx, ctx, userId, models.InviteCodeDto{UsageCount: inviteCode.Body.UsageCount + 1}, nil)
	if err != nil {
		tx.Rollback()
		return dtos.InternalErrMS[bool]("creating Error"), err
	}
	commit := tx.Commit()
	if commit.Error != nil {
		return dtos.InternalErrMS[bool](commit.Error.Error()), commit.Error
	}
	return dtos.SuccessS(true, updateResp.RowsAffected), nil
}

/*

later tasks

1. invitation signup: when user signups if he has an invitation code he will be associated with that company from the start
2. multiple user invitations: send invitations to multiple emails
3. are users defined: responder, client, ..., or are they dynamically created by the company

*/
