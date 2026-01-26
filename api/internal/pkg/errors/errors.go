package errors

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Common application errors.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrBadRequest        = errors.New("bad request")
	ErrInternalServer    = errors.New("internal server error")
	ErrConflict          = errors.New("conflict")
	ErrUnprocessable     = errors.New("unprocessable entity")
)

// AppError represents an application error with additional context.
type AppError struct {
	Err     error
	Message string
	Code    int
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(code int, err error, message string) *AppError {
	return &AppError{
		Code:    code,
		Err:     err,
		Message: message,
	}
}

// Newf creates a new AppError with a formatted message.
func Newf(code int, err error, format string, args ...any) *AppError {
	return &AppError{
		Code:    code,
		Err:     err,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an error with a message.
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an error with a formatted message.
func Wrapf(err error, format string, args ...any) error {
	return fmt.Errorf(fmt.Sprintf(format, args...)+": %w", err)
}

// IsNotFound checks if the error is a "not found" error.
// This includes both our ErrNotFound and GORM's ErrRecordNotFound.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound)
}

// IsConstraintViolation checks if the error is a database constraint violation.
// This includes unique constraint, foreign key, and other database constraints.
func IsConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// SQLite constraint violation patterns
	return strings.Contains(errStr, "UNIQUE constraint failed") ||
		strings.Contains(errStr, "FOREIGN KEY constraint failed") ||
		strings.Contains(errStr, "NOT NULL constraint failed") ||
		strings.Contains(errStr, "CHECK constraint failed")
}
