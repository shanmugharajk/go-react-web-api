package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate is the package-level validator instance.
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// ValidationError represents a validation error with field details.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error implements the error interface.
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}

	var msgs []string
	for _, e := range v {
		msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(msgs, "; ")
}

// Struct validates a struct and returns ValidationErrors if invalid.
func Struct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return toValidationErrors(validationErrs)
	}

	return err
}

// toValidationErrors converts validator.ValidationErrors to ValidationErrors.
func toValidationErrors(errs validator.ValidationErrors) ValidationErrors {
	result := make(ValidationErrors, 0, len(errs))

	for _, err := range errs {
		result = append(result, ValidationError{
			Field:   toCamelCase(err.Field()),
			Message: formatMessage(err),
		})
	}

	return result
}

// formatMessage creates a human-readable message for a validation error.
func formatMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", err.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", err.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", err.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", err.Param())
	case "lt":
		return fmt.Sprintf("must be less than %s", err.Param())
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid URL"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	default:
		return fmt.Sprintf("failed validation: %s", err.Tag())
	}
}

// toCamelCase converts a PascalCase field name to camelCase.
func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	// Lowercase the first character
	runes := []rune(s)
	runes[0] = rune(strings.ToLower(string(runes[0]))[0])
	return string(runes)
}

// IsValidationError checks if an error is a ValidationErrors type.
func IsValidationError(err error) bool {
	var validationErrs ValidationErrors
	return errors.As(err, &validationErrs)
}

// GetValidationErrors extracts ValidationErrors from an error.
func GetValidationErrors(err error) (ValidationErrors, bool) {
	var validationErrs ValidationErrors
	if errors.As(err, &validationErrs) {
		return validationErrs, true
	}
	return nil, false
}
