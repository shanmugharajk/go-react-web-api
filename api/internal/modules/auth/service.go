package auth

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/crypto"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service handles business logic for auth.
type Service struct {
	repo *Repository
}

// NewService creates a new auth service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Login authenticates a user with email and password.
// Uses constant-time comparison to prevent timing attacks.
func (s *Service) Login(req LoginRequest, ip, userAgent string) (*User, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		// Hash dummy password to keep timing constant (prevent enumeration)
		_, _ = crypto.HashPassword("dummy-password-to-prevent-timing-attack")
		
		// Audit log failed attempt
		logger.Warn("Login failed - user not found", 
			"email", req.Email, 
			"ip", ip,
			"user_agent", userAgent)
		
		// Return generic error to prevent account enumeration
		return nil, errors.New("invalid credentials")
	}

	// Verify password using constant-time comparison
	if err := crypto.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		// Audit log failed attempt
		logger.Warn("Login failed - invalid password", 
			"email", req.Email, 
			"ip", ip,
			"user_agent", userAgent)
		
		// Return generic error
		return nil, errors.New("invalid credentials")
	}

	// Audit log successful login
	logger.Info("Login successful", 
		"email", req.Email, 
		"user_id", user.ID,
		"ip", ip,
		"user_agent", userAgent)

	return user, nil
}

// Register creates a new user account with password hashing.
func (s *Service) Register(req RegisterRequest, ip, userAgent string) (*User, error) {
	// Validate password strength (12+ characters)
	if err := crypto.ValidatePasswordStrength(req.Password); err != nil {
		logger.Warn("Registration failed - weak password", 
			"email", req.Email, 
			"ip", ip,
			"error", err.Error())
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if user already exists
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		// User already exists - return generic error to prevent enumeration
		logger.Warn("Registration failed - email already exists", 
			"email", req.Email, 
			"ip", ip)
		return nil, errors.New("registration failed")
	}
	
	// Only proceed if error is "not found"
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password using Argon2id
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		logger.Error("Registration failed - password hashing error", 
			"email", req.Email, 
			"error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         req.Name,
	}

	if err := s.repo.Create(user); err != nil {
		logger.Error("Registration failed - database error", 
			"email", req.Email, 
			"error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Audit log successful registration
	logger.Info("Registration successful", 
		"email", req.Email, 
		"user_id", user.ID,
		"ip", ip,
		"user_agent", userAgent)

	return user, nil
}

// GetUserByID retrieves a user by ID.
func (s *Service) GetUserByID(id uuid.UUID) (*User, error) {
	return s.repo.FindByID(id)
}
