package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/iputil"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/jwt"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/sessions"
)

// Handler handles HTTP requests for auth.
type Handler struct {
	service      *Service
	sessionStore *sessions.SQLiteStore
	jwtService   *jwt.TokenService
	isDev        bool
	sessionTTL   time.Duration
	jwtTTL       time.Duration
	trustProxy   bool
}

// NewHandler creates a new auth handler.
func NewHandler(database *db.DB, sessionStore *sessions.SQLiteStore, jwtService *jwt.TokenService, isDev bool, sessionTTL, jwtTTL time.Duration, trustProxy bool) *Handler {
	repo := NewRepository(database)
	service := NewService(repo)
	return &Handler{
		service:      service,
		sessionStore: sessionStore,
		jwtService:   jwtService,
		isDev:        isDev,
		sessionTTL:   sessionTTL,
		jwtTTL:       jwtTTL,
		trustProxy:   trustProxy,
	}
}

// Routes returns the public auth routes (no authentication required).
func (h *Handler) PublicRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Get("/csrf", h.GetCSRFToken)
	return r
}

// ProtectedRoutes returns the protected auth routes (authentication required).
func (h *Handler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/me", h.GetCurrentUser)
	r.Post("/logout", h.Logout)
	return r
}

// Login handles login requests and creates a session.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get client IP and user agent for audit logging
	ip := iputil.ExtractClientIP(r, h.trustProxy)
	userAgent := r.UserAgent()

	// Authenticate user
	user, err := h.service.Login(req, ip, userAgent)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Create session
	session, err := h.sessionStore.Create(int64(user.ID), h.sessionTTL)
	if err != nil {
		logger.Error("Failed to create session", "error", err, "user_id", user.ID)
		response.Error(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	// Set HttpOnly session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   !h.isDev, // true in production, false in dev
		SameSite: http.SameSiteLaxMode, // Lax allows top-level navigation
		MaxAge:   int(h.sessionTTL.Seconds()),
	})

	// Return user and CSRF token
	loginResp := LoginResponse{
		User:      user,
		CSRFToken: csrf.Token(r),
	}

	response.Success(w, loginResp)
}

// Register handles registration requests.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get client IP and user agent for audit logging
	ip := iputil.ExtractClientIP(r, h.trustProxy)
	userAgent := r.UserAgent()

	user, err := h.service.Register(req, ip, userAgent)
	if err != nil {
		// Check if it's a password strength error
		if err.Error() == "password validation failed: password must be at least 12 characters" ||
		   err.Error() == "password validation failed: password must not exceed 128 characters" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		// Generic error for everything else (prevents enumeration)
		response.Error(w, http.StatusBadRequest, "registration failed")
		return
	}

	response.Created(w, user)
}

// Logout handles logout requests and destroys the session.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		// No session cookie, but return success anyway
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Delete session from store
	if err := h.sessionStore.Delete(cookie.Value); err != nil {
		logger.Error("Failed to delete session", "error", err)
	}

	// Get user info from context for audit log
	ip := iputil.ExtractClientIP(r, h.trustProxy)
	logger.Info("Logout successful", "ip", ip)

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   !h.isDev,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Delete cookie
	})

	w.WriteHeader(http.StatusNoContent)
}

// GetCSRFToken returns the CSRF token for the current session.
func (h *Handler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	response.Success(w, map[string]string{
		"csrf_token": token,
	})
}

// GetCurrentUser returns the currently authenticated user.
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (injected by auth middleware)
	userID, err := GetUserIDFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Fetch user from database
	user, err := h.service.GetUserByID(uint(userID))
	if err != nil {
		logger.Error("Failed to get user", "error", err, "user_id", userID)
		response.Error(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	response.Success(w, user)
}

// TokenLogin handles login requests for API clients and returns a JWT token.
// This endpoint is designed for non-browser clients (Postman, mobile apps, CLI tools, etc.)
// that prefer token-based authentication over cookie-based sessions.
func (h *Handler) TokenLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get client IP and user agent for audit logging
	ip := iputil.ExtractClientIP(r, h.trustProxy)
	userAgent := r.UserAgent()

	// Authenticate user
	user, err := h.service.Login(req, ip, userAgent)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate JWT token
	token, err := h.jwtService.Generate(int64(user.ID), h.jwtTTL)
	if err != nil {
		logger.Error("Failed to generate JWT token", "error", err, "user_id", user.ID)
		response.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Return token response (RFC 6750 - Bearer Token Usage)
	tokenResp := TokenLoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(h.jwtTTL.Seconds()),
		User:        user,
	}

	response.Success(w, tokenResp)
}
