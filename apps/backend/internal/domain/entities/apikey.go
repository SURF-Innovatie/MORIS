package entities

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents a user's personal API key for external tool access
type APIKey struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Name       string // User-defined label, e.g., "Power BI"
	KeyPrefix  string // First 8 chars for identification (shown to user)
	CreatedAt  time.Time
	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	IsActive   bool
}

// APIKeyWithSecret is returned only on creation, includes the plain key
type APIKeyWithSecret struct {
	APIKey
	PlainKey string // The actual key, only shown once on creation
}

// IsExpired checks if the API key has expired
func (k *APIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*k.ExpiresAt)
}

// IsValid checks if the API key is active and not expired
func (k *APIKey) IsValid() bool {
	return k.IsActive && !k.IsExpired()
}
