package auth

import (
	"context"
	"errors"
)

// contextKey is a type-safe key for context values.
type contextKey string

const (
	// UserIDKey is the context key for storing user ID.
	// This is exported so the http middleware can use it.
	UserIDKey contextKey = "userID"
)

// GetUserIDFromContext extracts the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
