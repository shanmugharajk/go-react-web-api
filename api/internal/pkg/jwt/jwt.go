package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when token validation fails
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken is returned when token has expired
	ErrExpiredToken = errors.New("token expired")

	// ErrInvalidClaims is returned when claims are invalid
	ErrInvalidClaims = errors.New("invalid token claims")
)

// Claims represents the JWT claims for authentication.
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenService handles JWT token generation and validation.
type TokenService struct {
	secret []byte
}

// NewTokenService creates a new JWT token service.
func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
	}
}

// Generate creates a new JWT token for the given user ID with specified duration.
func (s *TokenService) Generate(userID int64, duration time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-react-web-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Validate validates a JWT token and returns the user ID.
func (s *TokenService) Validate(tokenString string) (int64, error) {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrExpiredToken
		}
		return 0, ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, ErrInvalidClaims
	}

	return claims.UserID, nil
}

// ValidateAndRefresh validates a token and returns a new token if it's close to expiry.
// Returns (userID, newToken, error). newToken is empty string if refresh not needed.
func (s *TokenService) ValidateAndRefresh(tokenString string, refreshThreshold, newDuration time.Duration) (int64, string, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, "", ErrExpiredToken
		}
		return 0, "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", ErrInvalidClaims
	}

	// Check if token needs refresh (within threshold of expiry)
	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)
	if timeUntilExpiry <= refreshThreshold {
		// Generate new token
		newToken, err := s.Generate(claims.UserID, newDuration)
		if err != nil {
			return claims.UserID, "", fmt.Errorf("failed to refresh token: %w", err)
		}
		return claims.UserID, newToken, nil
	}

	return claims.UserID, "", nil
}
