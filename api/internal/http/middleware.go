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
