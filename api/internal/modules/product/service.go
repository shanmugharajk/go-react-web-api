package product

// Service handles business logic for products.
type Service struct {
	repo *Repository
}

// NewService creates a new product service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetAll retrieves all products.
func (s *Service) GetAll() ([]Product, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a product by ID.
func (s *Service) GetByID(id uint) (*Product, error) {
	return s.repo.FindByID(id)
}

// Create creates a new product.
func (s *Service) Create(req CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// Update updates an existing product.
func (s *Service) Update(id uint, req UpdateProductRequest) (*Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete deletes a product by ID.
func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}
