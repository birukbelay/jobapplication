package sql_db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	// "gorm.io/driver/sqlserver"
	cmn "github.com/birukbelay/gocmn/src/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	conf "github.com/projTemplate/goauth/src/models/config"
)

func NewSqlDb(config *conf.SqlDbConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var confg gorm.Config
	var dsn string
	if os.Getenv("ENABLE_GORM_LOGGER") != "" {
		confg = gorm.Config{}
	} else {
		confg = gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	switch strings.ToLower(config.Driver) {
	// case "mysql":
	// 	dsn := config.Username + ":" + config.Password + "@tcp(" + config.MongoHost + ":" + strconv.Itoa(config.MongoPort) + ")/" + config.DbName + "?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=UTC"
	// 	db, err = gorm.Open(mysql.Open(dsn), &confg)
	// 	brea

	case "postgresql", "postgres":
		//dsn := "user=" + config.Username + " password=" + config.CoverImage + " dbname=" + config.DbName + " host=" + config.MongoHost + " port=" + strconv.Itoa(config.ServerPort) + " TimeZone=UTC"

		// dsn := fmt.Sprintf("user=postgres dbname=%s   sslmode=disable", config.DbName)
		// dsn := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable", config.Username, config.DbName, config.Password, config.Postgres)
		dsn = fmt.Sprintf("user=%s dbname=%s password=%s port=%s host=%s sslmode=%s", config.Username, config.DbName, config.Password, config.SqlPort, config.SqlHost, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &confg)

		// case "sqlserver", "mssql":
		// 	dsn := "sqlserver://" + config.Username + ":" + config.Password + "@" + config.MongoHost + ":" + strconv.Itoa(config.MongoPort) + "?database=" + config.DbName
		// 	db, err = gorm.Open(sqlserver.Open(dsn), &confg)
		// 	break
	}
	if err != nil || db == nil {
		cmn.LogTrace("failed to connect to database:", err.Error())
		cmn.LogTrace("conn string is :", dsn)
		log.Fatal(err)
		//panic("Failed to connect database")
	}

	cmn.LogTrace("Connection Opened to Postgress Database", config.DbName)

	return db, err
}
