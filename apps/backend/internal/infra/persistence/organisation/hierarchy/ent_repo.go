package hierarchy

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) AncestorIDs(ctx context.Context, nodeID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.cli.OrganisationNodeClosure.
		Query().
		Where(entclosure.DescendantIDEQ(nodeID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.AncestorID)
	}
	return ids, nil
}

func (r *EntRepo) IsAncestor(ctx context.Context, ancestorID, descendantID uuid.UUID) (bool, error) {
	return r.cli.OrganisationNodeClosure.
		Query().
		Where(
			entclosure.AncestorIDEQ(ancestorID),
			entclosure.DescendantIDEQ(descendantID),
		).
		Exist(ctx)
}
