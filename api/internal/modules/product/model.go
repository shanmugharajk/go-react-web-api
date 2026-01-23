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
	CategoryID  *uint  `json:"category_id"`
	common.AuditFields
}

// CreateProductRequest represents a request to create a product.
type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
	Stock       int    `json:"stock"`
	CategoryID  *uint  `json:"category_id"`
}

// UpdateProductRequest represents a request to update a product.
type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
	Stock       int    `json:"stock"`
	CategoryID  *uint  `json:"category_id"`
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
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateProductCategoryRequest represents a request to update a product category.
type UpdateProductCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
