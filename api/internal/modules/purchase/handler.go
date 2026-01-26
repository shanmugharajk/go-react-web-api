package purchase

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

// Handler handles HTTP requests for purchase orders.
type Handler struct {
	service *PurchaseOrderService
}

// NewHandler creates a new Handler instance.
func NewHandler(database *db.DB) *Handler {
	repo := NewPurchaseOrderRepository(database.DB)
	service := NewPurchaseOrderService(repo)
	return &Handler{service: service}
}

// Routes returns the purchase order routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	r.Get("/vendor/{vendorId}", h.GetByVendorID)
	return r
}

// GetAll retrieves all purchase orders.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch purchase orders")
		return
	}
	response.Success(w, orders)
}

// GetByID retrieves a purchase order by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid purchase order ID")
		return
	}

	order, err := h.service.GetByID(id)
	if err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Purchase order not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to fetch purchase order")
		return
	}

	response.Success(w, order)
}

// GetByVendorID retrieves all purchase orders for a vendor.
func (h *Handler) GetByVendorID(w http.ResponseWriter, r *http.Request) {
	vendorIDStr := chi.URLParam(r, "vendorId")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vendor ID")
		return
	}

	orders, err := h.service.GetByVendorID(vendorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch purchase orders")
		return
	}

	response.Success(w, orders)
}

// Create creates a new purchase order.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	order, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create purchase order")
		return
	}

	response.Created(w, order)
}

// Update updates an existing purchase order.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid purchase order ID")
		return
	}

	var req UpdatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	order, err := h.service.Update(id, req, user)
	if err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Purchase order not found")
			return
		}
		if err == ErrCannotUpdateNonDraft {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update purchase order")
		return
	}

	response.Success(w, order)
}

// Delete cancels a purchase order (soft delete).
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid purchase order ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Purchase order not found")
			return
		}
		if err == ErrCannotDeleteNonDraft {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete purchase order")
		return
	}

	response.NoContent(w)
}
