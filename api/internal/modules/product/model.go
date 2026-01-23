package product

import (
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// Product represents a product in the system.
type Product struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Stock       int    `gorm:"not null;default:0" json:"stock"`
	CategoryID  *uint  `json:"categoryId"`
	common.AuditFields
}

// CreateProductRequest represents a request to create a product.
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	CategoryID  *uint   `json:"category_id" validate:"omitempty,gt=0"`
}

// UpdateProductRequest represents a request to update a product.
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description" validate:"max=1000"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	CategoryID  *uint   `json:"category_id" validate:"omitempty,gt=0"`
}

// ProductCategory represents a product category in the system.
type ProductCategory struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
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
