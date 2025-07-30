package models

import (
	"github.com/birukbelay/gocmn/src/consts"
	"github.com/lib/pq"
)

//============ Operation access

type OperationAccess struct {
	OperationId  consts.OperationId `gorm:"primaryKey" json:"operation_id"`
	AllowedRoles pq.StringArray     `gorm:"type:text[]" json:"allowedRoles,omitempty"`
	GroupName    string
	Description  string
}
