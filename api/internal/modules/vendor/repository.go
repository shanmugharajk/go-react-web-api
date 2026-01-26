package vendor

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/errors"
	"gorm.io/gorm"
)

// VendorRepository handles database operations for vendors.
type VendorRepository struct {
	db *gorm.DB
}

// NewVendorRepository creates a new VendorRepository instance.
func NewVendorRepository(db *gorm.DB) *VendorRepository {
	return &VendorRepository{db: db}
}

// FindAll retrieves all active vendors.
func (r *VendorRepository) FindAll() ([]Vendor, error) {
	var vendors []Vendor
	if err := r.db.Where("active = ?", true).Find(&vendors).Error; err != nil {
		return nil, err
	}
	return vendors, nil
}

// FindByID retrieves a vendor by ID.
func (r *VendorRepository) FindByID(id uuid.UUID) (*Vendor, error) {
	var vendor Vendor
	if err := r.db.First(&vendor, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

// Create inserts a new vendor into the database.
func (r *VendorRepository) Create(vendor *Vendor) error {
	if vendor.ID == uuid.Nil {
		vendor.ID = uuid.New()
	}
	return r.db.Create(vendor).Error
}

// Update updates an existing vendor in the database.
func (r *VendorRepository) Update(vendor *Vendor) error {
	return r.db.Save(vendor).Error
}

// Delete soft deletes a vendor by setting active to false.
func (r *VendorRepository) Delete(id uuid.UUID) error {
	result := r.db.Model(&Vendor{}).Where("id = ?", id).Update("active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// UpdateBalance updates the vendor's balance.
func (r *VendorRepository) UpdateBalance(id uuid.UUID, amount float64) error {
	result := r.db.Model(&Vendor{}).Where("id = ?", id).Update("balance", amount)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}
	return nil
}
