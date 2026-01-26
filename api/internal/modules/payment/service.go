package payment

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/purchase"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/vendor"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
	"gorm.io/gorm"
)

var (
	ErrVendorNotFound         = errors.New("vendor not found")
	ErrPaymentExceedsBalance  = errors.New("payment amount exceeds vendor balance")
	ErrInsufficientBalance    = errors.New("vendor has no outstanding balance")
)

// VendorPaymentService handles business logic for vendor payments.
type VendorPaymentService struct {
	repo       *VendorPaymentRepository
	vendorRepo *vendor.VendorRepository
	poRepo     *purchase.PurchaseOrderRepository
	db         *gorm.DB
}

// NewVendorPaymentService creates a new VendorPaymentService instance.
func NewVendorPaymentService(
	repo *VendorPaymentRepository,
	vendorRepo *vendor.VendorRepository,
	poRepo *purchase.PurchaseOrderRepository,
	db *gorm.DB,
) *VendorPaymentService {
	return &VendorPaymentService{
		repo:       repo,
		vendorRepo: vendorRepo,
		poRepo:     poRepo,
		db:         db,
	}
}

// GetAll retrieves all vendor payments.
func (s *VendorPaymentService) GetAll() ([]VendorPayment, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a vendor payment by ID.
func (s *VendorPaymentService) GetByID(id uuid.UUID) (*VendorPayment, error) {
	return s.repo.FindByID(id)
}

// GetByVendorID retrieves all payments for a vendor.
func (s *VendorPaymentService) GetByVendorID(vendorID uuid.UUID) ([]VendorPayment, error) {
	return s.repo.FindByVendorID(vendorID)
}

// Create creates a new vendor payment and allocates it to purchase orders using FIFO.
func (s *VendorPaymentService) Create(req CreateVendorPaymentRequest, user *auth.User) (*VendorPayment, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	// Get vendor
	vendorRecord, err := s.vendorRepo.FindByID(req.VendorID)
	if err != nil {
		return nil, ErrVendorNotFound
	}

	// Check if vendor has outstanding balance
	if vendorRecord.Balance <= 0 {
		return nil, ErrInsufficientBalance
	}

	// Check if payment exceeds balance
	if req.Amount > vendorRecord.Balance {
		return nil, ErrPaymentExceedsBalance
	}

	// Prepare payment ID before transaction
	paymentID := uuid.New()

	// Execute all operations in a transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create tx-scoped repositories
		txPaymentRepo := NewVendorPaymentRepository(tx)
		txVendorRepo := vendor.NewVendorRepository(tx)
		txPORepo := purchase.NewPurchaseOrderRepository(tx)

		// Create payment record
		payment := &VendorPayment{
			ID:            paymentID,
			VendorID:      req.VendorID,
			Amount:        req.Amount,
			PaymentDate:   req.PaymentDate,
			PaymentMethod: req.PaymentMethod,
			Reference:     req.Reference,
			Notes:         req.Notes,
			AuditFields: common.AuditFields{
				CreatedBy: user.ID,
				UpdatedBy: user.ID,
			},
		}

		if err := txPaymentRepo.Create(payment); err != nil {
			return err
		}

		// Deduct from vendor balance
		newBalance := vendorRecord.Balance - req.Amount
		if err := txVendorRepo.UpdateBalance(req.VendorID, newBalance); err != nil {
			return err
		}

		// Get unpaid/partial POs for FIFO allocation
		unpaidPOs, err := txPORepo.FindUnpaidByVendorID(req.VendorID)
		if err != nil {
			return err
		}

		// Allocate payment using FIFO
		remaining := req.Amount
		now := time.Now()

		for _, po := range unpaidPOs {
			if remaining <= 0 {
				break
			}

			outstanding := po.TotalAmount - po.PaidAmount
			if outstanding <= 0 {
				continue
			}

			// Calculate allocation
			allocate := remaining
			if allocate > outstanding {
				allocate = outstanding
			}

			// Update PO payment fields
			newPaidAmount := po.PaidAmount + allocate
			var newPaymentStatus string
			if newPaidAmount >= po.TotalAmount {
				newPaymentStatus = purchase.PaymentStatusPaid
			} else {
				newPaymentStatus = purchase.PaymentStatusPartial
			}

			if err := txPORepo.UpdatePayment(po.ID, newPaidAmount, newPaymentStatus, now); err != nil {
				return err
			}

			remaining -= allocate
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload and return the payment
	return s.repo.FindByID(paymentID)
}
