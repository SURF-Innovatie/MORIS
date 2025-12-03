package auth

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Service interface {
	Register(ctx context.Context, req userdto.Request) (*entities.UserAccount, error)
	Login(ctx context.Context, email, password string) (string, *entities.UserAccount, error)
	ValidateToken(tokenString string) (*entities.UserAccount, error)
}
