package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shanmugharajk/go-react-web-api/api/internal/config"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
)

// Server represents the HTTP server.
type Server struct {
	router *chi.Mux
	server *http.Server
	db     *db.DB
	cfg    *config.Config
}

// New creates a new HTTP server instance.
func New(cfg *config.Config, database *db.DB) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     database,
		cfg:    cfg,
	}

	// Setup middleware
	s.setupMiddleware()

	// Setup routes
	s.setupRoutes()

	s.server = &http.Server{
		Addr:         s.cfg.ServerAddr(),
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	logger.Info("Starting HTTP server", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
