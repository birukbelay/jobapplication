# Implementation Plan

- [ ] 1. Add ApplicationStatus enum to enums package
  - Create new ApplicationStatus enum type with proper constants (Applied, Reviewed, Interview, Rejected, Hired)
  - Add String() method for enum conversion
  - _Requirements: 2.2, 3.4_

- [ ] 2. Update Job model structure and tags
  - Fix JobDto field naming inconsistency (Title field should have proper JSON tag)
  - Ensure all GORM tags are properly set for database constraints
  - Add Applications relationship to Job model for reverse lookup
  - _Requirements: 3.1, 3.3, 4.1, 4.4_

- [ ] 3. Completely restructure Application model
  - Replace existing Application struct with job application focused structure
  - Implement ApplicationDto with required fields (ApplicantID, JobID, ResumeLink, CoverLetter, Status, AppliedAt)
  - Add proper GORM tags for foreign key relationships and constraints
  - Add proper JSON tags for API serialization
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 3.1, 3.2, 3.3, 4.1, 4.2_

- [ ] 4. Update Application supporting DTOs
  - Rewrite ApplicationUpdateDto to support status updates and cover letter changes
  - Update ApplicationFilter to support filtering by applicant, job, and status
  - Update ApplicationQuery to include relevant fields for job application queries
  - _Requirements: 2.1, 2.2, 2.3, 4.2, 4.4_

- [ ] 5. Remove obsolete Application model methods
  - Remove SetOnCreate method that was specific to invite code functionality
  - Clean up any invite code related fields or methods
  - _Requirements: 4.2, 4.4_

- [ ] 6. Update database migration file
  - Add migration logic to handle the Application model restructure
  - Ensure proper foreign key constraints are created
  - Handle any existing data migration if needed
  - _Requirements: 3.3, 3.4_