package purchase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

var (
	ErrCannotUpdateNonDraft = errors.New("can only update purchase orders in draft status")
	ErrCannotDeleteNonDraft = errors.New("can only delete purchase orders in draft status")
)

// PurchaseOrderService handles business logic for purchase orders.
type PurchaseOrderService struct {
	repo *PurchaseOrderRepository
}

// NewPurchaseOrderService creates a new PurchaseOrderService instance.
func NewPurchaseOrderService(repo *PurchaseOrderRepository) *PurchaseOrderService {
	return &PurchaseOrderService{repo: repo}
}

// GetAll retrieves all purchase orders.
func (s *PurchaseOrderService) GetAll() ([]PurchaseOrder, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a purchase order by ID.
func (s *PurchaseOrderService) GetByID(id uuid.UUID) (*PurchaseOrder, error) {
	return s.repo.FindByID(id)
}

// GetByVendorID retrieves all purchase orders for a vendor.
func (s *PurchaseOrderService) GetByVendorID(vendorID uuid.UUID) ([]PurchaseOrder, error) {
	return s.repo.FindByVendorID(vendorID)
}

// Create creates a new purchase order with items.
func (s *PurchaseOrderService) Create(req CreatePurchaseOrderRequest, user *auth.User) (*PurchaseOrder, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	// Calculate total amount
	var totalAmount float64
	items := make([]PurchaseOrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		itemTotal := float64(itemReq.QuantityOrdered) * itemReq.CostPrice
		totalAmount += itemTotal

		items[i] = PurchaseOrderItem{
			ID:              uuid.New(),
			ProductID:       itemReq.ProductID,
			QuantityOrdered: itemReq.QuantityOrdered,
			CostPrice:       itemReq.CostPrice,
			SellingPrice:    itemReq.SellingPrice,
			ExpiresAt:       itemReq.ExpiresAt,
			AuditFields: common.AuditFields{
				CreatedBy: user.ID,
				UpdatedBy: user.ID,
			},
		}
	}

	order := &PurchaseOrder{
		VendorID:      req.VendorID,
		OrderDate:     req.OrderDate,
		Status:        StatusDraft,
		TotalAmount:   totalAmount,
		PaidAmount:    0,
		PaymentStatus: PaymentStatusUnpaid,
		Notes:         req.Notes,
		Items:         items,
		AuditFields: common.AuditFields{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// Reload to get the complete order with items
	return s.repo.FindByID(order.ID)
}

// Update updates an existing purchase order (only if draft).
func (s *PurchaseOrderService) Update(id uuid.UUID, req UpdatePurchaseOrderRequest, user *auth.User) (*PurchaseOrder, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Only allow updates to draft orders
	if order.Status != StatusDraft {
		return nil, ErrCannotUpdateNonDraft
	}

	// Calculate total amount
	var totalAmount float64
	items := make([]PurchaseOrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		itemTotal := float64(itemReq.QuantityOrdered) * itemReq.CostPrice
		totalAmount += itemTotal

		items[i] = PurchaseOrderItem{
			ID:              itemReq.ID,
			PurchaseOrderID: id,
			ProductID:       itemReq.ProductID,
			QuantityOrdered: itemReq.QuantityOrdered,
			CostPrice:       itemReq.CostPrice,
			SellingPrice:    itemReq.SellingPrice,
			ExpiresAt:       itemReq.ExpiresAt,
			AuditFields: common.AuditFields{
				CreatedBy: user.ID,
				UpdatedBy: user.ID,
			},
		}
	}

	order.VendorID = req.VendorID
	order.OrderDate = req.OrderDate
	order.Status = req.Status
	order.TotalAmount = totalAmount
	order.Notes = req.Notes
	order.UpdatedBy = user.ID

	if err := s.repo.UpdateWithItems(order, items); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

// Delete cancels a purchase order (only if draft).
func (s *PurchaseOrderService) Delete(id uuid.UUID) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if order.Status != StatusDraft {
		return ErrCannotDeleteNonDraft
	}

	return s.repo.Delete(id)
}
