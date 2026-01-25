package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the GORM database connection.
type DB struct {
	*gorm.DB
}

// New initializes and returns a new database connection.
func New(dsn string) (*DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure SQLite pragmas for optimal transaction behavior
	if err := gormDB.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	if err := gormDB.Exec("PRAGMA journal_mode = WAL").Error; err != nil {
		return nil, fmt.Errorf("failed to set journal_mode: %w", err)
	}
	if err := gormDB.Exec("PRAGMA synchronous = NORMAL").Error; err != nil {
		return nil, fmt.Errorf("failed to set synchronous mode: %w", err)
	}
	if err := gormDB.Exec("PRAGMA busy_timeout = 5000").Error; err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	return &DB{DB: gormDB}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Health performs a simple health check on the database.
func (db *DB) Health() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Ping()
}
