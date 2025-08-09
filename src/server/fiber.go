package server

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/app/account/auth"
	"github.com/projTemplate/goauth/src/app/account/profile"
	admins "github.com/projTemplate/goauth/src/app/admin/admin_users"
	"github.com/projTemplate/goauth/src/app/admin/companies"
	"github.com/projTemplate/goauth/src/app/general/upload"
	"github.com/projTemplate/goauth/src/app/owner/company"
	"github.com/projTemplate/goauth/src/app/owner/inviteCode"
	"github.com/projTemplate/goauth/src/app/owner/user"
	"github.com/projTemplate/goauth/src/models"
	conf "github.com/projTemplate/goauth/src/models/config"
	"github.com/projTemplate/goauth/src/providers"
)

func SetHumaCoreRoutes(humaRouter huma.API, dbs *providers.IProviderS, conf *conf.EnvConfig) {

	//core

	//account routes
	auth.SetupAdminAuthRoutes(humaRouter, dbs, auth.NewAdminAuthServH[models.Admin](conf, dbs))
	auth.SetupUserAuthRoutes(humaRouter, dbs, auth.NewAdminAuthServH[models.User](conf, dbs))
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
	upload.SetupUploadRoutes(humaRouter, dbs, upload.NewUploadGormServ(dbs))

	// //common routes
	// upload2.SetupUploadRoutes(humaRouter, cmnService, upload2.NewUploadGormServ(dbs.Gorm, cmnService))
}
