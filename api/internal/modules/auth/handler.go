package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
)

// Handler handles HTTP requests for auth.
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler.
func NewHandler(database *db.DB) *Handler {
	repo := NewRepository(database)
	service := NewService(repo)
	return &Handler{service: service}
}

// Routes returns the auth routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	return r
}

// Login handles login requests.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Login(req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	response.Success(w, user)
}

// Register handles registration requests.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	response.Created(w, user)
}
