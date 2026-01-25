package product

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// Product represents a product in the system.
type Product struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`

	CategoryID uuid.UUID       `gorm:"type:char(36);index;not null" json:"categoryId"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`

	common.AuditFields
}

// CreateProductRequest represents a request to create a product.
type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
	Stock       int       `json:"stock" validate:"gte=0"`
	CategoryID  uuid.UUID `json:"categoryId" validate:"required"`
}

// UpdateProductRequest represents a request to update a product.
type UpdateProductRequest struct {
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"max=1000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
	Stock       int       `json:"stock" validate:"gte=0"`
	CategoryID  uuid.UUID `json:"categoryId" validate:"required"`
}

// ProductCategory represents a product category in the system.
type ProductCategory struct {
	ID          uuid.UUID `gorm:"type:char(36);primarykey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	common.AuditFields
}

// CreateProductCategoryRequest represents a request to create a product category.
type CreateProductCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateProductCategoryRequest represents a request to update a product category.
type UpdateProductCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}
