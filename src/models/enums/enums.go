package enums

import (
	"database/sql/driver"
	"fmt"
)

//=========================

//=================================   USER Roles  ============================

type Role string

func (r *Role) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*r = Role(v)
	case string:
		*r = Role(v)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}
func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

const (
	PLATFORM_ADMIN   = Role("PLATFORM_ADMIN")
	USER             = Role("USER")
	OWNER            = Role("OWNER")
	UNVERIFIED_OWNER = Role("UNVERIFIED_OWNER")
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

type CompanyStatus string

const (
	CompanyPendingApproval = CompanyStatus("pending_approval")
	CompanyApproved        = CompanyStatus("approved")
	CompanyDisabled        = CompanyStatus("disabled")
)
