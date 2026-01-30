package hierarchy

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error)
	IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error)
}
