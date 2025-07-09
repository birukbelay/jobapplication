package providers

import (
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/models/config"
)

type IProviderS struct {
	GormConn *gorm.DB
	EnvConf  *config.EnvConfig
	// UploadServ upload.FileUploadInterface
	//will have upload services, email services
}
