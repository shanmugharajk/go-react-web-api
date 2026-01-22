package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port int
	Host string
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	DSN string // Data Source Name
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	SessionSecret   string // Secret for signing session IDs (HMAC)
	CSRFSecret      string // Secret for CSRF token generation
	JWTSecret       string // Secret for JWT token signing
	IsDevelopment   bool   // Toggle for development mode (affects cookie security)
	SessionDuration int    // Session TTL in seconds
	JWTDuration     int    // JWT TTL in seconds
	TrustProxy      bool   // Whether to trust X-Forwarded-For and X-Real-IP headers
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("PORT", 8080),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DATABASE_DSN", "file:./data/pos.db"),
		},
		Auth: AuthConfig{
			SessionSecret:   getEnv("AUTH_SESSION_SECRET", "dev-secret-change-in-production-min-32-chars"),
			CSRFSecret:      getEnv("AUTH_CSRF_SECRET", "dev-csrf-secret-change-in-production-32-chars"),
			JWTSecret:       getEnv("AUTH_JWT_SECRET", "dev-jwt-secret-change-in-production-min-32-chars"),
			IsDevelopment:   getEnvAsBool("AUTH_IS_DEVELOPMENT", true),
			SessionDuration: getEnvAsInt("AUTH_SESSION_DURATION", 86400),  // 24 hours
			JWTDuration:     getEnvAsInt("AUTH_JWT_DURATION", 86400),      // 24 hours
			TrustProxy:      getEnvAsBool("AUTH_TRUST_PROXY", false),      // Only trust proxy headers if explicitly enabled
		},
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as int or returns a default value.
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvAsBool retrieves an environment variable as bool or returns a default value.
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// ServerAddr returns the full server address.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	// Validate secrets in production
	if !c.Auth.IsDevelopment {
		if err := validateSecret(c.Auth.SessionSecret, "AUTH_SESSION_SECRET"); err != nil {
			return err
		}
		if err := validateSecret(c.Auth.CSRFSecret, "AUTH_CSRF_SECRET"); err != nil {
			return err
		}
		if err := validateSecret(c.Auth.JWTSecret, "AUTH_JWT_SECRET"); err != nil {
			return err
		}
	}

	// Validate session duration
	if c.Auth.SessionDuration < 60 {
		return fmt.Errorf("AUTH_SESSION_DURATION must be at least 60 seconds")
	}
	if c.Auth.SessionDuration > 2592000 { // 30 days
		return fmt.Errorf("AUTH_SESSION_DURATION cannot exceed 30 days (2592000 seconds)")
	}

	// Validate JWT duration
	if c.Auth.JWTDuration < 60 {
		return fmt.Errorf("AUTH_JWT_DURATION must be at least 60 seconds")
	}
	if c.Auth.JWTDuration > 2592000 { // 30 days
		return fmt.Errorf("AUTH_JWT_DURATION cannot exceed 30 days (2592000 seconds)")
	}

	return nil
}

// validateSecret validates that a secret meets security requirements.
func validateSecret(secret, name string) error {
	// Check minimum length
	if len(secret) < 32 {
		return fmt.Errorf("%s must be at least 32 characters long (got %d)", name, len(secret))
	}

	// Check maximum length (reasonable upper bound)
	if len(secret) > 256 {
		return fmt.Errorf("%s cannot exceed 256 characters (got %d)", name, len(secret))
	}

	// Check if using default development secret in production
	if secret == "dev-secret-change-in-production-min-32-chars" ||
		secret == "dev-csrf-secret-change-in-production-32-chars" ||
		secret == "dev-jwt-secret-change-in-production-min-32-chars" {
		return fmt.Errorf("%s is using default development value in production - MUST be changed", name)
	}

	// Check for weak secrets (all same character)
	if isRepeatingChar(secret) {
		return fmt.Errorf("%s appears to be weak (repeating characters)", name)
	}

	return nil
}

// isRepeatingChar checks if a string is all the same character.
func isRepeatingChar(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}
