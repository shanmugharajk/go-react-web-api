package common

import "time"

type AuditFields struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy uint      `json:"createdBy"`
	UpdatedBy uint      `json:"updatedBy"`
}
