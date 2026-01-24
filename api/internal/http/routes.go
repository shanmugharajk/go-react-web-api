package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/customer"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/product"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
)

// setupRoutes configures all application routes with proper CSRF protection.
func (s *Server) setupRoutes() {
	// Health check (public, no CSRF, no auth)
	s.router.Get("/healthz", s.handleHealth)

	// API routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Create CSRF middleware with custom error handler for JSON responses
		csrfMiddleware := csrf.Protect(
			[]byte(s.config.Auth.CSRFSecret),
			csrf.Secure(!s.config.IsDevelopment),
			csrf.Path("/"),
			csrf.HttpOnly(true),
			csrf.SameSite(csrf.SameSiteLaxMode), // Lax allows top-level navigation
			csrf.RequestHeader("X-CSRF-Token"),
			csrf.FieldName("csrf_token"),
			csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response.Error(w, http.StatusForbidden, "invalid CSRF token")
			})),
		)

		// Wrap CSRF middleware to exempt Bearer token requests
		csrfProtection := csrfProtectUnlessBearerToken(csrfMiddleware)

		// Create auth handler with all dependencies
		sessionTTL := time.Duration(s.config.Auth.SessionDuration) * time.Second
		jwtTTL := time.Duration(s.config.Auth.JWTDuration) * time.Second
		authHandler := auth.NewHandler(
			s.db,
			s.sessionStore,
			s.jwtService,
			s.config.IsDevelopment,
			sessionTTL,
			jwtTTL,
			s.config.Auth.TrustProxy,
		)

		// =================================================================
		// TOKEN-BASED AUTH ROUTES (for API clients: Postman, mobile, etc.)
		// No CSRF required - Bearer tokens are not vulnerable to CSRF
		// =================================================================
		r.Group(func(r chi.Router) {
			// Rate limit: 5 requests per 15 minutes
			r.Use(RateLimitMiddleware(5, 15*time.Minute, s.config.Auth.TrustProxy))

			// Token login - returns JWT for API clients
			r.Post("/auth/token/login", authHandler.TokenLogin)

			// Token register - returns JWT for API clients
			r.Post("/auth/token/register", authHandler.TokenRegister)
		})

		// =================================================================
		// PUBLIC AUTH ROUTES WITH CSRF (for browser/SPA clients)
		// CSRF protection REQUIRED - these are state-changing operations
		// =================================================================
		r.Group(func(r chi.Router) {
			// Apply CSRF protection first (exempts Bearer tokens)
			r.Use(csrfProtection)

			// Rate limit: 5 requests per 15 minutes
			r.Use(RateLimitMiddleware(5, 15*time.Minute, s.config.Auth.TrustProxy))

			// Session-based login - returns session cookie + CSRF token
			r.Post("/auth/login", authHandler.Login)

			// Registration
			r.Post("/auth/register", authHandler.Register)
		})

		// =================================================================
		// CSRF TOKEN ENDPOINT (public, CSRF applied but permissive rate limit)
		// =================================================================
		r.Group(func(r chi.Router) {
			// Apply CSRF middleware so csrf.Token(r) works correctly
			r.Use(csrfProtection)

			// Rate limit: 10 requests per minute
			r.Use(RateLimitMiddleware(10, time.Minute, s.config.Auth.TrustProxy))

			// Get CSRF token
			r.Get("/auth/csrf", authHandler.GetCSRFToken)
		})

		// =================================================================
		// AUTHENTICATED ROUTES WITH CSRF (state-changing operations)
		// =================================================================
		r.Group(func(r chi.Router) {
			// Apply CSRF protection (exempts Bearer tokens)
			r.Use(csrfProtection)

			// Require authentication (supports both session + JWT)
			r.Use(RequireAuth(s.sessionStore, s.jwtService, s.authService))

			// Logout - MUST have CSRF protection
			r.Post("/auth/logout", authHandler.Logout)
		})

		// =================================================================
		// AUTHENTICATED READ-ONLY ROUTES (no CSRF needed)
		// GET requests don't change state, so CSRF is unnecessary
		// =================================================================
		r.Group(func(r chi.Router) {
			// Require authentication (supports both session + JWT)
			r.Use(RequireAuth(s.sessionStore, s.jwtService, s.authService))

			// Get current user
			r.Get("/auth/me", authHandler.GetCurrentUser)

			// Product read operations
			productHandler := product.NewHandler(s.db)
			r.Get("/products", productHandler.GetAll)
			r.Get("/products/{id}", productHandler.GetByID)

			// Product category read operations
			categoryHandler := product.NewCategoryHandler(s.db)
			r.Get("/products/categories", categoryHandler.GetAll)
			r.Get("/products/categories/{id}", categoryHandler.GetByID)

			// Customer read operations
			customerHandler := customer.NewHandler(s.db)
			r.Get("/customers", customerHandler.GetAll)
			r.Get("/customers/{id}", customerHandler.GetByID)
		})

		// =================================================================
		// PROTECTED PRODUCT MUTATION ROUTES (auth + CSRF required)
		// =================================================================
		r.Group(func(r chi.Router) {
			// Apply CSRF protection (exempts Bearer tokens)
			r.Use(csrfProtection)

			// Require authentication (supports both session + JWT)
			r.Use(RequireAuth(s.sessionStore, s.jwtService, s.authService))

			// Product mutations
			productHandler := product.NewHandler(s.db)
			r.Post("/products", productHandler.Create)
			r.Put("/products/{id}", productHandler.Update)
			r.Delete("/products/{id}", productHandler.Delete)

			// Product category mutations
			categoryHandler := product.NewCategoryHandler(s.db)
			r.Post("/products/categories", categoryHandler.Create)
			r.Put("/products/categories/{id}", categoryHandler.Update)
			r.Delete("/products/categories/{id}", categoryHandler.Delete)

			// Customer mutations
			customerHandler := customer.NewHandler(s.db)
			r.Post("/customers", customerHandler.Create)
			r.Put("/customers/{id}", customerHandler.Update)
			r.Delete("/customers/{id}", customerHandler.Delete)
		})
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
