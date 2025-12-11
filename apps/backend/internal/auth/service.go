package auth

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, *entities.UserAccount, error)
	ValidateToken(tokenString string) (*entities.UserAccount, error)
	GetOIDCAuthURL(ctx context.Context) (string, error)
	LoginOIDC(ctx context.Context, code string) (string, *entities.UserAccount, error)
}
