package customer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/response"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/validator"
)

// Handler handles HTTP requests for customers.
type Handler struct {
	service *CustomerService
}

// NewHandler creates a new Handler instance.
func NewHandler(database *db.DB) *Handler {
	repo := NewCustomerRepository(database)
	service := NewCustomerService(repo)
	return &Handler{service: service}
}

// Routes returns the customer routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// GetAll retrieves all customers.
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch customers")
		return
	}
	response.Success(w, customers)
}

// GetByID retrieves a customer by ID.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.GetByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Customer not found")
		return
	}

	response.Success(w, customer)
}

// Create creates a new customer.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	customer, err := h.service.Create(req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to create customer")
		}
		return
	}

	response.Created(w, customer)
}

// Update updates an existing customer.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var req UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	customer, err := h.service.Update(id, req, user)
	if err != nil {
		if validator.IsValidationError(err) {
			response.Error(w, http.StatusBadRequest, err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to update customer")
		}
		return
	}

	response.Success(w, customer)
}

// Delete deletes a customer (soft delete).
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete customer")
		return
	}

	response.NoContent(w)
}
