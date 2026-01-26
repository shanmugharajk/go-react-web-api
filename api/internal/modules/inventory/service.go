package inventory

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// BatchService handles business logic for product batches.
type BatchService struct {
	repo *BatchRepository
}

// NewBatchService creates a new batch service.
func NewBatchService(repo *BatchRepository) *BatchService {
	return &BatchService{repo: repo}
}

// GetAll retrieves all product batches.
func (s *BatchService) GetAll() ([]ProductBatch, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a product batch by ID.
func (s *BatchService) GetByID(id uuid.UUID) (*ProductBatch, error) {
	return s.repo.FindByID(id)
}

// GetByProductID retrieves all batches for a product.
func (s *BatchService) GetByProductID(productID uuid.UUID) ([]ProductBatch, error) {
	return s.repo.FindByProductID(productID)
}

// Create creates a new product batch.
func (s *BatchService) Create(req CreateProductBatchRequest, user *auth.User) (*ProductBatch, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	batch := &ProductBatch{
		ProductID:         req.ProductID,
		CostPrice:         req.CostPrice,
		SellingPrice:      req.SellingPrice,
		QuantityAvailable: req.QuantityAvailable,
		PurchasedAt:       req.PurchasedAt,
		ExpiresAt:         req.ExpiresAt,
		AuditFields: common.AuditFields{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.repo.Create(batch); err != nil {
		return nil, err
	}

	return batch, nil
}

// Update updates an existing product batch.
func (s *BatchService) Update(id uuid.UUID, req UpdateProductBatchRequest, user *auth.User) (*ProductBatch, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	batch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	batch.CostPrice = req.CostPrice
	batch.SellingPrice = req.SellingPrice
	batch.QuantityAvailable = req.QuantityAvailable
	batch.PurchasedAt = req.PurchasedAt
	batch.ExpiresAt = req.ExpiresAt
	batch.UpdatedBy = user.ID

	if err := s.repo.Update(batch); err != nil {
		return nil, err
	}

	return batch, nil
}

// Delete deletes a product batch by ID.
func (s *BatchService) Delete(id uuid.UUID, user *auth.User) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	logger.Info("product batch deleted", "batch_id", id, "deleted_by", user.ID)
	return nil
}
