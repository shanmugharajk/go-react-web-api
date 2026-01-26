package payment

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/purchase"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/vendor"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/errors"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// Handler handles HTTP requests for vendor payments.
type Handler struct {
	service *VendorPaymentService
}

// NewHandler creates a new Handler instance.
func NewHandler(database *db.DB) *Handler {
	repo := NewVendorPaymentRepository(database.DB)
	vendorRepo := vendor.NewVendorRepository(database.DB)
	poRepo := purchase.NewPurchaseOrderRepository(database.DB)
	service := NewVendorPaymentService(repo, vendorRepo, poRepo, database.DB)
	return &Handler{service: service}
}

// Routes returns the vendor payment routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/vendor/{vendorId}", h.GetByVendorID)
	return r
}

// GetAll retrieves all vendor payments.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	payments, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch vendor payments")
		return
	}
	response.Success(w, payments)
}

// GetByID retrieves a vendor payment by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.service.GetByID(id)
	if err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Payment not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to fetch payment")
		return
	}

	response.Success(w, payment)
}

// GetByVendorID retrieves all payments for a vendor.
func (h *Handler) GetByVendorID(w http.ResponseWriter, r *http.Request) {
	vendorIDStr := chi.URLParam(r, "vendorId")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid vendor ID")
		return
	}

	payments, err := h.service.GetByVendorID(vendorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	response.Success(w, payments)
}

// Create creates a new vendor payment.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateVendorPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	payment, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		switch err {
		case ErrVendorNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case ErrPaymentExceedsBalance:
			response.Error(w, http.StatusBadRequest, err.Error())
		case ErrInsufficientBalance:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Failed to create payment")
		}
		return
	}

	response.Created(w, payment)
}
