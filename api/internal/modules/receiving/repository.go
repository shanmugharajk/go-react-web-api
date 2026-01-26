package receiving

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StockReceiptRepository handles database operations for stock receipts.
type StockReceiptRepository struct {
	db *gorm.DB
}

// NewStockReceiptRepository creates a new StockReceiptRepository instance.
func NewStockReceiptRepository(db *gorm.DB) *StockReceiptRepository {
	return &StockReceiptRepository{db: db}
}

// FindAll retrieves all stock receipts.
func (r *StockReceiptRepository) FindAll() ([]StockReceipt, error) {
	var receipts []StockReceipt
	if err := r.db.Preload("Items").Order("received_date DESC").Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

// FindByID retrieves a stock receipt by ID with items.
func (r *StockReceiptRepository) FindByID(id uuid.UUID) (*StockReceipt, error) {
	var receipt StockReceipt
	if err := r.db.Preload("Items").First(&receipt, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &receipt, nil
}

// FindByPurchaseOrderID retrieves all stock receipts for a purchase order.
func (r *StockReceiptRepository) FindByPurchaseOrderID(purchaseOrderID uuid.UUID) ([]StockReceipt, error) {
	var receipts []StockReceipt
	if err := r.db.Preload("Items").Where("purchase_order_id = ?", purchaseOrderID).Order("received_date DESC").Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

// Create inserts a new stock receipt (generates receipt number).
func (r *StockReceiptRepository) Create(receipt *StockReceipt) error {
	if receipt.ID == uuid.Nil {
		receipt.ID = uuid.New()
	}
	if receipt.ReceiptNumber == "" {
		receipt.ReceiptNumber = generateReceiptNumber(receipt.ReceivedDate)
	}
	return r.db.Create(receipt).Error
}

// CreateWithoutNumber inserts a new stock receipt (receipt number must be pre-set).
func (r *StockReceiptRepository) CreateWithoutNumber(receipt *StockReceipt) error {
	if receipt.ID == uuid.Nil {
		receipt.ID = uuid.New()
	}
	return r.db.Create(receipt).Error
}

// GetDB returns the underlying database for transaction support.
func (r *StockReceiptRepository) GetDB() *gorm.DB {
	return r.db
}

// generateReceiptNumber generates a unique receipt number.
func generateReceiptNumber(receivedDate time.Time) string {
	return fmt.Sprintf("SR-%s-%s", receivedDate.Format("20060102"), uuid.New().String()[:8])
}
