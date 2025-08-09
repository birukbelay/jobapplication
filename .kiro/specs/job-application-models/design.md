# Design Document

## Overview

This design outlines the restructuring of the Job and Application models to support a proper job application system. The current Application model is designed for invite codes and needs to be completely redesigned to represent job applications with proper applicant tracking functionality.

## Architecture

The updated models will follow the existing codebase patterns:
- Use the `Base` struct for common fields (ID, CreatedAt, UpdatedAt)
- Separate DTO structs for data transfer and validation
- Proper GORM tags for database relationships
- JSON tags for API serialization
- Enum types for status fields

## Components and Interfaces

### Job Model Updates

The existing Job model needs minor corrections:
- Fix field naming inconsistencies (Title vs Name)
- Ensure proper JSON tags match field names
- Maintain existing relationships with User model

### Application Model Redesign

The Application model will be completely restructured:
- Remove invite code functionality
- Add job application specific fields
- Establish proper foreign key relationships
- Implement application status tracking

### New Application Status Enum

A new enum type `ApplicationStatus` will be added to track application states through the hiring process.

## Data Models

### Updated Job Model Structure

```go
type Job struct {
    Base    `mapstructure:",squash"`
    JobDto  `mapstructure:",squash"`
    Employees []User `gorm:"foreignKey:CompanyID"`
    CreatedBy User   `gorm:"foreignKey:OwnerID"`
    Applications []Application `gorm:"foreignKey:JobID"`
}

type JobDto struct {
    Title       string          `json:"title" gorm:"not null"`
    Description string          `json:"description"`
    Location    string          `json:"location"`
    JobStatus   enums.JobStatus `json:"status" gorm:"default:Draft"`
    OwnerID     string          `json:"owner_id" gorm:"not null;index"`
}
```

### New Application Model Structure

```go
type Application struct {
    Base           `mapstructure:",squash"`
    ApplicationDto `mapstructure:",squash"`
    Applicant      User `gorm:"foreignKey:ApplicantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Job            Job  `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ApplicationDto struct {
    ApplicantID  string                    `json:"applicant_id" gorm:"not null;index"`
    JobID        string                    `json:"job_id" gorm:"not null;index"`
    ResumeLink   string                    `json:"resume_link" gorm:"not null"`
    CoverLetter  string                    `json:"cover_letter"`
    Status       enums.ApplicationStatus   `json:"status" gorm:"default:Applied"`
    AppliedAt    *time.Time               `json:"applied_at" gorm:"autoCreateTime"`
}
```

### New ApplicationStatus Enum

```go
type ApplicationStatus string

const (
    APPLICATION_APPLIED    = ApplicationStatus("Applied")
    APPLICATION_REVIEWED   = ApplicationStatus("Reviewed")
    APPLICATION_INTERVIEW  = ApplicationStatus("Interview")
    APPLICATION_REJECTED   = ApplicationStatus("Rejected")
    APPLICATION_HIRED      = ApplicationStatus("Hired")
)
```

### Supporting DTOs

```go
type ApplicationUpdateDto struct {
    Status      enums.ApplicationStatus `json:"status"`
    CoverLetter string                  `json:"cover_letter"`
}

type ApplicationFilter struct {
    ID          string                    `query:"id"`
    ApplicantID string                    `query:"applicant_id"`
    JobID       string                    `query:"job_id"`
    Status      enums.ApplicationStatus   `query:"status"`
}

type ApplicationQuery struct {
    SelectedFields []string `query:"selected_fields" enum:"applicant_id,job_id,resume_link,cover_letter,status,applied_at,id,created_at,updated_at"`
    Sort           string   `query:"sort" enum:"applied_at,status,created_at,updated_at"`
}
```

## Error Handling

- Foreign key constraints will prevent orphaned applications
- Enum validation will ensure only valid status transitions
- Required field validation will ensure data integrity
- Proper GORM tags will handle database-level constraints

## Testing Strategy

- Unit tests for model validation
- Integration tests for database relationships
- Tests for enum value validation
- Tests for proper JSON serialization/deserialization
- Migration tests to ensure data integrity during updates