package upload

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/birukbelay/gocmn/src/crypto"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/util"
	"github.com/danielgtaylor/huma/v2"

	"github.com/projTemplate/goauth/src/common"
	"github.com/projTemplate/goauth/src/models"
)

func (uph *GinUploadHandler) SingleFileUpload(ctx context.Context, input *struct {
	RawBody multipart.Form
	dtos.AuthParam
	//dtos.InputId
}) (*dtos.HumaResponse[dtos.GResp[*models.Upload]], error) {

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	files := input.RawBody.File["filename"]
	uploads, err := uph.UploadServ.SaveFileHeader(ctx, files[0], v.UserId, nil)
	if err != nil {
		return dtos.HumaReturnG(uploads, err)
	}
	return dtos.HumaReturnG(uploads, nil)
}

func (uph *GinUploadHandler) FileUpload(ctx context.Context, input *struct {
	RawBody huma.MultipartFormFiles[struct {
		MyFile huma.FormFile `form:"file" contentType:"text/plain" required:"true"`
		// SomeOtherFiles   []huma.FormFile `form:"other-files" contentType:"text/plain" required:"true"`
		// NoTagBindingFile huma.FormFile   `contentType:"text/plain"`
		// MyGreeting                string          `form:"greeting", minLength:"6"`
		SomeNumbers []int `form:"numbers"`
	}]
	//dtos.InputId
}) (*dtos.HumaResponse[dtos.GResp[*models.Upload]], error) {
	logger.LogTrace("file is", input.RawBody.Data().SomeNumbers)

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	uploads, err := uph.UploadServ.SaveFile(ctx, input.RawBody.Data().MyFile, v.UserId, nil)
	if err != nil {
		return dtos.HumaReturnG(uploads, err)
	}
	return dtos.HumaReturnG(uploads, nil)
}
func (uph *GinUploadHandler) OffsetPaginatedUploads(ctx context.Context, filter *struct {
	models.UploadFilter
	models.UploadQuery
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]models.Upload]], error) {
	filter.PaginationInput.Select = filter.SelectedFields
	filter.PaginationInput.SortBy = filter.Sort

	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := generic.DbFetchManyWithOffset[models.Upload](uph.CmnServ.GormConn, ctx, filter.UploadFilter, filter.PaginationInput, &generic.Opt{AuthKey: util.Ptr("user_id"), AuthVal: &v.UserId})
	return dtos.PHumaReturn(resp, err)
}
func (uph *GinUploadHandler) GetOneUploadByFilter(ctx context.Context, filter *models.UploadFilter) (*dtos.HumaResponse[dtos.GResp[models.Upload]], error) {
	v, ok := ctx.Value(common.CtxClaims.Str()).(crypto.CustomClaims)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := generic.DbGetOne[models.Upload](uph.CmnServ.GormConn, ctx, filter, &generic.Opt{AuthKey: util.Ptr("user_id"), AuthVal: &v.UserId})
	return dtos.HumaReturnG(resp, err)
}
