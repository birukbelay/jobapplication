package enums

//=========================

//=================================   USER Roles  ============================

type Role string

const (
	PLATFORM_ADMIN = Role("PLATFORM_ADMIN")
	APPLICANT      = Role("APPLICANT")
	COMPANY        = Role("COMPANY")
)

func (r Role) S() string {
	return string(r)
}

//========= User Roles

type UserRoles string

func (r UserRoles) S() string {
	return string(r)
}

//=================================   !  AccountStatus  ============================

type AccountStatus string

const (
	AccountPendingVerification = AccountStatus("pending_verification") //when it is pending email Verification
	AccountVerified            = AccountStatus("verified")             //email verified but needs to be accepted

	AccountActive = AccountStatus("active")

	AccountDisabled = AccountStatus("disabled") //when it is Disabled by the admin

	AccountDeleted = AccountStatus("deleted") //When the user deletes his own account
)

func (r AccountStatus) S() string {
	return string(r)
}

//=================================   !  CompanyStatus  ============================

type JobStatus string

const (
	STATUS_DRAFT  = JobStatus("Draft")
	STATUS_OPEN   = JobStatus("Open")
	STATUS_CLOSED = JobStatus("Closed")
)
