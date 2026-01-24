package customer

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// CustomerService handles business logic for customers.
type CustomerService struct {
	repo *CustomerRepository
}

// NewCustomerService creates a new CustomerService instance.
func NewCustomerService(repo *CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// GetAll retrieves all customers.
func (s *CustomerService) GetAll() ([]Customer, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a customer by ID.
func (s *CustomerService) GetByID(id uuid.UUID) (*Customer, error) {
	return s.repo.FindByID(id)
}

// Create creates a new customer.
func (s *CustomerService) Create(req CreateCustomerRequest, user *auth.User) (*Customer, error) {
	// Validate request
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	customer := &Customer{
		Name:    req.Name,
		Email:   req.Email,
		Mobile:  req.Mobile,
		Balance: req.Balance,
		Active:  true,
		AuditFields: common.AuditFields{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// Update updates an existing customer.
func (s *CustomerService) Update(id uuid.UUID, req UpdateCustomerRequest, user *auth.User) (*Customer, error) {
	// Validate request
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	customer.Email = req.Email
	customer.Balance = req.Balance
	customer.Active = req.Active
	customer.UpdatedBy = user.ID

	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// Delete soft deletes a customer.
func (s *CustomerService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
