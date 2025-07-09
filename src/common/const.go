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

type OperationId string

func (o OperationId) Str() string {
	return string(o)
}
