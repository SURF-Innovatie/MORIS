package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity/readmodels"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// StatusResponse represents a standard status response
type StatusResponse struct {
	Status    string `json:"status" example:"ok"`
	Timestamp string `json:"timestamp" example:"2025-11-12T10:00:00Z"`
}

// BackendError swagger:model BackendError
// BackendError is a standardized error response structure, referenced by Swagger
// used by Swagger for API documentation
type BackendError struct {
	Code    int    `json:"code" example:"400"`
	Status  string `json:"status" example:"Bad Request"`
	Errors  any    `json:"errors,omitempty"`                                       // Can be map[string]string or []string or null
	Message string `json:"message,omitempty" example:"Detailed error description"` // Custom message
}

// ContextKey is a custom type for context keys
type ContextKey string

const (
	// ContextKeyErrorDetails is the key used to store error details in the request context
	ContextKeyErrorDetails ContextKey = "error_details"
	// ContextKeyUser is the key used to store user info in context
	ContextKeyUser ContextKey = "user"
)

// GetUserFromContext retrieves the authUser from the request context
func GetUserFromContext(ctx context.Context) (*readmodels.UserAccount, bool) {
	user, ok := ctx.Value(ContextKeyUser).(*readmodels.UserAccount)
	return user, ok
}

// GetUserIDFromContext helper to extract user ID safely
func GetUserIDFromContext(ctx context.Context) *uuid.UUID {
	userCtx, ok := GetUserFromContext(ctx)
	if !ok || userCtx == nil {
		return nil
	}
	return &userCtx.User.ID
}

// WriteError writes a standardized error response
// It automatically handles environment-specific masking of errors.
// It also stores the full error details in the request entity for middleware logging.
func WriteError(w http.ResponseWriter, r *http.Request, code int, message string, errs any) error {
	// Store full error details in the request context specifically for the middleware to pick up.
	// This allows the middleware to log the original error details even if the response is sanitized for production.
	if container, ok := r.Context().Value(ContextKeyErrorDetails).(*ErrorDetailsContainer); ok {
		container.Message = message
		container.Errors = errs
		container.StatusCode = code
	}

	resp := BackendError{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Errors:  errs,
	}

	// Environment based masking
	if env.IsProd() {
		resp.Message = http.StatusText(code) // Default to generic status text
		resp.Errors = nil                    // Hide detailed errors
	} else {
		// Log error to console in development
		log.Error().Int("code", code).Interface("errors", errs).Msg(message)
	}

	return WriteJSON(w, code, resp)
}

// WriteStatus writes a standard status response with "ok" status and current timestamp
func WriteStatus(w http.ResponseWriter) error {
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return WriteJSON(w, http.StatusOK, resp)
}

// WriteJSON writes a JSON response with the given status code and data.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// ReadJSON reads a JSON request body into v.
// It returns true if successful, or writes an error response and returns false if failed.
func ReadJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		WriteError(w, r, http.StatusBadRequest, "invalid request body", err.Error())
		return false
	}
	return true
}

// ErrorDetailsContainer is used to pass error details from handler to middleware via context
type ErrorDetailsContainer struct {
	Message    string
	Errors     any
	StatusCode int
}

// ParseUUIDParam parses a UUID from the URL path parameters.
func ParseUUIDParam(r *http.Request, key string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, key)
	if idStr == "" {
		return uuid.Nil, errors.New("missing param: " + key)
	}
	return uuid.Parse(idStr)
}

// ParseUUIDQuery parses a UUID from the URL query parameters.
func ParseUUIDQuery(r *http.Request, key string) (uuid.UUID, error) {
	idStr := r.URL.Query().Get(key)
	if idStr == "" {
		return uuid.Nil, errors.New("missing query param: " + key)
	}
	return uuid.Parse(idStr)
}

// ParseIntQuery parses an int from the URL query parameters with a default value.
func ParseIntQuery(r *http.Request, key string, defaultValue int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}
