package common

type EnvVar string

const (
	Environment  = EnvVar("ENVIRONMENT")
	RefreshToken = "refresh-token"
	AccessToken  = "access-token"
)
const ApiV1 = "/api/v1"

type ContextKey string

var CtxClaims = ContextKey("USER_CLAIMS")
var CTXCompany_ID = ContextKey("CTX_COMPANY_ID")
var CTXUser_ID = ContextKey("CTX_USER_ID")

func (o ContextKey) Str() string {
	return string(o)
}

type EmailTemplatePaths string

const (
	VerificationTemplate = EmailTemplatePaths("")
)
