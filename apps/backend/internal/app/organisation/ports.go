package organisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx Repository) error) error

	CreateNode(ctx context.Context, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error)
	GetNode(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error)
	UpdateNode(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error)

	ListRoots(ctx context.Context) ([]entities.OrganisationNode, error)
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error)
	Search(ctx context.Context, query string, limit int) ([]entities.OrganisationNode, error)

	InsertClosure(ctx context.Context, ancestorID, descendantID uuid.UUID, depth int) error
	ListClosuresByDescendant(ctx context.Context, descendantID uuid.UUID) ([]entities.OrganisationNodeClosure, error)
	ListClosuresByAncestor(ctx context.Context, ancestorID uuid.UUID) ([]entities.OrganisationNodeClosure, error)

	DeleteClosures(ctx context.Context, ancestorIDs, descendantIDs []uuid.UUID) error
	CreateClosuresBulk(ctx context.Context, rows []entities.OrganisationNodeClosure) error
}

type PersonRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Person, error)
	Update(ctx context.Context, id uuid.UUID, p entities.Person) (*entities.Person, error)
}
