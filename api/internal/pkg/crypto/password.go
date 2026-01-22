package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters (OWASP recommended)
	argon2Time       = 1         // Number of iterations
	argon2Memory     = 64 * 1024 // Memory in KiB (64 MB)
	argon2Threads    = 4         // Number of threads
	argon2KeyLength  = 32        // Length of generated key in bytes
	argon2SaltLength = 16        // Length of salt in bytes
)

var (
	// ErrInvalidPasswordHash is returned when the password hash format is invalid
	ErrInvalidPasswordHash = errors.New("invalid password hash format")

	// ErrPasswordTooShort is returned when the password is too short
	ErrPasswordTooShort = errors.New("password must be at least 12 characters")

	// ErrPasswordTooLong is returned when the password is too long (DoS prevention)
	ErrPasswordTooLong = errors.New("password must not exceed 128 characters")
)

// HashPassword hashes a password using Argon2id with a random salt.
// Returns the hash in the format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
func HashPassword(password string) (string, error) {
	// Validate password strength
	if err := ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	// Generate random salt
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLength,
	)

	// Encode to base64
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	// Return in PHC string format
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Memory,
		argon2Time,
		argon2Threads,
		saltB64,
		hashB64,
	), nil
}

// VerifyPassword verifies a password against a hash using constant-time comparison.
// Returns nil if the password matches, otherwise returns an error.
func VerifyPassword(hashedPassword, password string) error {
	// Parse the hash
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return ErrInvalidPasswordHash
	}

	if parts[1] != "argon2id" {
		return ErrInvalidPasswordHash
	}

	// Parse parameters
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return ErrInvalidPasswordHash
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return ErrInvalidPasswordHash
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}

	// Decode hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	// Generate hash from provided password using same parameters
	passwordHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(hash)),
	)

	// Constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, passwordHash) == 1 {
		return nil
	}

	return errors.New("invalid password")
}

// ValidatePasswordStrength validates password meets minimum security requirements.
func ValidatePasswordStrength(password string) error {
	if len(password) < 12 {
		return ErrPasswordTooShort
	}

	if len(password) > 128 {
		return ErrPasswordTooLong
	}

	// Additional strength requirements can be added here:
	// - Check for common passwords
	// - Require mix of characters
	// - Check against breach databases

	return nil
}
