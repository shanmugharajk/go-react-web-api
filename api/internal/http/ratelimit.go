package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/iputil"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
)

// rateLimitEntry tracks request attempts for a single IP address.
type rateLimitEntry struct {
	attempts  int
	firstSeen time.Time
}

// rateLimiter implements a sliding window rate limiter.
type rateLimiter struct {
	mu          sync.RWMutex
	entries     map[string]*rateLimitEntry
	maxAttempts int
	window      time.Duration
}

// newRateLimiter creates a new rate limiter.
func newRateLimiter(maxAttempts int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		entries:     make(map[string]*rateLimitEntry),
		maxAttempts: maxAttempts,
		window:      window,
	}
	
	// Start background cleanup
	go rl.cleanup()
	
	return rl
}

// allow checks if a request from the given IP should be allowed.
// Returns true if allowed, false if rate limit exceeded.
func (rl *rateLimiter) allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	entry, exists := rl.entries[ip]
	
	if !exists {
		// First request from this IP
		rl.entries[ip] = &rateLimitEntry{
			attempts:  1,
			firstSeen: now,
		}
		return true, 0
	}
	
	// Check if the window has expired
	if now.Sub(entry.firstSeen) > rl.window {
		// Reset the window
		entry.attempts = 1
		entry.firstSeen = now
		return true, 0
	}
	
	// Increment attempts
	entry.attempts++
	
	// Check if limit exceeded (> because first request creates with attempts=1)
	// With maxAttempts=5: allows requests 1-5, blocks 6+
	if entry.attempts > rl.maxAttempts {
		retryAfter := rl.window - now.Sub(entry.firstSeen)
		return false, retryAfter
	}
	
	return true, 0
}

// cleanup periodically removes expired entries to prevent memory leaks.
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		
		for ip, entry := range rl.entries {
			if now.Sub(entry.firstSeen) > rl.window {
				delete(rl.entries, ip)
			}
		}
		
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a middleware that limits requests per IP.
// trustProxy determines whether to trust X-Forwarded-For and X-Real-IP headers.
func RateLimitMiddleware(maxAttempts int, window time.Duration, trustProxy bool) func(http.Handler) http.Handler {
	limiter := newRateLimiter(maxAttempts, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP using secure extraction
			ip := iputil.ExtractClientIP(r, trustProxy)

			// Check if request is allowed
			allowed, retryAfter := limiter.allow(ip)
			if !allowed {
				// Set Retry-After header
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(retryAfter.Seconds())))

				response.Error(w, http.StatusTooManyRequests,
					"too many requests, please try again later")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
