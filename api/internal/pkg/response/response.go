package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Success writes a successful JSON response.
func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

// Error writes an error JSON response.
func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, map[string]any{
		"success": false,
		"error":   message,
	})
}

// Created writes a 201 Created response.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data":    data,
	})
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
