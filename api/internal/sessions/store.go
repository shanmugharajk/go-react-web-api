package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	// ErrSessionNotFound is returned when a session is not found
	ErrSessionNotFound = errors.New("session not found")
	
	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = errors.New("session expired")
)

// Session represents a user session stored in the database.
type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for the Session model.
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session has expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Store defines the interface for session storage operations.
type Store interface {
	Create(userID int64, duration time.Duration) (*Session, error)
	Get(sessionID string) (*Session, error)
	Delete(sessionID string) error
	DeleteExpired() error
}

// SQLiteStore implements the Store interface using SQLite/GORM.
type SQLiteStore struct {
	db              *gorm.DB
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
	wg              sync.WaitGroup
}

// NewStore creates a new session store with automatic cleanup.
func NewStore(db *gorm.DB, cleanupInterval time.Duration) *SQLiteStore {
	store := &SQLiteStore{
		db:              db,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}
	
	// Start background cleanup goroutine
	store.startCleanup()
	
	return store
}

// Create creates a new session with a cryptographically secure random ID.
func (s *SQLiteStore) Create(userID int64, duration time.Duration) (*Session, error) {
	// Generate cryptographically secure random session ID (32 bytes)
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// Get retrieves a session by ID.
func (s *SQLiteStore) Get(sessionID string) (*Session, error) {
	var session Session
	
	err := s.db.Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session has expired
	if session.IsExpired() {
		// Delete expired session
		_ = s.Delete(sessionID)
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// Delete removes a session by ID.
func (s *SQLiteStore) Delete(sessionID string) error {
	result := s.db.Where("id = ?", sessionID).Delete(&Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete session: %w", result.Error)
	}
	
	return nil
}

// DeleteExpired removes all expired sessions.
func (s *SQLiteStore) DeleteExpired() error {
	result := s.db.Where("expires_at < ?", time.Now()).Delete(&Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", result.Error)
	}
	
	if result.RowsAffected > 0 {
		fmt.Printf("Cleaned up %d expired sessions\n", result.RowsAffected)
	}
	
	return nil
}

// startCleanup starts a background goroutine to periodically clean up expired sessions.
func (s *SQLiteStore) startCleanup() {
	s.wg.Add(1)
	
	go func() {
		defer s.wg.Done()
		
		ticker := time.NewTicker(s.cleanupInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if err := s.DeleteExpired(); err != nil {
					fmt.Printf("Error cleaning up expired sessions: %v\n", err)
				}
			case <-s.stopCleanup:
				return
			}
		}
	}()
}

// Stop stops the background cleanup goroutine.
func (s *SQLiteStore) Stop() {
	close(s.stopCleanup)
	s.wg.Wait()
}

// generateSessionID generates a cryptographically secure random session ID.
func generateSessionID() (string, error) {
	// Generate 32 random bytes
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	
	// Encode to base64 URL-safe string
	return base64.URLEncoding.EncodeToString(b), nil
}
