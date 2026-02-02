package auth

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type CurrentUserProvider interface {
	Current(ctx context.Context) (entities.Principal, error)
}

// Repository is the auth-specific data access port.
type Repository interface {
	// GetAccountByEmail returns the full user account (user + person) used for claims.
	GetAccountByEmail(ctx context.Context, email string) (*entities.UserAccount, error)

	// GetPasswordHash returns the stored password hash for a user.
	// Empty string means "no password set" (OAuth-only).
	GetPasswordHash(ctx context.Context, userID uuid.UUID) (string, error)

	// GetAccountByID returns the full account for token validation.
	GetAccountByID(ctx context.Context, userID uuid.UUID) (*entities.UserAccount, error)
}
