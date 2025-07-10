package enums

import (
	"database/sql/driver"
	"fmt"
)

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
