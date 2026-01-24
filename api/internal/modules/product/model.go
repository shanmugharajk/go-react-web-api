package product

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// Product represents a product in the system.
type Product struct {
	ID          uuid.UUID  `gorm:"type:char(36);primarykey" json:"id"`
	Name        string     `gorm:"unique;not null" json:"name"`
	Description string     `json:"description"`
	CategoryID  *uuid.UUID `gorm:"type:char(36)" json:"categoryId"`
	IsActive    bool       `gorm:"not null;default:true" json:"isActive"`
	common.AuditFields
}

// CreateProductRequest represents a request to create a product.
type CreateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	CategoryID  *uuid.UUID `json:"categoryId" validate:"omitempty"`
	IsActive    *bool      `json:"isActive"`
}

// UpdateProductRequest represents a request to update a product.
type UpdateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	CategoryID  *uuid.UUID `json:"categoryId" validate:"omitempty"`
	IsActive    *bool      `json:"isActive"`
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
