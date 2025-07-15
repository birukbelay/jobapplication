package main

import (
	"embed"
	"io/fs"

	cmnConf "github.com/birukbelay/gocmn/src/config"
	// cloudinaryServ2 "github.com/Iscanner1/api/src/providers/upload/cloudinaryServ"

	"github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/models/migration"
	Idb "github.com/projTemplate/goauth/src/providers"
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
	// migration.MigrateDb2(Db)
	//Creating upload service
	//fileServ := diskUpload.NewDidkUploader(conf)
	// cloudinaryServ := cloudinaryServ2.NewCloudinaryUploader(conf)
	// providerService := providers.NewProviderServ(Db, conf, cloudinaryServ)
	emailSender := email.NewVerificationCodeSender(conf.SmtpHost, conf.SmtpPort, conf.SmtpPwd, conf.SmtpUsername)
	ginApp := IGin.CreateFiber(Idb.IProviderS{GormConn: Db, EnvConf: conf, VerificationCodeSender: emailSender}, conf)
	_ = ginApp.Listen()
}
