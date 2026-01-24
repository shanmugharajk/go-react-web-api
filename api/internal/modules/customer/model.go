package customer

import (
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
)

// Customer represents a customer in the system.
type Customer struct {
	ID      uint    `gorm:"primarykey" json:"id"`
	Name    string  `gorm:"not null" json:"name"`
	Email   string  `gorm:"unique;not null" json:"email"`
	Mobile  string  `gorm:"not null" json:"mobile"`
	Balance float64 `gorm:"default:0" json:"balance"`
	Active  bool    `gorm:"default:true" json:"active"`
	common.AuditFields
}

// CreateCustomerRequest represents the request payload for creating a customer.
type CreateCustomerRequest struct {
	Name   string `json:"name" validate:"required,min=1,max=255"`
	Email  string `json:"email" validate:"required,email"`
	Mobile string `json:"mobile" validate:"required,min=1,max=20"`
	Balance float64 `json:"balance" validate:"gte=0"`
}

// UpdateCustomerRequest represents the request payload for updating a customer.
// Only allows updating email, balance, and active status to maintain historical data integrity.
// Name and mobile are immutable to avoid confusion with historical records and transactions.
type UpdateCustomerRequest struct {
	Email   string  `json:"email" validate:"required,email"`
	Balance float64 `json:"balance" validate:"gte=0"`
	Active  bool    `json:"active"`
}
