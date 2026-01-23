package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// Handler handles HTTP requests for products.
type Handler struct {
	service *ProductService
}

// NewHandler creates a new product handler.
func NewHandler(database *db.DB) *Handler {
	repo := NewProductRepository(database)
	service := NewProductService(repo)
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

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	product, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
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

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	product, err := h.service.Update(uint(id), req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
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

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	if err := h.service.Delete(uint(id), user); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete product")
		return
	}

	response.NoContent(w)
}

// CategoryHandler handles HTTP requests for product categories.
type CategoryHandler struct {
	service *CategoryService
}

// NewCategoryHandler creates a new product category handler.
func NewCategoryHandler(database *db.DB) *CategoryHandler {
	repo := NewCategoryRepository(database)
	service := NewCategoryService(repo)
	return &CategoryHandler{service: service}
}

// Routes returns the product category routes.
func (h *CategoryHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// GetAll handles retrieving all product categories.
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve product categories")
		return
	}

	response.Success(w, categories)
}

// GetByID handles retrieving a product category by ID.
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	category, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(w, http.StatusNotFound, "product category not found")
		return
	}

	response.Success(w, category)
}

// Create handles creating a new product category.
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	category, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create product category")
		return
	}

	response.Created(w, category)
}

// Update handles updating a product category.
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	var req UpdateProductCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	category, err := h.service.Update(uint(id), req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update product category")
		return
	}

	response.Success(w, category)
}

// Delete handles deleting a product category.
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	if err := h.service.Delete(uint(id), user); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete product category")
		return
	}

	response.NoContent(w)
}
