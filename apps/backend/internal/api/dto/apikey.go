package dto

import (
	"time"

	"github.com/google/uuid"
)

// APIKeyResponse represents an API key (without secret)
type APIKeyResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	IsActive   bool       `json:"isActive"`
}

// APIKeyWithSecretResponse is returned only on creation
type APIKeyWithSecretResponse struct {
	APIKeyResponse
	PlainKey string `json:"plainKey"` // Only shown once
}

// CreateAPIKeyRequest is the input for creating an API key
type CreateAPIKeyRequest struct {
	Name      string     `json:"name" binding:"required"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

// APIKeyListResponse wraps a list of API keys
type APIKeyListResponse struct {
	APIKeys []APIKeyResponse `json:"apiKeys"`
}
