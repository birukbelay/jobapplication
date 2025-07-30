package upload

import (
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/generic"
	"github.com/birukbelay/gocmn/src/resp_const"
	"github.com/birukbelay/gocmn/src/util"
	"gorm.io/gorm"

	"github.com/projTemplate/goauth/src/models"
	"github.com/projTemplate/goauth/src/providers"
)

// GormUploadServ implements the item.ItemRepository interface
type GormUploadServ struct {
	conn    *gorm.DB
	CmnServ *providers.IProviderS
}

func NewUploadGormServ(cmnServ *providers.IProviderS) *GormUploadServ {
	if cmnServ == nil {
		log.Panic("provider is nill")
	}
	if cmnServ.GormConn == nil {
		log.Panic("gorm Connection Not FoundPlease make Sure connection Exists")
	}
	return &GormUploadServ{conn: cmnServ.GormConn, CmnServ: cmnServ}
}

type FileOpt struct {
	ModelId   *string
	ModelType *string
}

func (ups *GormUploadServ) SaveSingleFile(ctx context.Context, file *multipart.FileHeader, ownerId string, fileOPt *FileOpt) (dtos.GResp[*models.Upload], error) {

	fileResp, err := ups.CmnServ.UploadServ.UploadSingleFile(file)
	if err != nil {
		return dtos.RespStatusMsgS[*models.Upload](fileResp.Error, fileResp.Status), err
	}
	uploadDtos, err := util.MarshalToStruct[models.UploadDto](fileResp.Body)
	if err != nil {
		return dtos.RespStatusMsgS[*models.Upload](fileResp.Error, fileResp.Status), err
	}
	//TODO: add the ownerId, groupId & prefix
	uploadDtos.UserID = ownerId
	if fileOPt != nil {
		if fileOPt.ModelId != nil && fileOPt.ModelType != nil {
			uploadDtos.ModelID = *fileOPt.ModelId
			uploadDtos.ModelType = *fileOPt.ModelType
		}

	}

	upld, err := generic.DbCreateOne[models.Upload](ups.conn, ctx, uploadDtos, nil)
	if err != nil {
		return dtos.InternalErrMS[*models.Upload](err.Error()), err
	}
	return dtos.SuccessS(&upld.Body, upld.RowsAffected), nil
}
func (ups *GormUploadServ) removeFiles(c context.Context, filesToRemove []string) (dtos.GResp[[]models.Upload], error) {
	if len(filesToRemove) < 1 {
		return dtos.InternalErrMS[[]models.Upload]("file names can't be empty"), errors.New("file names can't be empty")
	}
	uploadTobeRemoved, err := generic.DbFetchWihtIn[models.Upload](ups.conn, c, "name", filesToRemove, nil)
	if err != nil {
		return dtos.InternalErrMS[[]models.Upload](" the files to be removed are not found "), err
	}
	var deletedFiles []models.Upload
	for _, value := range uploadTobeRemoved.Body {
		eror := ups.CmnServ.UploadServ.DeleteFileWithName(value.Name)
		if eror != nil {
			return dtos.InternalErrMS[[]models.Upload]("Failed To remove the file: " + value.Name), eror
		}
		// Remove the file From database
		upld, err := generic.DbDeleteByFilter[models.Upload](ups.conn, c, models.UploadFilter{Name: value.Name}, nil)
		if err != nil {
			return dtos.InternalErrMS[[]models.Upload]("Failed To remove the file from Database: " + value.Name), err
		}
		deletedFiles = append(deletedFiles, upld.Body)
	}

	return dtos.SuccessCS(deletedFiles, resp_const.DeleteSuccess, int64(len(deletedFiles))), nil
}
