package models

import (
	"github.com/projTemplate/goauth/src/models/enums"
)

type Job struct {
	Base   `mapstructure:",squash"`
	JobDto `mapstructure:",squash"`
	// Applications []Application `gorm:"foreignKey:JobID"`
	CreatedBy User `gorm:"foreignKey:CompanyID"`
}
type JobDto struct {
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Location    string          `json:"location,omitempty"`
	JobStatus   enums.JobStatus `json:"job_status,omitempty" enum:"Draft,Open,Closed"`

	//Relationships
	CompanyID string `json:"company_id,omitempty" gorm:"not null;index"`
}

func (d JobDto) SetOnCreate(key string) {
	d.CompanyID = key
}

type JobUpdateDto struct {
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Location    string          `json:"location,omitempty"`
	JobStatus   enums.JobStatus `json:"job_status,omitempty" enum:"Draft,Open,Closed"`
}
type JobFilter struct {
	ID        string          `query:"id"`
	Title     string          `query:"title"`
	CompanyID string          `query:"company_id"`
	Location  string          `query:"location"`
	JobStatus enums.JobStatus `query:"job_status" enum:"Draft,Open,Closed"`
}
type JobQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"title,description,location,job_status,company_id,id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"title,location,job_status,created_at,updated_at"`
}

func (q JobQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}

// ==============. Invite Codes

type Application struct {
	Base           `mapstructure:",squash"`
	ApplicationDto `mapstructure:",squash"`
	Job            Job    `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Applicant      User   `gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Status         string `json:"status,omitempty" gorm:"default:'Applied'"`
}
type ApplicationDto struct {
	ApplicantID string `json:"-" gorm:"not null;index"`
	CompanyID   string `json:"-" `
	JobID       string `json:"job_id,omitempty" gorm:"not null;index"`
	CoverLetter string `json:"cover_letter,omitempty"`
	ResumeUrl   string `json:"resume_url,omitempty"`
}

func (d ApplicationDto) SetOnCreate(key string) {
	d.ApplicantID = key
}

type StatusUpdateDto struct {
	Status string `json:"status,omitempty" enums:"Applied,Reviewed,Interview,Rejected,Hired"`
}

type ApplicationUpdateDto struct {
	CoverLetter string `json:"cover_letter,omitempty"`
	Resume      string `json:"resume,omitempty"`
}
type ApplicationFilter struct {
	ID          string `query:"id"`
	ApplicantID string `query:"applicant_id"`
	JobID       string `query:"job_id"`
	Status      string `query:"status"`
}
type ApplicationQuery struct {
	SelectedFields []string `query:"selected_fields" enum:"applicant_id,job_id,status,cover_letter,resume,id,created_at,updated_at"`
	Sort           string   `query:"sort" enum:"applicant_id,job_id,status,created_at,updated_at"`
}

func (q ApplicationQuery) GetQueries() (string, []string) {
	return q.Sort, q.SelectedFields
}
