package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/product"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
)

// setupRoutes configures all application routes.
func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/healthz", s.handleHealth)

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Auth module
		authHandler := auth.NewHandler(s.db)
		r.Mount("/auth", authHandler.Routes())

		// Product module
		productHandler := product.NewHandler(s.db)
		r.Mount("/products", productHandler.Routes())
	})
}

// handleHealth handles health check requests.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.db.Health(); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	response.Success(w, map[string]string{
		"status": "ok",
	})
}
