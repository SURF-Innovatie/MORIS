package queries

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

type ProjectReadRepository interface {
	PeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error)
	ProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]role.ProjectRole, error)
	ProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]product.Product, error)
	OrganisationNodeByID(ctx context.Context, id uuid.UUID) (organisation.OrganisationNode, error)

	ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error)

	ProjectIDsStarted(ctx context.Context) ([]uuid.UUID, error)
	ListAncestors(ctx context.Context, orgID uuid.UUID) ([]uuid.UUID, error)
	ProjectIDBySlug(ctx context.Context, slug string) (uuid.UUID, error)
}

type EventStore interface {
	Load(ctx context.Context, projectID uuid.UUID) ([]events.Event, int, error)
}

type ProjectRoleRepository interface {
	List(ctx context.Context) ([]role.ProjectRole, error)
	ListByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]role.ProjectRole, error)
}
