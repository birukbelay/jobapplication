package common

import "github.com/birukbelay/gocmn/src/provider/email"

const (
	RefreshToken = "refresh-token"
	AccessToken  = "access-token"
)

const (
	VerificationTemplate   = email.EmailTemplates("verification_code.html")
	UserInvitationTemplate = email.EmailTemplates("user_invitation.html")
)
