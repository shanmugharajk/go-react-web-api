package auth

import "time"

// User represents a user in the system.
type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"column:password;not null" json:"-"` // Never expose password hash
	Name         string    `gorm:"not null" json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// LoginRequest represents a login request payload.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents a registration request payload.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginResponse represents the response after successful login (cookie-based).
type LoginResponse struct {
	User      *User  `json:"user"`
	CSRFToken string `json:"csrfToken"`
}

// TokenLoginResponse represents the response after successful token-based login.
// Follows RFC 6750 - The OAuth 2.0 Authorization Framework: Bearer Token Usage
type TokenLoginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"` // Always "Bearer"
	ExpiresIn   int    `json:"expiresIn"` // Token lifetime in seconds
	User        *User  `json:"user"`
}
