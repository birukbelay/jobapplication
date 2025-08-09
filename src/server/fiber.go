package server

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/app/account/auth"
	"github.com/projTemplate/goauth/src/app/account/profile"
	admins "github.com/projTemplate/goauth/src/app/admin/admin_users"
	"github.com/projTemplate/goauth/src/app/general/upload"
	"github.com/projTemplate/goauth/src/app/job/applications"
	"github.com/projTemplate/goauth/src/app/job/jobs"
	"github.com/projTemplate/goauth/src/models"
	conf "github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/providers"
)

func SetHumaCoreRoutes(humaRouter huma.API, dbs *providers.IProviderS, conf *conf.EnvConfig) {

	//core

	//account routes

	auth.SetupUserAuthRoutes(humaRouter, dbs, auth.NewAdminAuthServH[models.User](conf, dbs))
	//profile related routes

	profile.SetUserProfileRoutes(humaRouter, dbs)
	//platform owners management routes

	admins.SetupManageAdminUsersRoutes(humaRouter, dbs, admins.NewService(dbs))
	//company owners management routes

	upload.SetupUploadRoutes(humaRouter, dbs, upload.NewUploadGormServ(dbs))
	SetupJobModule(humaRouter, dbs)

	// //common routes
	// upload2.SetupUploadRoutes(humaRouter, cmnService, upload2.NewUploadGormServ(dbs.Gorm, cmnService))
}

// SetupJobModule initializes all job-related routes
func SetupJobModule(humaRouter huma.API, cmnServ *providers.IProviderS) {
	// Initialize services
	jobService := jobs.NewService(cmnServ)
	applicationService := applications.NewService(cmnServ)

	// Setup routes
	jobs.SetupJobRoutes(humaRouter, cmnServ, jobService)
	applications.SetupApplicationRoutes(humaRouter, cmnServ, applicationService)
}
