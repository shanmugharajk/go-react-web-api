package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shanmugharajk/go-react-web-api/api/internal/config"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/jwt"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
	"github.com/shanmugharajk/go-react-web-api/api/internal/sessions"
)

// Server represents the HTTP server.
type Server struct {
	router       *chi.Mux
	server       *http.Server
	db           *db.DB
	config       *config.Config
	sessionStore *sessions.SQLiteStore
	jwtService   *jwt.TokenService
	authService  *auth.Service
}

// New creates a new HTTP server instance.
func New(cfg *config.Config, database *db.DB) *Server {
	// Initialize session store with 1-hour cleanup interval
	sessionStore := sessions.NewStore(
		database.DB,
		time.Hour,
	)

	// Initialize JWT service
	jwtService := jwt.NewTokenService(cfg.Auth.JWTSecret)

	// Initialize auth service for user lookup
	authRepo := auth.NewRepository(database.DB)
	authService := auth.NewService(authRepo)

	s := &Server{
		router:       chi.NewRouter(),
		db:           database,
		config:       cfg,
		sessionStore: sessionStore,
		jwtService:   jwtService,
		authService:  authService,
	}

	// Setup middleware
	s.setupMiddleware()

	// Setup routes (CSRF will be applied selectively within routes)
	s.setupRoutes()

	s.server = &http.Server{
		Addr:         cfg.ServerAddr(),
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

	// Stop session cleanup goroutine
	s.sessionStore.Stop()

	return s.server.Shutdown(ctx)
}
