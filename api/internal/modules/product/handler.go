package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
)

// Handler handles HTTP requests for products.
type Handler struct {
	service *Service
}

// NewHandler creates a new product handler.
func NewHandler(database *db.DB) *Handler {
	repo := NewRepository(database)
	service := NewService(repo)
	return &Handler{service: service}
}

// Routes returns the product routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// GetAll handles retrieving all products.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve products")
		return
	}

	response.Success(w, products)
}

// GetByID handles retrieving a product by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(w, http.StatusNotFound, "product not found")
		return
	}

	response.Success(w, product)
}

// Create handles creating a new product.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.service.Create(req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create product")
		return
	}

	response.Created(w, product)
}

// Update handles updating a product.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.service.Update(uint(id), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update product")
		return
	}

	response.Success(w, product)
}

// Delete handles deleting a product.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete product")
		return
	}

	response.NoContent(w)
}
