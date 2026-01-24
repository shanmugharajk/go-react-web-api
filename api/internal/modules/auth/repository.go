package auth

import (
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
)

// Repository handles data access for auth.
type Repository struct {
	db *db.DB
}

// NewRepository creates a new auth repository.
func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

// FindByEmail finds a user by email (case-insensitive).
func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	// Case-insensitive email lookup to prevent enumeration
	if err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user.
func (r *Repository) Create(user *User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return r.db.Create(user).Error
}

// FindByID finds a user by ID.
func (r *Repository) FindByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
