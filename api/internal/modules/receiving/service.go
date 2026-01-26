package receiving

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/inventory"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/purchase"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/vendor"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
	"gorm.io/gorm"
)

var (
	ErrPurchaseOrderNotFound     = errors.New("purchase order not found")
	ErrPurchaseOrderItemNotFound = errors.New("purchase order item not found")
	ErrQuantityExceedsOrdered    = errors.New("quantity received exceeds quantity ordered")
	ErrOrderNotReceivable        = errors.New("purchase order is not in a receivable status")
)

// StockReceiptService handles business logic for stock receipts.
type StockReceiptService struct {
	repo       *StockReceiptRepository
	poRepo     *purchase.PurchaseOrderRepository
	vendorRepo *vendor.VendorRepository
	batchRepo  *inventory.BatchRepository
	db         *gorm.DB
}

// NewStockReceiptService creates a new StockReceiptService instance.
func NewStockReceiptService(
	repo *StockReceiptRepository,
	poRepo *purchase.PurchaseOrderRepository,
	vendorRepo *vendor.VendorRepository,
	batchRepo *inventory.BatchRepository,
	db *gorm.DB,
) *StockReceiptService {
	return &StockReceiptService{
		repo:       repo,
		poRepo:     poRepo,
		vendorRepo: vendorRepo,
		batchRepo:  batchRepo,
		db:         db,
	}
}

// GetAll retrieves all stock receipts.
func (s *StockReceiptService) GetAll() ([]StockReceipt, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a stock receipt by ID.
func (s *StockReceiptService) GetByID(id uuid.UUID) (*StockReceipt, error) {
	return s.repo.FindByID(id)
}

// GetByPurchaseOrderID retrieves all stock receipts for a purchase order.
func (s *StockReceiptService) GetByPurchaseOrderID(purchaseOrderID uuid.UUID) ([]StockReceipt, error) {
	return s.repo.FindByPurchaseOrderID(purchaseOrderID)
}

// Create creates a new stock receipt, creates product batches, and updates balances.
func (s *StockReceiptService) Create(req CreateStockReceiptRequest, user *auth.User) (*StockReceipt, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	// Get purchase order (read operation - outside transaction is fine)
	po, err := s.poRepo.FindByID(req.PurchaseOrderID)
	if err != nil {
		return nil, ErrPurchaseOrderNotFound
	}

	// Check if order is receivable (ordered or partial status)
	if po.Status != purchase.StatusOrdered && po.Status != purchase.StatusPartial {
		return nil, ErrOrderNotReceivable
	}

	// Build a map of PO items for quick lookup
	poItemMap := make(map[uuid.UUID]*purchase.PurchaseOrderItem)
	for i := range po.Items {
		poItemMap[po.Items[i].ID] = &po.Items[i]
	}

	// Validate items and calculate total (validation - outside transaction)
	var totalAmount float64
	for _, itemReq := range req.Items {
		poItem, exists := poItemMap[itemReq.PurchaseOrderItemID]
		if !exists {
			return nil, ErrPurchaseOrderItemNotFound
		}

		remaining := poItem.QuantityOrdered - poItem.QuantityReceived
		if itemReq.QuantityReceived > remaining {
			return nil, ErrQuantityExceedsOrdered
		}

		itemTotal := float64(itemReq.QuantityReceived) * poItem.CostPrice
		totalAmount += itemTotal
	}

	// Generate receipt number before transaction
	receiptNumber := fmt.Sprintf("SR-%s-%s", req.ReceivedDate.Format("20060102"), uuid.New().String()[:8])

	// Prepare receipt ID
	receiptID := uuid.New()

	// Execute all write operations in a transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create tx-scoped repositories
		txReceiptRepo := NewStockReceiptRepository(tx)
		txBatchRepo := inventory.NewBatchRepository(tx)
		txPORepo := purchase.NewPurchaseOrderRepository(tx)
		txVendorRepo := vendor.NewVendorRepository(tx)

		// Create receipt
		receipt := &StockReceipt{
			ID:              receiptID,
			PurchaseOrderID: req.PurchaseOrderID,
			ReceiptNumber:   receiptNumber,
			ReceivedDate:    req.ReceivedDate,
			TotalAmount:     totalAmount,
			Notes:           req.Notes,
			AuditFields: common.AuditFields{
				CreatedBy: user.ID,
				UpdatedBy: user.ID,
			},
		}

		if err := txReceiptRepo.CreateWithoutNumber(receipt); err != nil {
			return err
		}

		// Process each item
		for _, itemReq := range req.Items {
			poItem := poItemMap[itemReq.PurchaseOrderItemID]

			// Create product batch
			batch := &inventory.ProductBatch{
				ID:                uuid.New(),
				ProductID:         poItem.ProductID,
				CostPrice:         poItem.CostPrice,
				SellingPrice:      poItem.SellingPrice,
				QuantityAvailable: itemReq.QuantityReceived,
				PurchasedAt:       req.ReceivedDate,
				ExpiresAt:         poItem.ExpiresAt,
				AuditFields: common.AuditFields{
					CreatedBy: user.ID,
					UpdatedBy: user.ID,
				},
			}

			if err := txBatchRepo.Create(batch); err != nil {
				return err
			}

			// Create receipt item
			receiptItem := &StockReceiptItem{
				ID:                  uuid.New(),
				StockReceiptID:      receiptID,
				PurchaseOrderItemID: itemReq.PurchaseOrderItemID,
				ProductBatchID:      batch.ID,
				QuantityReceived:    itemReq.QuantityReceived,
				AuditFields: common.AuditFields{
					CreatedBy: user.ID,
					UpdatedBy: user.ID,
				},
			}

			if err := tx.Create(receiptItem).Error; err != nil {
				return err
			}

			// Update PO item quantity received
			newQtyReceived := poItem.QuantityReceived + itemReq.QuantityReceived
			if err := txPORepo.UpdateItemQuantityReceived(poItem.ID, newQtyReceived); err != nil {
				return err
			}

			// Update local map for status calculation
			poItem.QuantityReceived = newQtyReceived
		}

		// Determine new PO status based on updated quantities
		allReceived := true
		anyReceived := false
		for _, poItem := range poItemMap {
			if poItem.QuantityReceived > 0 {
				anyReceived = true
			}
			if poItem.QuantityReceived < poItem.QuantityOrdered {
				allReceived = false
			}
		}

		var newStatus string
		if allReceived {
			newStatus = purchase.StatusReceived
		} else if anyReceived {
			newStatus = purchase.StatusPartial
		} else {
			newStatus = po.Status
		}

		if newStatus != po.Status {
			po.Status = newStatus
			po.UpdatedBy = user.ID
			if err := txPORepo.Update(po); err != nil {
				return err
			}
		}

		// Update vendor balance
		vendorRecord, err := txVendorRepo.FindByID(po.VendorID)
		if err != nil {
			return err
		}

		newBalance := vendorRecord.Balance + totalAmount
		if err := txVendorRepo.UpdateBalance(po.VendorID, newBalance); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload and return the receipt
	return s.repo.FindByID(receiptID)
}
