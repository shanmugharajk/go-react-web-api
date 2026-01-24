package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// contextKey is a type-safe key for context values.
type contextKey string

const (
	// UserIDKey is the context key for storing user ID.
	// This is exported so the http middleware can use it.
	UserIDKey contextKey = "userID"
	// UserKey is the context key for storing the authenticated user.
	// This is exported so handlers can access the full user object without DB queries.
	UserKey contextKey = "user"
)

// GetUserIDFromContext extracts the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetUserFromContext extracts the authenticated user from the request context.
func GetUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(UserKey).(*User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
