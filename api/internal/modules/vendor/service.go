package vendor

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/common"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// VendorService handles business logic for vendors.
type VendorService struct {
	repo *VendorRepository
}

// NewVendorService creates a new VendorService instance.
func NewVendorService(repo *VendorRepository) *VendorService {
	return &VendorService{repo: repo}
}

// GetAll retrieves all active vendors.
func (s *VendorService) GetAll() ([]Vendor, error) {
	return s.repo.FindAll()
}

// GetByID retrieves a vendor by ID.
func (s *VendorService) GetByID(id uuid.UUID) (*Vendor, error) {
	return s.repo.FindByID(id)
}

// Create creates a new vendor.
func (s *VendorService) Create(req CreateVendorRequest, user *auth.User) (*Vendor, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	vendor := &Vendor{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Balance:       0,
		Active:        true,
		AuditFields: common.AuditFields{
			CreatedBy: user.ID,
			UpdatedBy: user.ID,
		},
	}

	if err := s.repo.Create(vendor); err != nil {
		return nil, err
	}

	return vendor, nil
}

// Update updates an existing vendor.
func (s *VendorService) Update(id uuid.UUID, req UpdateVendorRequest, user *auth.User) (*Vendor, error) {
	if err := validator.Struct(req); err != nil {
		return nil, err
	}

	vendor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	vendor.Name = req.Name
	vendor.ContactPerson = req.ContactPerson
	vendor.Phone = req.Phone
	vendor.Email = req.Email
	vendor.Address = req.Address
	vendor.Active = req.Active
	vendor.UpdatedBy = user.ID

	if err := s.repo.Update(vendor); err != nil {
		return nil, err
	}

	return vendor, nil
}

// Delete soft deletes a vendor.
func (s *VendorService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
