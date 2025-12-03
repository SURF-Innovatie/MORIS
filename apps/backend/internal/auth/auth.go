package auth

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	PersonID uuid.UUID
	Password string
}

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*entities.UserAccount, error)
	Login(ctx context.Context, email, password string) (string, *entities.UserAccount, error)
	ValidateToken(tokenString string) (*entities.UserAccount, error)
}
