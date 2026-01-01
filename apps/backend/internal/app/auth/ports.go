package auth

import (
	"context"

	"github.com/google/uuid"
)

type CurrentUser interface {
	UserID() uuid.UUID
	PersonID() uuid.UUID
}

type CurrentUserProvider interface {
	Current(ctx context.Context) (CurrentUser, error)
}
