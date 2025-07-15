package providers

import (
	"net/http"
	"strings"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/util"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers/email"
)

type IProviderS struct {
	GormConn               *gorm.DB
	EmailSender            email.EmailSender
	VerificationCodeSender email.VerificationSender

	EnvConf *config.EnvConfig
	// UploadServ upload.FileUploadInterface
	//will have upload services, email services
}

// Authorization needs database service as well as configs
func (gs *IProviderS) Authorization(operationId common.OperationId, needsAuth bool, allowedRoles ...enums.Role) func(ctx huma.Context, next func(huma.Context)) {

	return func(ctx huma.Context, next func(huma.Context)) {
		if !needsAuth {
			next(ctx)
			return
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

		ctx = huma.WithValue(ctx, string(common.CtxClaims), claims)
		//=====================   Authorization ===================
		if len(allowedRoles) > 0 {
			if !util.ElementExists(enums.Role(claims.Role), allowedRoles...) {
				ctx.SetStatus(http.StatusUnauthorized)
				_, _ = ctx.BodyWriter().Write([]byte("Not Authorized"))
				return
			}
		}

		// Call the next middleware in the chain. This eventually calls the
		// operation handler as well.
		next(ctx)
	}
}
