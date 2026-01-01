package queries

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type ProjectReadRepository interface {
	PeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Person, error)
	ProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.ProjectRole, error)
	ProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]entities.Product, error)
	OrganisationNodeByID(ctx context.Context, id uuid.UUID) (entities.OrganisationNode, error)

	ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error)

	ListProjectRoles(ctx context.Context) ([]entities.ProjectRole, error)
	ProjectIDsStarted(ctx context.Context) ([]uuid.UUID, error)
}

type EventStore interface {
	Load(ctx context.Context, projectID uuid.UUID) ([]events.Event, int, error)
}
