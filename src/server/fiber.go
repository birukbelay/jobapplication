package server

import (
	"fmt"
	"net/http"

	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	zlog "github.com/rs/zerolog/log"

	"github.com/projTemplate/goauth/src/app/account/auth"
	"github.com/projTemplate/goauth/src/app/account/profile"
	admins "github.com/projTemplate/goauth/src/app/admin/admin_users"
	"github.com/projTemplate/goauth/src/app/admin/companies"
	"github.com/projTemplate/goauth/src/app/owner/company"
	"github.com/projTemplate/goauth/src/app/owner/inviteCode"
	"github.com/projTemplate/goauth/src/app/owner/user"
	"github.com/projTemplate/goauth/src/models"
	conf "github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/providers"
)

type FiberServer struct {
	Engine  *fiber.App
	EnvConf *conf.EnvConfig
}

func CreateFiber(dbs providers.IProviderS, conf *conf.EnvConfig) *FiberServer {

	app := fiber.New()

	app.Get("/fiber", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	serv := &FiberServer{
		Engine:  app,
		EnvConf: conf,
	}
	//huma related path
	config := huma.DefaultConfig("", "")
	config.DefaultFormat = "application/json"
	config.DocsPath = "/"
	serv.SetupMiddleware()

	humaRouter := humafiber.NewWithGroup(app, app, config)
	app.Get("/docs", ServFiberDock)
	// v1 := app.Group("/api/v1")

	serv.SetHumaCoreRoutes(humaRouter, &dbs)
	for _, x := range app.GetRoutes() {
		fmt.Printf("=> %s:%s\n", x.Method, x.Path)
	}

	return serv
}
func (s *FiberServer) SetupMiddleware() {
	// s.Engine.Use(gin.Logger())
	// s.Engine.Use(gin.Recovery())
	s.Engine.Static("/assets", "./public/assets")
	s.Engine.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(EmbeddedAssets),
		// PathPrefix: "static",
	}))

	// s.Engine.StaticFS("/static", http.FS(EmbeddedAssets))
}

func (s *FiberServer) Listen() error {
	// s.Engine.GET("/ping", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })
	s.Engine.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	cmn.LogTrace("server started at", fmt.Sprintf("http://127.0.0.1:%s/docs", s.EnvConf.ServerPort))
	err := s.Engine.Listen(s.EnvConf.ServerHost + ":" + s.EnvConf.ServerPort)

	if err != nil {
		zlog.Panic().Err(err).Msg("listen Error")
	}
	return err
}

func (s *FiberServer) SetHumaCoreRoutes(humaRouter huma.API, dbs *providers.IProviderS) {

	//core

	//account routes
	auth.SetupAdminAuthRoutes(humaRouter, dbs, auth.NewAdminAuthServH[models.Admin](s.EnvConf, dbs))
	auth.SetupUserAuthRoutes(humaRouter, dbs, auth.NewAdminAuthServH[models.User](s.EnvConf, dbs))
	//profile related routes
	profile.SetAdminProfileRoutes(humaRouter, dbs)
	profile.SetUserProfileRoutes(humaRouter, dbs)
	//platform owners management routes
	companies.SetupManageCompaniesRoutes(humaRouter, dbs, companies.NewService(dbs))
	admins.SetupManageAdminUsersRoutes(humaRouter, dbs, admins.NewService(dbs))
	//company owners management routes
	company.SetupOwnerCompanyRoutes(humaRouter, dbs, company.NewService(dbs))
	inviteCode.SetupInviteCodeRoutes(humaRouter, dbs, inviteCode.NewService(dbs))
	user.SetupCompanyUserRoutes(humaRouter, dbs, user.NewService(dbs))
	// //Admin Routes

	// //common routes
	// upload2.SetupUploadRoutes(humaRouter, cmnService, upload2.NewUploadGormServ(dbs.Gorm, cmnService))
}
