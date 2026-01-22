package product

import (
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
)

// Repository handles data access for products.
type Repository struct {
	db *db.DB
}

// NewRepository creates a new product repository.
func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

// FindAll retrieves all products.
func (r *Repository) FindAll() ([]Product, error) {
	var products []Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// FindByID retrieves a product by ID.
func (r *Repository) FindByID(id uint) (*Product, error) {
	var product Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Create creates a new product.
func (r *Repository) Create(product *Product) error {
	return r.db.Create(product).Error
}

// Update updates an existing product.
func (r *Repository) Update(product *Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product by ID.
func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}
