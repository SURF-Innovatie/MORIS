package auth

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, *entities.UserAccount, error)
	// LoginByEmail issues a MORIS JWT for an existing user identified by email.
	// Used for SSO/OIDC flows where we don't have a local password.
	LoginByEmail(ctx context.Context, email string) (string, *entities.UserAccount, error)
	ValidateToken(tokenString string) (*entities.UserAccount, error)
}
