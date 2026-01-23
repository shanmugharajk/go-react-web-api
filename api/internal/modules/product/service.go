package product

import (
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
)

// ProductService handles business logic for products.
type ProductService struct {
	repo *ProductRepository
}

// NewProductService creates a new product service.
func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAll retrieves all products.
func (s *ProductService) GetAll() ([]Product, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a product by ID.
func (s *ProductService) GetByID(id uint) (*Product, error) {
	return s.repo.FindByID(id)
}

// Create creates a new product.
func (s *ProductService) Create(req CreateProductRequest, user *auth.User) (*Product, error) {
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		AuditFields: common.AuditFields{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// Update updates an existing product.
func (s *ProductService) Update(id uint, req UpdateProductRequest, user *auth.User) (*Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.UpdatedBy = user.ID

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete deletes a product by ID.
func (s *ProductService) Delete(id uint, user *auth.User) error {
	return s.repo.Delete(id)
}

// CategoryService handles business logic for product categories.
type CategoryService struct {
	repo *CategoryRepository
}

// NewCategoryService creates a new product category service.
func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll retrieves all product categories.
func (s *CategoryService) GetAll() ([]ProductCategory, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a product category by ID.
func (s *CategoryService) GetByID(id uint) (*ProductCategory, error) {
	return s.repo.FindByID(id)
}

// Create creates a new product category.
func (s *CategoryService) Create(req CreateProductCategoryRequest, user *auth.User) (*ProductCategory, error) {
	category := &ProductCategory{
		Name:        req.Name,
		Description: req.Description,
		AuditFields: common.AuditFields{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Update updates an existing product category.
func (s *CategoryService) Update(id uint, req UpdateProductCategoryRequest, user *auth.User) (*ProductCategory, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	category.UpdatedBy = user.ID

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete deletes a product category by ID.
func (s *CategoryService) Delete(id uint, user *auth.User) error {
	return s.repo.Delete(id)
}
