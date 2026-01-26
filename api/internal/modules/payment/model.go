package payment

import (
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// PaymentMethod constants
const (
	PaymentMethodCash         = "cash"
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodCheque       = "cheque"
	PaymentMethodUPI          = "upi"
)

// VendorPayment represents a payment made to a vendor.
type VendorPayment struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	VendorID      uuid.UUID `gorm:"type:char(36);index;not null" json:"vendorId"`
	PaymentNumber string    `gorm:"unique;not null" json:"paymentNumber"`
	Amount        float64   `gorm:"not null" json:"amount"`
	PaymentDate   time.Time `gorm:"not null" json:"paymentDate"`
	PaymentMethod string    `gorm:"not null" json:"paymentMethod"`
	Reference     string    `json:"reference"`
	Notes         string    `json:"notes"`
	common.AuditFields
}

// CreateVendorPaymentRequest represents the request to create a vendor payment.
type CreateVendorPaymentRequest struct {
	VendorID      uuid.UUID `json:"vendorId" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	PaymentDate   time.Time `json:"paymentDate" validate:"required"`
	PaymentMethod string    `json:"paymentMethod" validate:"required,oneof=cash bank_transfer cheque upi"`
	Reference     string    `json:"reference" validate:"max=255"`
	Notes         string    `json:"notes" validate:"max=1000"`
}
