package vendor

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// Vendor represents a supplier/vendor in the system.
type Vendor struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name          string    `gorm:"unique;not null" json:"name"`
	ContactPerson string    `json:"contactPerson"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Address       string    `json:"address"`
	Balance       float64   `gorm:"default:0" json:"balance"`
	Active        bool      `gorm:"default:true" json:"active"`
	common.AuditFields
}

// CreateVendorRequest represents the request payload for creating a vendor.
type CreateVendorRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	ContactPerson string `json:"contactPerson" validate:"max=255"`
	Phone         string `json:"phone" validate:"max=20"`
	Email         string `json:"email" validate:"omitempty,email"`
	Address       string `json:"address" validate:"max=500"`
}

// UpdateVendorRequest represents the request payload for updating a vendor.
type UpdateVendorRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	ContactPerson string `json:"contactPerson" validate:"max=255"`
	Phone         string `json:"phone" validate:"max=20"`
	Email         string `json:"email" validate:"omitempty,email"`
	Address       string `json:"address" validate:"max=500"`
	Active        bool   `json:"active"`
}
