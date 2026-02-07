package organisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/readmodels"
	"github.com/google/uuid"
)

type repository interface {
	CreateNode(ctx context.Context, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string, slug string) (*organisation.OrganisationNode, error)
	GetNode(ctx context.Context, id uuid.UUID) (*organisation.OrganisationNode, error)
	UpdateNode(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*organisation.OrganisationNode, error)

	ListRoots(ctx context.Context) ([]organisation.OrganisationNode, error)
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]organisation.OrganisationNode, error)
	ListAll(ctx context.Context) ([]organisation.OrganisationNode, error)
	Search(ctx context.Context, query string, limit int) ([]organisation.OrganisationNode, error)

	InsertClosure(ctx context.Context, ancestorID, descendantID uuid.UUID, depth int) error
	ListClosuresByDescendant(ctx context.Context, descendantID uuid.UUID) ([]readmodels.OrganisationNodeClosure, error)
	ListClosuresByAncestor(ctx context.Context, ancestorID uuid.UUID) ([]readmodels.OrganisationNodeClosure, error)

	DeleteClosures(ctx context.Context, ancestorIDs, descendantIDs []uuid.UUID) error
	CreateClosuresBulk(ctx context.Context, rows []readmodels.OrganisationNodeClosure) error
}
