package organisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entorgnode "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
	tx  *ent.Tx
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func newTxRepo(tx *ent.Tx) *EntRepo {
	return &EntRepo{tx: tx}
}

func (r *EntRepo) node() *ent.OrganisationNodeClient {
	if r.tx != nil {
		return r.tx.OrganisationNode
	}
	return r.cli.OrganisationNode
}

func (r *EntRepo) closure() *ent.OrganisationNodeClosureClient {
	if r.tx != nil {
		return r.tx.OrganisationNodeClosure
	}
	return r.cli.OrganisationNodeClosure
}

func (r *EntRepo) WithTx(ctx context.Context, fn func(ctx context.Context, tx organisation.Repository) error) error {
	if r.cli == nil {
		// already in tx repo; just run
		return fn(ctx, r)
	}

	tx, err := r.cli.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := fn(ctx, newTxRepo(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *EntRepo) CreateNode(ctx context.Context, name string, parentID *uuid.UUID, rorID *string) (*entities.OrganisationNode, error) {
	create := r.node().Create().SetName(name).SetNillableRorID(rorID)

	if parentID != nil {
		parent, err := r.node().Get(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		create = create.SetParent(parent)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &entities.OrganisationNode{
		ID:       row.ID,
		ParentID: row.ParentID,
		Name:     row.Name,
		RorID:    row.RorID,
	}, nil
}

func (r *EntRepo) GetNode(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error) {
	row, err := r.node().
		Query().
		Where(entorgnode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &entities.OrganisationNode{
		ID:       row.ID,
		ParentID: row.ParentID,
		Name:     row.Name,
		RorID:    row.RorID,
	}, nil
}

func (r *EntRepo) UpdateNode(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string) (*entities.OrganisationNode, error) {
	upd := r.node().UpdateOneID(id).SetName(name).SetNillableRorID(rorID)

	if parentID == nil {
		upd = upd.ClearParent()
	} else {
		parent, err := r.node().Get(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		upd = upd.SetParent(parent)
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	return &entities.OrganisationNode{
		ID:       row.ID,
		ParentID: row.ParentID,
		Name:     row.Name,
		RorID:    row.RorID,
	}, nil
}

func (r *EntRepo) ListRoots(ctx context.Context) ([]entities.OrganisationNode, error) {
	rows, err := r.node().
		Query().
		Where(entorgnode.Not(entorgnode.HasParent())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrganisationNode, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.OrganisationNode{
			ID:       row.ID,
			ParentID: row.ParentID,
			Name:     row.Name,
			RorID:    row.RorID,
		})
	}
	return out, nil
}

func (r *EntRepo) ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error) {
	rows, err := r.node().
		Query().
		Where(entorgnode.HasParentWith(entorgnode.IDEQ(parentID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrganisationNode, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.OrganisationNode{
			ID:       row.ID,
			ParentID: row.ParentID,
			Name:     row.Name,
			RorID:    row.RorID,
		})
	}
	return out, nil
}

func (r *EntRepo) InsertClosure(ctx context.Context, ancestorID, descendantID uuid.UUID, depth int) error {
	_, err := r.closure().
		Create().
		SetAncestorID(ancestorID).
		SetDescendantID(descendantID).
		SetDepth(depth).
		Save(ctx)
	return err
}

func (r *EntRepo) ListClosuresByDescendant(ctx context.Context, descendantID uuid.UUID) ([]entities.OrganisationNodeClosure, error) {
	rows, err := r.closure().
		Query().
		Where(entclosure.DescendantIDEQ(descendantID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrganisationNodeClosure, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.OrganisationNodeClosure{
			AncestorID:   row.AncestorID,
			DescendantID: row.DescendantID,
			Depth:        row.Depth,
		})
	}
	return out, nil
}

func (r *EntRepo) ListClosuresByAncestor(ctx context.Context, ancestorID uuid.UUID) ([]entities.OrganisationNodeClosure, error) {
	rows, err := r.closure().
		Query().
		Where(entclosure.AncestorIDEQ(ancestorID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrganisationNodeClosure, 0, len(rows))
	for _, row := range rows {
		out = append(out, entities.OrganisationNodeClosure{
			AncestorID:   row.AncestorID,
			DescendantID: row.DescendantID,
			Depth:        row.Depth,
		})
	}
	return out, nil
}

func (r *EntRepo) DeleteClosures(ctx context.Context, ancestorIDs, descendantIDs []uuid.UUID) error {
	if len(ancestorIDs) == 0 || len(descendantIDs) == 0 {
		return nil
	}

	_, err := r.closure().
		Delete().
		Where(
			entclosure.AncestorIDIn(ancestorIDs...),
			entclosure.DescendantIDIn(descendantIDs...),
		).
		Exec(ctx)
	return err
}

func (r *EntRepo) CreateClosuresBulk(ctx context.Context, rows []entities.OrganisationNodeClosure) error {
	if len(rows) == 0 {
		return nil
	}

	bulk := make([]*ent.OrganisationNodeClosureCreate, 0, len(rows))
	for _, c := range rows {
		bulk = append(bulk, r.closure().
			Create().
			SetAncestorID(c.AncestorID).
			SetDescendantID(c.DescendantID).
			SetDepth(c.Depth),
		)
	}
	_, err := r.closure().CreateBulk(bulk...).Save(ctx)
	return err
}
