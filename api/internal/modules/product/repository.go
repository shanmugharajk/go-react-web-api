package product

import (
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
)

// ProductRepository handles data access for products.
type ProductRepository struct {
	db *db.DB
}

// NewProductRepository creates a new product repository.
func NewProductRepository(database *db.DB) *ProductRepository {
	return &ProductRepository{db: database}
}

// FindAll retrieves all products.
func (r *ProductRepository) FindAll() ([]Product, error) {
	var products []Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// FindByID retrieves a product by ID.
func (r *ProductRepository) FindByID(id uint) (*Product, error) {
	var product Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Create creates a new product.
func (r *ProductRepository) Create(product *Product) error {
	return r.db.Create(product).Error
}

// Update updates an existing product.
func (r *ProductRepository) Update(product *Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product by ID.
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

// CategoryRepository handles data access for product categories.
type CategoryRepository struct {
	db *db.DB
}

// NewCategoryRepository creates a new product category repository.
func NewCategoryRepository(database *db.DB) *CategoryRepository {
	return &CategoryRepository{db: database}
}

// FindAll retrieves all product categories.
func (r *CategoryRepository) FindAll() ([]ProductCategory, error) {
	var categories []ProductCategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByID retrieves a product category by ID.
func (r *CategoryRepository) FindByID(id uint) (*ProductCategory, error) {
	var category ProductCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Create creates a new product category.
func (r *CategoryRepository) Create(category *ProductCategory) error {
	return r.db.Create(category).Error
}

// Update updates an existing product category.
func (r *CategoryRepository) Update(category *ProductCategory) error {
	return r.db.Save(category).Error
}

// Delete deletes a product category by ID.
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&ProductCategory{}, id).Error
}
