package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
)

// setupMiddleware configures global middleware for the router.
func (s *Server) setupMiddleware() {
	// Request ID
	s.router.Use(middleware.RequestID)

	// Real IP
	s.router.Use(middleware.RealIP)

	// Logger
	s.router.Use(requestLogger)

	// Recoverer
	s.router.Use(middleware.Recoverer)

	// Security headers (OWASP recommended)
	s.router.Use(securityHeaders(s.config.Auth.IsDevelopment))

	// Timeout
	s.router.Use(middleware.Timeout(30 * time.Second))

	// Content-Type JSON
	s.router.Use(middleware.SetHeader("Content-Type", "application/json"))
}

// requestLogger is a custom logging middleware using slog.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration", time.Since(start).String(),
				"remote", r.RemoteAddr,
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

// securityHeaders adds OWASP-recommended security headers to all responses.
func securityHeaders(isDevelopment bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME-sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")
			
			// Prevent clickjacking
			w.Header().Set("X-Frame-Options", "DENY")
			
			// XSS protection (legacy, but still useful)
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			
			// Content Security Policy - only allow resources from same origin
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			
			// Referrer policy - don't leak referrer information
			w.Header().Set("Referrer-Policy", "no-referrer")
			
			// Permissions policy - disable dangerous features
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			
			// HSTS - enforce HTTPS (only in production)
			if !isDevelopment {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
