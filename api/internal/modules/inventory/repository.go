package inventory

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/errors"
)

// BatchRepository handles data access for product batches.
type BatchRepository struct {
	db *db.DB
}

// NewBatchRepository creates a new batch repository.
func NewBatchRepository(database *db.DB) *BatchRepository {
	return &BatchRepository{db: database}
}

// FindAll retrieves all product batches.
func (r *BatchRepository) FindAll() ([]ProductBatch, error) {
	var batches []ProductBatch
	if err := r.db.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

// FindByID retrieves a product batch by ID.
func (r *BatchRepository) FindByID(id uuid.UUID) (*ProductBatch, error) {
	var batch ProductBatch
	if err := r.db.First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

// FindByProductID retrieves all batches for a product.
func (r *BatchRepository) FindByProductID(productID uuid.UUID) ([]ProductBatch, error) {
	var batches []ProductBatch
	if err := r.db.Where("product_id = ?", productID).Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

// Create creates a new product batch.
func (r *BatchRepository) Create(batch *ProductBatch) error {
	if batch.ID == uuid.Nil {
		batch.ID = uuid.New()
	}
	return r.db.Create(batch).Error
}

// Update updates an existing product batch.
func (r *BatchRepository) Update(batch *ProductBatch) error {
	return r.db.Save(batch).Error
}

// Delete deletes a product batch by ID.
func (r *BatchRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&ProductBatch{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}
