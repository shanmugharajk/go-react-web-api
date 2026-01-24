package inventory

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// BatchHandler handles HTTP requests for product batches.
type BatchHandler struct {
	service *BatchService
}

// NewBatchHandler creates a new product batch handler.
func NewBatchHandler(database *db.DB) *BatchHandler {
	repo := NewBatchRepository(database)
	service := NewBatchService(repo)
	return &BatchHandler{service: service}
}

// GetAll handles retrieving all product batches.
func (h *BatchHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	batches, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve batches")
		return
	}

	response.Success(w, batches)
}

// GetByID handles retrieving a product batch by ID.
func (h *BatchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid batch ID")
		return
	}

	batch, err := h.service.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "batch not found")
		return
	}

	response.Success(w, batch)
}

// GetByProductID handles retrieving all batches for a specific product.
func (h *BatchHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	batches, err := h.service.GetByProductID(productID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve batches")
		return
	}

	response.Success(w, batches)
}

// Create handles creating a new product batch.
func (h *BatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	batch, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create batch")
		return
	}

	response.Created(w, batch)
}

// Update handles updating a product batch.
func (h *BatchHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid batch ID")
		return
	}

	var req UpdateBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	batch, err := h.service.Update(id, req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update batch")
		return
	}

	response.Success(w, batch)
}

// Delete handles deleting a product batch.
func (h *BatchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid batch ID")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	if err := h.service.Delete(id, user); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete batch")
		return
	}

	response.NoContent(w)
}
