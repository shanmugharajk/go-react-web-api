package vendor

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/errors"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// Handler handles HTTP requests for vendors.
type Handler struct {
	service *VendorService
}

// NewHandler creates a new Handler instance.
func NewHandler(database *db.DB) *Handler {
	repo := NewVendorRepository(database.DB)
	service := NewVendorService(repo)
	return &Handler{service: service}
}

// Routes returns the vendor routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// GetAll retrieves all vendors.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	vendors, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch vendors")
		return
	}
	response.Success(w, vendors)
}

// GetByID retrieves a vendor by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vendor ID")
		return
	}

	vendor, err := h.service.GetByID(id)
	if err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Vendor not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to fetch vendor")
		return
	}

	response.Success(w, vendor)
}

// Create creates a new vendor.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	vendor, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to create vendor")
		}
		return
	}

	response.Created(w, vendor)
}

// Update updates an existing vendor.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vendor ID")
		return
	}

	var req UpdateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	vendor, err := h.service.Update(id, req, user)
	if err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Vendor not found")
			return
		}
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update vendor")
		return
	}

	response.Success(w, vendor)
}

// Delete deletes a vendor (soft delete).
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vendor ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Vendor not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete vendor")
		return
	}

	response.NoContent(w)
}
