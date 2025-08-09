package main

import (
	cmnConf "github.com/birukbelay/gocmn/src/config"
	"github.com/birukbelay/gocmn/src/provider/db/redis"
	sql_db "github.com/birukbelay/gocmn/src/provider/db/sql"
	email "github.com/birukbelay/gocmn/src/provider/email/smtp"
	"github.com/birukbelay/gocmn/src/provider/upload/cloudinaryServ"
	"github.com/birukbelay/gocmn/src/server"

	"github.com/projTemplate/goauth/public/static/templates"
	"github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/models/migration"
	Idb "github.com/projTemplate/goauth/src/providers"
	locSrv "github.com/projTemplate/goauth/src/server"
)

func main() {

	//Loading the config
	conf := cmnConf.LoadConfigT[config.EnvConfig]()
	//Creating the database connection
	Db, _ := sql_db.NewSqlDb(&conf.SqlDbConfig)
	migration.MigrateDb2(Db)
	//Creating redis service
	redis, _ := redis.NewRedis(&conf.KeyValConfig)
	//email sending provider, also serves as verification code sender
	emailSender := email.NewSmtp(conf.SmtpHost, conf.SmtpPort, conf.SmtpPwd, conf.SmtpUsername, templates.Embedded)
	//Creating upload service
	cloudinary := cloudinaryServ.NewCloudinaryUploader(&conf.CloudinaryConfig) //file upload provider
	//createing the Provider, with, db, email, redis, fileUpload services
	provider := Idb.NewProvider(Db, conf, emailSender, emailSender, redis, cloudinary)
	//create the server with default configs
	srvr := server.CreateFiber(conf.ServerHost, conf.ServerPort)
	//Setting up the routes
	locSrv.SetHumaCoreRoutes(srvr.HumaRouter, provider, conf)
	_ = srvr.Listen()
}
