package migration

import (
	"log"

	"github.com/birukbelay/gocmn/src/logger"
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/models"
)

func MigrateDb2(Db *gorm.DB) {
	err := Db.AutoMigrate(&models.Upload{}, &models.User{}, &models.Company{}, &models.Admin{})
	if err != nil {
		logger.LogTrace("failed to AutoMigrate Users:", err.Error())
		log.Panicln(err.Error())
		return
	}
}
