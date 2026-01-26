package receiving

import (
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// StockReceipt represents a stock receipt for a purchase order.
type StockReceipt struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	PurchaseOrderID uuid.UUID `gorm:"type:char(36);index;not null" json:"purchaseOrderId"`
	ReceiptNumber   string    `gorm:"unique;not null" json:"receiptNumber"`
	ReceivedDate    time.Time `gorm:"not null" json:"receivedDate"`
	TotalAmount     float64   `gorm:"not null;default:0" json:"totalAmount"`
	Notes           string    `json:"notes"`
	common.AuditFields

	// Associations (not stored, for response)
	Items []StockReceiptItem `gorm:"foreignKey:StockReceiptID" json:"items,omitempty"`
}

// StockReceiptItem represents a line item in a stock receipt.
type StockReceiptItem struct {
	ID                  uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	StockReceiptID      uuid.UUID `gorm:"type:char(36);index;not null" json:"stockReceiptId"`
	PurchaseOrderItemID uuid.UUID `gorm:"type:char(36);index;not null" json:"purchaseOrderItemId"`
	ProductBatchID      uuid.UUID `gorm:"type:char(36);index;not null" json:"productBatchId"`
	QuantityReceived    int       `gorm:"not null" json:"quantityReceived"`
	common.AuditFields
}

// CreateStockReceiptRequest represents the request to create a stock receipt.
type CreateStockReceiptRequest struct {
	PurchaseOrderID uuid.UUID                       `json:"purchaseOrderId" validate:"required"`
	ReceivedDate    time.Time                       `json:"receivedDate" validate:"required"`
	Notes           string                          `json:"notes" validate:"max=1000"`
	Items           []CreateStockReceiptItemRequest `json:"items" validate:"required,min=1,dive"`
}

// CreateStockReceiptItemRequest represents a line item in create receipt request.
type CreateStockReceiptItemRequest struct {
	PurchaseOrderItemID uuid.UUID `json:"purchaseOrderItemId" validate:"required"`
	QuantityReceived    int       `json:"quantityReceived" validate:"required,gte=1"`
}
