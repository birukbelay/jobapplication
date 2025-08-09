package models

import (
	"github.com/projTemplate/goauth/src/models/enums"
)

type Job struct {
	Base   `mapstructure:",squash"`
	JobDto `mapstructure:",squash"`
	// Applications []Application `gorm:"foreignKey:JobID"`
	CreatedBy User `gorm:"foreignKey:CreatedBy"`
}
type JobDto struct {
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Location    string          `json:"location,omitempty"`
	JobStatus   enums.JobStatus `json:"job_status,omitempty" enum:"Draft,Open,Closed"`

	//Relationships
	CreatedBy string `json:"created_by,omitempty"`
}

func (d JobDto) SetOnCreate(key string) {
	d.CreatedBy = key
}

type JobUpdateDto struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}
type JobFilter struct {
	ID        string          `query:"id"`
	Title     string          `query:"title"`
	CreatedBy string          `query:"owner_id"`
	Location  string          `query:"location"`
	JobStatus enums.JobStatus `query:"job_status"`
}
type JobQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"title,description,location,job_status,owner_id,id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"title,location,job_status,created_at,updated_at"`
}

func (q JobQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}

// ==============. Invite Codes

type Application struct {
	Base           `mapstructure:",squash"`
	ApplicationDto `mapstructure:",squash"`
	Job            Job  `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Applicant      User `gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
type ApplicationDto struct {
	ApplicantID string `json:"applicant_id,omitempty" gorm:"not null;index"`
	JobID       string `json:"job_id,omitempty" gorm:"not null;index"`
	Status      string `json:"status,omitempty" gorm:"default:'pending'"`
	CoverLetter string `json:"cover_letter,omitempty"`
	Resume      string `json:"resume,omitempty"`
}

func (d ApplicationDto) SetOnCreate(key string) {
	d.ApplicantID = key
}

type ApplicationUpdateDto struct {
	Status      string `json:"status,omitempty"`
	CoverLetter string `json:"cover_letter,omitempty"`
	Resume      string `json:"resume,omitempty"`
}
type ApplicationFilter struct {
	ID          string `query:"id,omitempty"`
	ApplicantID string `query:"applicant_id,omitempty"`
	JobID       string `query:"job_id,omitempty"`
	Status      string `query:"status,omitempty"`
}
type ApplicationQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"applicant_id,job_id,status,cover_letter,resume,id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"applicant_id,job_id,status,created_at,updated_at"`
}

func (q ApplicationQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}
