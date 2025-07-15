package models

import "github.com/lib/pq"

type UploadStatus string

const (
	Draft  = UploadStatus("DRAFT") //when it is pending email Verification
	ACTIVE = UploadStatus("active")
)

type Upload struct {
	Base        `mapstructure:",squash"`
	UploadDto   `mapstructure:",squash"`
	ChildImages []Upload `json:"child_images,omitempty" gorm:"foreignKey:ParentImgId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type UploadDto struct {
	Name string `gorm:"index:,unique;not null" json:"name,omitempty"` // the unique name fo the file
	//ImageFor    string ` json:"image_for,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Url         string `json:"url"`                // the url which we will get this file
	OwnerId     string `json:"owner_id,omitempty"` // the id of the user who uploaded the file
	Path        string `json:"-"`
	//Todo the files hash
	Hash     string       `json:"hash,omitempty"`      // the unique hash of the file
	FileType string       `json:"file_type,omitempty"` // the filetype got from the header
	Size     int64        `json:"size,omitempty"`      //the actual size of the file
	ImgModel string       `json:"img_model,omitempty"`
	Status   UploadStatus `json:"status,omitempty"`

	//=====  below are properties for images that are part of a group
	GroupId string         `json:"group_id,omitempty"` //if the images are part of array of images
	Images  pq.StringArray `json:"images,omitempty" gorm:"type:text[]"`
	//poly morphic associations
	ModelID   string `json:"model_id,omitempty"`
	ModelType string `json:"model_type,omitempty"`
	//parent images
	ParentImgId *string `json:"parent_img,omitempty"`
	IsPrimary   *bool   `json:"is_primary,omitempty"`
}
type UploadFilter struct {
	ID        string       `query:"id" query:"id"`
	Name      string       `query:"name,omitempty" query:"name"`
	Hash      string       `query:"hash" query:"hash"`
	OwnerId   string       `query:"owner_id,omitempty"`
	Status    UploadStatus `query:"status,omitempty"`
	ParentImg string       `query:"parent_img,omitempty"`
	IsPrimary bool         `query:"is_primary,omitempty"`
}
type UploadQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"name,display_name,url,user_id,path,hash,file_type,size,id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"name,path,file_type,size,id,created_at,updated_at"`
}
