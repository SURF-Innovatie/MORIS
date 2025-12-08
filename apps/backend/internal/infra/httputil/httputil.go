package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	Code    int         `json:"code" example:"400"`
	Status  string      `json:"status" example:"Bad Request"`
	Errors  interface{} `json:"errors,omitempty"`                                       // Can be map[string]string or []string or null
	Message string      `json:"message,omitempty" example:"Detailed error description"` // Custom message
}

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, code int, message string, errs interface{}) error {
	resp := BackendError{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
		Errors:  errs,
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return false
	}
	return true
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
