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
	ADMIN = Role("ADMIN")
	USER  = Role("USER")
	OWNER = Role("OWNER")
)

//=================================   !  AccountStatus  ============================

type AccountStatus string

const (
	AccountPendingVerification = AccountStatus("pending_verification") //when it is pending email Verification
	AccountVerified            = AccountStatus("verified")             //email verified but needs to be accepted

	AccountActive = AccountStatus("active")

	AccountDisabled = AccountStatus("disabled") //when it is Disabled by the admin
	//old
	AccountDeleted = AccountStatus("deleted") //When the user deletes his own account
)

//=================================   !  CompanyStatus  ============================

type CompanyStatus string
const (
	CompanyPendingApproval = CompanyStatus("pending_approval")
	CompanyApproved        = CompanyStatus("approved")
	CompanyDisabled        = CompanyStatus("disabled")
)
