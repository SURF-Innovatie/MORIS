package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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
