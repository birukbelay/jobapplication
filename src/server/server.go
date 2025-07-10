package server

import (
	"fmt"
	"io/fs"
	"net/http"

	cmn "github.com/birukbelay/gocmn/src/logger"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"

	"github.com/projTemplate/goauth/src/app/account/auth"
	conf "github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/providers"
)

type MainServer struct {
	Engine  *gin.Engine
	EnvConf *conf.EnvConfig
}

var EmbeddedAssets fs.FS

func Create(dbs providers.IProviderS, conf *conf.EnvConfig) *MainServer {
	router := gin.Default()
	//router.Use(CORSMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PATCH", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Content-Length"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "http://localhost:3000"
		//},
		//MaxAge: 12 * time.Hour,
	}))
	gin.SetMode(gin.DebugMode)
	router.GET("/", func(c *gin.Context) {
		c.String(200, "holla")
	})
	serv := &MainServer{
		Engine:  router,
		EnvConf: conf,
	}
	config := huma.DefaultConfig("", "")
	config.DefaultFormat = "application/json"
	config.DocsPath = ""

	humaRouter := humagin.New(router, config)
	router.GET("/docs", ServDock)
	//huma.AutoRegister(humaRouter, As{})

	serv.SetupMiddleware()

	serv.SetHumaCoreRoutes(humaRouter, &dbs)

	return serv
}
func (s *MainServer) SetupMiddleware() {
	s.Engine.Use(gin.Logger())
	s.Engine.Use(gin.Recovery())
	s.Engine.Static("/assets", "./public/assets")
	s.Engine.StaticFS("/static", http.FS(EmbeddedAssets))
}

func (s *MainServer) Listen() error {
	s.Engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	cmn.LogTrace("server started at", fmt.Sprintf("http://127.0.0.1:%s/docs", s.EnvConf.ServerPort))
	err := s.Engine.Run(s.EnvConf.ServerHost + ":" + s.EnvConf.ServerPort)

	if err != nil {
		zlog.Panic().Err(err).Msg("listen Error")
	}
	return err
}

func (s *MainServer) SetHumaCoreRoutes(humaRouter huma.API, dbs *providers.IProviderS) {

	//core

	//account routes
	auth.SetupAuthRoutes(humaRouter, dbs, auth.NewAuthServH(s.EnvConf, dbs))
	// profile2.SetProfileRoutes(humaRouter, cmnService)

	// //common routes
	// upload2.SetupUploadRoutes(humaRouter, cmnService, upload2.NewUploadGormServ(dbs.Gorm, cmnService))
}
