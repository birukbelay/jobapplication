package main

import (
	"embed"
	"io/fs"

	cmnConf "github.com/birukbelay/gocmn/src/config"
	// cloudinaryServ2 "github.com/Iscanner1/api/src/providers/upload/cloudinaryServ"
	"github.com/birukbelay/gocmn/src/provider/upload/cloudinaryServ"

	"github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/models/migration"
	Idb "github.com/projTemplate/goauth/src/providers"
	"github.com/projTemplate/goauth/src/providers/db/redis"
	sql_db "github.com/projTemplate/goauth/src/providers/db/sql_db"
	email "github.com/projTemplate/goauth/src/providers/email/smtp"
	IGin "github.com/projTemplate/goauth/src/server"
)

//go:embed public/static/styles
var EmbedAsset embed.FS

func main() {
	embeddedAssets, err := fs.Sub(EmbedAsset, "public/static/styles")
	if err != nil {
		panic(err)
	}
	IGin.EmbeddedAssets = embeddedAssets

	conf := cmnConf.LoadConfigT[config.EnvConfig]()
	Db, _ := sql_db.NewSqlDb(&conf.SqlDbConfig)
	migration.MigrateDb2(Db)
	redis, err := redis.NewRedis(&conf.KeyValConfig)

	//Creating upload service
	//fileServ := diskUpload.NewDidkUploader(conf)

	emailSender := email.NewSmtp(conf.SmtpHost, conf.SmtpPort, conf.SmtpPwd, conf.SmtpUsername) //email sending provider
	cloudinary := cloudinaryServ.NewCloudinaryUploader(&conf.CloudinaryConfig)                  //file upload provider
	provider := Idb.NewProvider(Db, conf, emailSender, emailSender, redis, cloudinary)
	ginApp := IGin.CreateFiber(provider, conf)
	_ = ginApp.Listen()
}
