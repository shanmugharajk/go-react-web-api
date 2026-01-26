package inventory

import (
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// ProductBatch represents a batch of products in inventory.
type ProductBatch struct {
	ID                uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID         uuid.UUID  `gorm:"type:char(36);index;not null" json:"productId"`
	CostPrice         float64    `gorm:"not null" json:"costPrice"`
	SellingPrice      float64    `gorm:"not null" json:"sellingPrice"`
	QuantityAvailable int        `gorm:"not null;default:0" json:"quantityAvailable"`
	PurchasedAt       time.Time  `gorm:"not null" json:"purchasedAt"`
	ExpiresAt         *time.Time `json:"expiresAt"`

	common.AuditFields
}

// CreateProductBatchRequest represents a request to create a product batch.
type CreateProductBatchRequest struct {
	ProductID         uuid.UUID  `json:"productId" validate:"required"`
	CostPrice         float64    `json:"costPrice" validate:"required,gte=0"`
	SellingPrice      float64    `json:"sellingPrice" validate:"required,gte=0"`
	QuantityAvailable int        `json:"quantityAvailable" validate:"required,gte=0"`
	PurchasedAt       time.Time  `json:"purchasedAt" validate:"required"`
	ExpiresAt         *time.Time `json:"expiresAt"`
}

// UpdateProductBatchRequest represents a request to update a product batch.
type UpdateProductBatchRequest struct {
	CostPrice         float64    `json:"costPrice" validate:"required,gte=0"`
	SellingPrice      float64    `json:"sellingPrice" validate:"required,gte=0"`
	QuantityAvailable int        `json:"quantityAvailable" validate:"required,gte=0"`
	PurchasedAt       time.Time  `json:"purchasedAt" validate:"required"`
	ExpiresAt         *time.Time `json:"expiresAt"`
}
