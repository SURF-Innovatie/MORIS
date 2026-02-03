package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*identity.User, error)
	Create(ctx context.Context, u identity.User) (*identity.User, error)
	Update(ctx context.Context, id uuid.UUID, u identity.User) (*identity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error

	ListUsers(ctx context.Context, limit, offset int) ([]identity.User, int, error)
	GetByPersonID(ctx context.Context, personID uuid.UUID) (*identity.User, error)
	SetZenodoTokens(ctx context.Context, userID uuid.UUID, access, refresh string) error
	ClearZenodoTokens(ctx context.Context, userID uuid.UUID) error
}

type ProjectMembershipRepository interface {
	ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error)
	PersonIDsForProjects(ctx context.Context, projectIDs []uuid.UUID) ([]uuid.UUID, error)
}
