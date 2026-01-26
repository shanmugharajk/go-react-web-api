package receiving

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/inventory"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/purchase"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/vendor"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/errors"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// Handler handles HTTP requests for stock receipts.
type Handler struct {
	service *StockReceiptService
}

// NewHandler creates a new Handler instance.
func NewHandler(database *db.DB) *Handler {
	repo := NewStockReceiptRepository(database.DB)
	poRepo := purchase.NewPurchaseOrderRepository(database.DB)
	vendorRepo := vendor.NewVendorRepository(database.DB)
	batchRepo := inventory.NewBatchRepository(database.DB)
	service := NewStockReceiptService(repo, poRepo, vendorRepo, batchRepo, database.DB)
	return &Handler{service: service}
}

// Routes returns the stock receipt routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/purchase-order/{purchaseOrderId}", h.GetByPurchaseOrderID)
	return r
}

// GetAll retrieves all stock receipts.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	receipts, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch stock receipts")
		return
	}
	response.Success(w, receipts)
}

// GetByID retrieves a stock receipt by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid stock receipt ID")
		return
	}

	receipt, err := h.service.GetByID(id)
	if err != nil {
		if errors.IsNotFound(err) {
			response.Error(w, http.StatusNotFound, "Stock receipt not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to fetch stock receipt")
		return
	}

	response.Success(w, receipt)
}

// GetByPurchaseOrderID retrieves all stock receipts for a purchase order.
func (h *Handler) GetByPurchaseOrderID(w http.ResponseWriter, r *http.Request) {
	poIDStr := chi.URLParam(r, "purchaseOrderId")
	poID, err := uuid.Parse(poIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid purchase order ID")
		return
	}

	receipts, err := h.service.GetByPurchaseOrderID(poID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch stock receipts")
		return
	}

	response.Success(w, receipts)
}

// Create creates a new stock receipt.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateStockReceiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	receipt, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		switch err {
		case ErrPurchaseOrderNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case ErrPurchaseOrderItemNotFound:
			response.Error(w, http.StatusBadRequest, err.Error())
		case ErrQuantityExceedsOrdered:
			response.Error(w, http.StatusBadRequest, err.Error())
		case ErrOrderNotReceivable:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "Failed to create stock receipt")
		}
		return
	}

	response.Created(w, receipt)
}
