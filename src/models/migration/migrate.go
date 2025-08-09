package migration

import (
	"log"

	"github.com/birukbelay/gocmn/src/logger"
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/models"
)

func MigrateDb2(Db *gorm.DB) {
	err := Db.AutoMigrate(&models.Upload{}, &models.User{}, &models.Job{}, &models.Application{})
	if err != nil {
		logger.LogTrace("failed to AutoMigrate Users:", err.Error())
		log.Panicln(err.Error())
		return
	}
	err = Db.AutoMigrate(&models.Session{}, &models.VerificationCode{})
	if err != nil {
		logger.LogTrace("failed to AutoMigrate session or Verification code:", err.Error())
		log.Panicln(err.Error())
		return
	}
}
