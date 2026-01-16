package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/apikey"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	// APIKeyPrefix is prepended to all generated keys
	APIKeyPrefix = "moris_"
	// APIKeyLength is the length of the random part of the key
	APIKeyLength = 32
)

// Handler handles API key management requests
type Handler struct {
	client *ent.Client
}

// NewHandler creates a new API key handler
func NewHandler(client *ent.Client) *Handler {
	return &Handler{client: client}
}

// RegisterRoutes registers API key management routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/profile/api-keys", func(r chi.Router) {
		r.Get("/", h.ListAPIKeys)
		r.Post("/", h.CreateAPIKey)
		r.Delete("/{keyId}", h.RevokeAPIKey)
	})
}

// ListAPIKeys godoc
// @Summary List user's API keys
// @Description Retrieves all API keys for the current user (secrets are not shown)
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIKeyListResponse
// @Failure 401 {string} string "unauthorized"
// @Router /profile/api-keys [get]
func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	userIDPtr := httputil.GetUserIDFromContext(r.Context())
	if userIDPtr == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userID := *userIDPtr

	keys, err := h.client.APIKey.
		Query().
		Where(apikey.UserID(userID)).
		Order(ent.Desc(apikey.FieldCreatedAt)).
		All(r.Context())

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := dto.APIKeyListResponse{
		APIKeys: make([]dto.APIKeyResponse, len(keys)),
	}
	for i, key := range keys {
		resp.APIKeys[i] = mapAPIKeyToResponse(key)
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// CreateAPIKey godoc
// @Summary Create a new API key
// @Description Creates a new API key for the current user. The secret is only shown once.
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateAPIKeyRequest true "API key creation request"
// @Success 201 {object} dto.APIKeyWithSecretResponse
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Router /profile/api-keys [post]
func (h *Handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userIDPtr := httputil.GetUserIDFromContext(r.Context())
	if userIDPtr == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userID := *userIDPtr

	var req dto.CreateAPIKeyRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	// Generate random key
	plainKey, err := generateAPIKey()
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to generate key", nil)
		return
	}

	// Hash the key for storage
	keyHash := hashAPIKey(plainKey)
	keyPrefix := plainKey[:8]

	// Create the key in database
	builder := h.client.APIKey.
		Create().
		SetUserID(userID).
		SetName(req.Name).
		SetKeyHash(keyHash).
		SetKeyPrefix(keyPrefix)

	if req.ExpiresAt != nil {
		builder = builder.SetExpiresAt(*req.ExpiresAt)
	}

	key, err := builder.Save(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Return response with the plain key (only shown once)
	resp := dto.APIKeyWithSecretResponse{
		APIKeyResponse: mapAPIKeyToResponse(key),
		PlainKey:       APIKeyPrefix + plainKey,
	}

	_ = httputil.WriteJSON(w, http.StatusCreated, resp)
}

// RevokeAPIKey godoc
// @Summary Revoke an API key
// @Description Revokes (deactivates) an API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyId path string true "API Key ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "invalid key id"
// @Failure 401 {string} string "unauthorized"
// @Failure 404 {string} string "key not found"
// @Router /profile/api-keys/{keyId} [delete]
func (h *Handler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	userIDPtr := httputil.GetUserIDFromContext(r.Context())
	if userIDPtr == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userID := *userIDPtr

	keyID, err := httputil.ParseUUIDParam(r, "keyId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid key id", nil)
		return
	}

	// Verify the key belongs to the user and deactivate it
	updated, err := h.client.APIKey.
		Update().
		Where(
			apikey.ID(keyID),
			apikey.UserID(userID),
		).
		SetIsActive(false).
		Save(r.Context())

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if updated == 0 {
		httputil.WriteError(w, r, http.StatusNotFound, "key not found", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ValidateAPIKey validates an API key and returns the user ID if valid
func (h *Handler) ValidateAPIKey(ctx http.Request, plainKey string) (uuid.UUID, error) {
	// Remove prefix if present
	if len(plainKey) > len(APIKeyPrefix) && plainKey[:len(APIKeyPrefix)] == APIKeyPrefix {
		plainKey = plainKey[len(APIKeyPrefix):]
	}

	keyHash := hashAPIKey(plainKey)

	key, err := h.client.APIKey.
		Query().
		Where(
			apikey.KeyHash(keyHash),
			apikey.IsActive(true),
		).
		Only(ctx.Context())

	if err != nil {
		return uuid.Nil, err
	}

	// Check expiration
	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return uuid.Nil, err
	}

	// Update last used timestamp
	_, _ = h.client.APIKey.
		UpdateOneID(key.ID).
		SetLastUsedAt(time.Now()).
		Save(ctx.Context())

	return key.UserID, nil
}

// Helper functions

func generateAPIKey() (string, error) {
	bytes := make([]byte, APIKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func mapAPIKeyToResponse(key *ent.APIKey) dto.APIKeyResponse {
	return dto.APIKeyResponse{
		ID:         key.ID,
		Name:       key.Name,
		KeyPrefix:  APIKeyPrefix + key.KeyPrefix + "...",
		CreatedAt:  key.CreatedAt,
		LastUsedAt: key.LastUsedAt,
		ExpiresAt:  key.ExpiresAt,
		IsActive:   key.IsActive,
	}
}
