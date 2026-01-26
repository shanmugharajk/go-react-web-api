package payment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VendorPaymentRepository handles database operations for vendor payments.
type VendorPaymentRepository struct {
	db *gorm.DB
}

// NewVendorPaymentRepository creates a new VendorPaymentRepository instance.
func NewVendorPaymentRepository(db *gorm.DB) *VendorPaymentRepository {
	return &VendorPaymentRepository{db: db}
}

// FindAll retrieves all vendor payments.
func (r *VendorPaymentRepository) FindAll() ([]VendorPayment, error) {
	var payments []VendorPayment
	if err := r.db.Order("payment_date DESC").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// FindByID retrieves a vendor payment by ID.
func (r *VendorPaymentRepository) FindByID(id uuid.UUID) (*VendorPayment, error) {
	var payment VendorPayment
	if err := r.db.First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByVendorID retrieves all payments for a vendor.
func (r *VendorPaymentRepository) FindByVendorID(vendorID uuid.UUID) ([]VendorPayment, error) {
	var payments []VendorPayment
	if err := r.db.Where("vendor_id = ?", vendorID).Order("payment_date DESC").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// Create inserts a new vendor payment.
func (r *VendorPaymentRepository) Create(payment *VendorPayment) error {
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}

	// Generate payment number
	payment.PaymentNumber = generatePaymentNumber(payment.PaymentDate)

	return r.db.Create(payment).Error
}

// GetDB returns the underlying database for transaction support.
func (r *VendorPaymentRepository) GetDB() *gorm.DB {
	return r.db
}

// generatePaymentNumber generates a unique payment number.
func generatePaymentNumber(paymentDate time.Time) string {
	return fmt.Sprintf("VP-%s-%s", paymentDate.Format("20060102"), uuid.New().String()[:8])
}
