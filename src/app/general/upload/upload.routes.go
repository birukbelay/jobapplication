package upload

import (
	"net/http"

	"github.com/birukbelay/gocmn/src/consts"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/models/enums"
	"github.com/projTemplate/goauth/src/providers"
)

// GinUploadHandler .
type GinUploadHandler struct {
	UploadServ *GormUploadServ
	CmnServ    *providers.IProviderS
}

// NewGinUploadHandler creates an upload handler from DomainService & Generic Gorm Service
func NewGinUploadHandler(repo *GormUploadServ, serv *providers.IProviderS) *GinUploadHandler {
	return &GinUploadHandler{UploadServ: repo, CmnServ: serv}
}

const (
	SingleFileUpload       = consts.OperationId("SingleFileUpload")
	UploadOne              = consts.OperationId("UploadOne")
	OffsetPaginatedUploads = consts.OperationId("OffsetPaginatedUploads")
	GetOneUploadByFilter   = consts.OperationId("UP-1-GetOneUploadByFilter")
)

var OperationMap = map[consts.OperationId]models.OperationAccess{

	SingleFileUpload:       {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: ".."},
	UploadOne:              {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: ".."},
	OffsetPaginatedUploads: {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: ".."},
	GetOneUploadByFilter:   {AllowedRoles: []string{enums.COMPANY.S(), enums.PLATFORM_ADMIN.S()}, Description: ".."},
}

func SetupUploadRoutes(humaRouter huma.API, cmnServ *providers.IProviderS, serv *GormUploadServ) {
	handler := NewGinUploadHandler(serv, cmnServ)
	tags := []string{"upload"}
	path := consts.ApiV1 + "/upload"

	huma.Register(humaRouter, huma.Operation{
		OperationID: SingleFileUpload.Str(),
		Method:      http.MethodPost,
		Path:        path,
		Middlewares: huma.Middlewares{cmnServ.Authorization(SingleFileUpload, OperationMap[SingleFileUpload].AllowedRoles, nil)},
		Tags:        tags}, handler.SingleFileUpload,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: UploadOne.Str(),
		Method:      http.MethodPost,
		Path:        path + "/one",
		// Middlewares: huma.Middlewares{cmnServ.Authorization(UploadOne, OperationMap[UploadOne].AllowedRoles, nil)},
		Tags: tags}, handler.FileUpload,
	)
	huma.Register(humaRouter, huma.Operation{
		OperationID: OffsetPaginatedUploads.Str(),
		Method:      http.MethodGet,
		Path:        path,
		Tags:        tags,
		Middlewares: huma.Middlewares{cmnServ.Authorization(OffsetPaginatedUploads, OperationMap[OffsetPaginatedUploads].AllowedRoles, nil)},
	}, handler.OffsetPaginatedUploads,
	)
	// huma.Register(humaRouter, huma.Operation{
	// 	OperationID: GetOneUploadByFilter.Str(),
	// 	Method:      http.MethodGet,
	// 	Path:        path + "/qu",
	// 	Tags:        tags,
	// 	Middlewares: huma.Middlewares{cmnServ.Authorization(GetOneUploadByFilter, OperationMap[GetOneUploadByFilter].AllowedRoles, nil)},
	// }, handler.GetOneUploadByFilter,
	// )

}
