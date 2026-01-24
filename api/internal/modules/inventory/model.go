package inventory

import (
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// ProductBatch represents a batch of products in inventory.
type ProductBatch struct {
	ID           uuid.UUID  `gorm:"type:char(36);primarykey" json:"id"`
	ProductID    uuid.UUID  `gorm:"type:char(36);not null;index" json:"productId"`
	CostPrice    float64    `gorm:"not null" json:"costPrice"`
	SellingPrice float64    `gorm:"not null" json:"sellingPrice"`
	Quantity     int        `gorm:"not null;default:0" json:"quantity"`
	PurchasedAt  time.Time  `gorm:"not null" json:"purchasedAt"`
	ExpiresAt    *time.Time `json:"expiresAt"`
	common.AuditFields
}

// CreateBatchRequest represents a request to create a product batch.
type CreateBatchRequest struct {
	ProductID    uuid.UUID  `json:"productId" validate:"required"`
	CostPrice    float64    `json:"costPrice" validate:"required,gte=0"`
	SellingPrice float64    `json:"sellingPrice" validate:"required,gte=0"`
	Quantity     int        `json:"quantity" validate:"required,gte=0"`
	PurchasedAt  time.Time  `json:"purchasedAt" validate:"required"`
	ExpiresAt    *time.Time `json:"expiresAt"`
}

// UpdateBatchRequest represents a request to update a product batch.
type UpdateBatchRequest struct {
	CostPrice    float64    `json:"costPrice" validate:"required,gte=0"`
	SellingPrice float64    `json:"sellingPrice" validate:"required,gte=0"`
	Quantity     int        `json:"quantity" validate:"required,gte=0"`
	PurchasedAt  time.Time  `json:"purchasedAt" validate:"required"`
	ExpiresAt    *time.Time `json:"expiresAt"`
}
