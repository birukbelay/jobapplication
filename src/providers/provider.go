package providers

import (
	"net/http"
	"strings"

	"github.com/birukbelay/gocmn/src/consts"
	common "github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/provider/db"
	"github.com/birukbelay/gocmn/src/provider/email"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/util"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/models/config"
)

type IProviderS struct {
	GormConn               *gorm.DB
	EmailSender            email.EmailSender
	VerificationCodeSender email.VerificationSender
	KeyValServ             db.KeyValServ
	UploadServ             upload.FileUploadInterface
	EnvConf                *config.EnvConfig
	// UploadServ upload.FileUploadInterface
	//will have upload services, email services
}

func NewProvider(conn *gorm.DB, env *config.EnvConfig, emailSender email.EmailSender, verificationSender email.VerificationSender, keyValServ db.KeyValServ, uploadServ upload.FileUploadInterface) *IProviderS {

	return &IProviderS{
		GormConn:               conn,
		EnvConf:                env,
		EmailSender:            emailSender,
		VerificationCodeSender: verificationSender,
		KeyValServ:             keyValServ,
		UploadServ:             uploadServ,
	}
}

type AuthOpts struct {
	SkipAuth   bool
	SetCompany bool
	SetUser    bool
}

// Authorization needs database service as well as configs
func (gs *IProviderS) Authorization(operationId consts.OperationId, allowedRoles []string, opt *AuthOpts) func(ctx huma.Context, next func(huma.Context)) {

	return func(ctx huma.Context, next func(huma.Context)) {
		if opt != nil {
			if opt.SkipAuth {
				next(ctx)
				return
			}
		}

		//ctx.SetHeader("My-Custom-Header", "Hello, world!")
		token := ctx.Header("Authorization")
		if token == "" {
			ctx.SetStatus(http.StatusForbidden)
			_, _ = ctx.BodyWriter().Write([]byte("Token Is Empty"))
			return
		}
		substrings := strings.Split(token, " ")
		if len(substrings) != 2 {
			ctx.SetStatus(http.StatusForbidden)
			_, _ = ctx.BodyWriter().Write([]byte("Token Not Valid"))
			return
		}
		claims, ok, err := crypto.Valid(substrings[1], gs.EnvConf.AccessSecret)
		if !ok || (err != nil) {
			logger.LogTrace("err", err.Error())
			ctx.SetStatus(http.StatusForbidden)
			_, _ = ctx.BodyWriter().Write([]byte("Token Not Valid"))
			return
		}

		//cmn.LogTrace("Authorization header:", token)
		//cmn.LogTrace("needsAuth:", needsAuth)

		//cmn.LogTrace("EnvConf:", gs.EnvConf.AccessSecret)

		ctx = huma.WithValue(ctx, common.CtxClaims.Str(), claims)
		ctx = huma.WithValue(ctx, common.CTXCompany_ID.Str(), claims.CompanyId)
		ctx = huma.WithValue(ctx, common.CTXUser_ID.Str(), claims.UserId)
		//=====================   Authorization ===================
		if len(allowedRoles) > 0 {
			if !util.ElementExists(claims.Role, allowedRoles...) {
				ctx.SetStatus(http.StatusUnauthorized)
				_, _ = ctx.BodyWriter().Write([]byte("Not Authorized"))
				return
			}
		}
		//TODO: check for blacklisted session on redis

		// Call the next middleware in the chain. This eventually calls the
		// operation handler as well.
		next(ctx)
	}
}
