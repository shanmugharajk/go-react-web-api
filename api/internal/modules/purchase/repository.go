package purchase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/errors"
	"gorm.io/gorm"
)

// PurchaseOrderRepository handles database operations for purchase orders.
type PurchaseOrderRepository struct {
	db *gorm.DB
}

// NewPurchaseOrderRepository creates a new PurchaseOrderRepository instance.
func NewPurchaseOrderRepository(db *gorm.DB) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{db: db}
}

// FindAll retrieves all purchase orders.
func (r *PurchaseOrderRepository) FindAll() ([]PurchaseOrder, error) {
	var orders []PurchaseOrder
	if err := r.db.Preload("Items").Order("order_date DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// FindByID retrieves a purchase order by ID with items.
func (r *PurchaseOrderRepository) FindByID(id uuid.UUID) (*PurchaseOrder, error) {
	var order PurchaseOrder
	if err := r.db.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByVendorID retrieves all purchase orders for a vendor.
func (r *PurchaseOrderRepository) FindByVendorID(vendorID uuid.UUID) ([]PurchaseOrder, error) {
	var orders []PurchaseOrder
	if err := r.db.Preload("Items").Where("vendor_id = ?", vendorID).Order("order_date DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// FindUnpaidByVendorID retrieves unpaid/partial purchase orders for a vendor, ordered by date ASC (for FIFO).
func (r *PurchaseOrderRepository) FindUnpaidByVendorID(vendorID uuid.UUID) ([]PurchaseOrder, error) {
	var orders []PurchaseOrder
	if err := r.db.Where("vendor_id = ? AND payment_status != ? AND status != ?", vendorID, PaymentStatusPaid, StatusCancelled).
		Order("order_date ASC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// Create inserts a new purchase order with items.
func (r *PurchaseOrderRepository) Create(order *PurchaseOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if order.ID == uuid.Nil {
			order.ID = uuid.New()
		}

		// Generate order number
		order.OrderNumber = generateOrderNumber(order.OrderDate)

		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		return nil
	})
}

// Update updates an existing purchase order.
func (r *PurchaseOrderRepository) Update(order *PurchaseOrder) error {
	return r.db.Save(order).Error
}

// UpdateWithItems updates a purchase order and replaces all items.
func (r *PurchaseOrderRepository) UpdateWithItems(order *PurchaseOrder, items []PurchaseOrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing items
		if err := tx.Where("purchase_order_id = ?", order.ID).Delete(&PurchaseOrderItem{}).Error; err != nil {
			return err
		}

		// Update order
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		// Create new items
		for i := range items {
			items[i].PurchaseOrderID = order.ID
			if items[i].ID == uuid.Nil {
				items[i].ID = uuid.New()
			}
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete cancels a purchase order (only if draft).
func (r *PurchaseOrderRepository) Delete(id uuid.UUID) error {
	result := r.db.Model(&PurchaseOrder{}).
		Where("id = ? AND status = ?", id, StatusDraft).
		Update("status", StatusCancelled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// UpdatePayment updates payment fields on a purchase order.
func (r *PurchaseOrderRepository) UpdatePayment(id uuid.UUID, paidAmount float64, paymentStatus string, lastPaymentAt time.Time) error {
	result := r.db.Model(&PurchaseOrder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"paid_amount":     paidAmount,
			"payment_status":  paymentStatus,
			"last_payment_at": lastPaymentAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// FindItemByID retrieves a purchase order item by ID.
func (r *PurchaseOrderRepository) FindItemByID(id uuid.UUID) (*PurchaseOrderItem, error) {
	var item PurchaseOrderItem
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItemQuantityReceived updates the quantity received for an item.
func (r *PurchaseOrderRepository) UpdateItemQuantityReceived(id uuid.UUID, quantityReceived int) error {
	result := r.db.Model(&PurchaseOrderItem{}).
		Where("id = ?", id).
		Update("quantity_received", quantityReceived)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// generateOrderNumber generates a unique order number.
func generateOrderNumber(orderDate time.Time) string {
	return fmt.Sprintf("PO-%s-%s", orderDate.Format("20060102"), uuid.New().String()[:8])
}
