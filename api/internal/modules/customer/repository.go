package customer

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
)

// CustomerRepository handles database operations for customers.
type CustomerRepository struct {
	db *db.DB
}

// NewCustomerRepository creates a new CustomerRepository instance.
func NewCustomerRepository(database *db.DB) *CustomerRepository {
	return &CustomerRepository{db: database}
}

// FindAll retrieves all customers.
func (r *CustomerRepository) FindAll() ([]Customer, error) {
	var customers []Customer
	if err := r.db.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// FindByID retrieves a customer by ID.
func (r *CustomerRepository) FindByID(id uuid.UUID) (*Customer, error) {
	var customer Customer
	if err := r.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// Create inserts a new customer into the database.
func (r *CustomerRepository) Create(customer *Customer) error {
	if customer.ID == uuid.Nil {
		customer.ID = uuid.New()
	}
	return r.db.Create(customer).Error
}

// Update updates an existing customer in the database.
func (r *CustomerRepository) Update(customer *Customer) error {
	return r.db.Save(customer).Error
}

// Delete soft deletes a customer by setting active to false.
func (r *CustomerRepository) Delete(id uuid.UUID) error {
	return r.db.Model(&Customer{}).Where("id = ?", id).Update("active", false).Error
}
