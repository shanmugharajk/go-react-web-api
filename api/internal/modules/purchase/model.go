package purchase

import (
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// PurchaseOrder status constants
const (
	StatusDraft     = "draft"
	StatusOrdered   = "ordered"
	StatusPartial   = "partial"
	StatusReceived  = "received"
	StatusCancelled = "cancelled"
)

// PaymentStatus constants
const (
	PaymentStatusUnpaid  = "unpaid"
	PaymentStatusPartial = "partial"
	PaymentStatusPaid    = "paid"
)

// PurchaseOrder represents a purchase order to a vendor.
type PurchaseOrder struct {
	ID            uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	VendorID      uuid.UUID  `gorm:"type:char(36);index;not null" json:"vendorId"`
	OrderNumber   string     `gorm:"unique;not null" json:"orderNumber"`
	OrderDate     time.Time  `gorm:"not null" json:"orderDate"`
	Status        string     `gorm:"not null;default:draft" json:"status"`
	TotalAmount   float64    `gorm:"not null;default:0" json:"totalAmount"`
	PaidAmount    float64    `gorm:"not null;default:0" json:"paidAmount"`
	PaymentStatus string     `gorm:"not null;default:unpaid" json:"paymentStatus"`
	LastPaymentAt *time.Time `json:"lastPaymentAt"`
	Notes         string     `json:"notes"`
	common.AuditFields

	// Associations (not stored, for response)
	Items []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
}

// PurchaseOrderItem represents a line item in a purchase order.
type PurchaseOrderItem struct {
	ID               uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	PurchaseOrderID  uuid.UUID  `gorm:"type:char(36);index;not null" json:"purchaseOrderId"`
	ProductID        uuid.UUID  `gorm:"type:char(36);index;not null" json:"productId"`
	QuantityOrdered  int        `gorm:"not null" json:"quantityOrdered"`
	QuantityReceived int        `gorm:"not null;default:0" json:"quantityReceived"`
	CostPrice        float64    `gorm:"not null" json:"costPrice"`
	SellingPrice     float64    `gorm:"not null" json:"sellingPrice"`
	ExpiresAt        *time.Time `json:"expiresAt"`
	common.AuditFields
}

// CreatePurchaseOrderRequest represents the request to create a purchase order.
type CreatePurchaseOrderRequest struct {
	VendorID  uuid.UUID                       `json:"vendorId" validate:"required"`
	OrderDate time.Time                       `json:"orderDate" validate:"required"`
	Notes     string                          `json:"notes" validate:"max=1000"`
	Items     []CreatePurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// CreatePurchaseOrderItemRequest represents a line item in create request.
type CreatePurchaseOrderItemRequest struct {
	ProductID       uuid.UUID  `json:"productId" validate:"required"`
	QuantityOrdered int        `json:"quantityOrdered" validate:"required,gte=1"`
	CostPrice       float64    `json:"costPrice" validate:"required,gte=0"`
	SellingPrice    float64    `json:"sellingPrice" validate:"required,gte=0"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

// UpdatePurchaseOrderRequest represents the request to update a purchase order.
type UpdatePurchaseOrderRequest struct {
	VendorID  uuid.UUID                       `json:"vendorId" validate:"required"`
	OrderDate time.Time                       `json:"orderDate" validate:"required"`
	Status    string                          `json:"status" validate:"required,oneof=draft ordered cancelled"`
	Notes     string                          `json:"notes" validate:"max=1000"`
	Items     []UpdatePurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// UpdatePurchaseOrderItemRequest represents a line item in update request.
type UpdatePurchaseOrderItemRequest struct {
	ID              uuid.UUID  `json:"id"`
	ProductID       uuid.UUID  `json:"productId" validate:"required"`
	QuantityOrdered int        `json:"quantityOrdered" validate:"required,gte=1"`
	CostPrice       float64    `json:"costPrice" validate:"required,gte=0"`
	SellingPrice    float64    `json:"sellingPrice" validate:"required,gte=0"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}
