package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type UserRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Create(ctx context.Context, u entities.User) (*entities.User, error)
	Update(ctx context.Context, id uuid.UUID, u entities.User) (*entities.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error

	ListUsers(ctx context.Context, limit, offset int) ([]entities.User, int, error)
	GetByPersonID(ctx context.Context, personID uuid.UUID) (*entities.User, error)
}

type PersonRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Person, error)
	GetByEmail(ctx context.Context, email string) (*entities.Person, error)
	Search(ctx context.Context, query string, limit int) ([]entities.Person, error)
}

type EventStore interface {
	LoadUserApprovedEvents(ctx context.Context, userID uuid.UUID) ([]events.Event, error)
}

type ProjectMembershipRepository interface {
	ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error)
	PersonIDsForProjects(ctx context.Context, projectIDs []uuid.UUID) ([]uuid.UUID, error)
}
