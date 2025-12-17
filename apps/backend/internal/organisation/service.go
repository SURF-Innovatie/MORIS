package organisation

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	orgnode "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	orgclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	CreateRoot(ctx context.Context, name string) (*entities.OrganisationNode, error)
	CreateChild(ctx context.Context, parentID uuid.UUID, name string) (*entities.OrganisationNode, error)

	Get(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error)
	Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID) (*entities.OrganisationNode, error)

	ListRoots(ctx context.Context) ([]entities.OrganisationNode, error)
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error)
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) insertSelfClosure(ctx context.Context, tx *ent.Tx, nodeID uuid.UUID) error {
	_, err := tx.OrganisationNodeClosure.
		Create().
		SetAncestorID(nodeID).
		SetDescendantID(nodeID).
		SetDepth(0).
		Save(ctx)
	return err
}

func (s *service) insertAncestorClosures(ctx context.Context, tx *ent.Tx, parentID, childID uuid.UUID) error {
	ancRows, err := tx.OrganisationNodeClosure.
		Query().
		Where(orgclosure.DescendantIDEQ(parentID)).
		All(ctx)
	if err != nil {
		return err
	}

	bulk := make([]*ent.OrganisationNodeClosureCreate, 0, len(ancRows))
	for _, a := range ancRows {
		bulk = append(bulk,
			tx.OrganisationNodeClosure.Create().
				SetAncestorID(a.AncestorID).
				SetDescendantID(childID).
				SetDepth(a.Depth+1),
		)
	}
	if len(bulk) == 0 {
		return nil
	}
	_, err = tx.OrganisationNodeClosure.CreateBulk(bulk...).Save(ctx)
	return err
}

func (s *service) CreateRoot(ctx context.Context, name string) (*entities.OrganisationNode, error) {
	tx, err := s.cli.Tx(ctx)
	if err != nil {
		return nil, err
	}

	row, err := tx.OrganisationNode.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.insertSelfClosure(ctx, tx, row.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (s *service) CreateChild(ctx context.Context, parentID uuid.UUID, name string) (*entities.OrganisationNode, error) {
	tx, err := s.cli.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	parent, err := tx.OrganisationNode.Get(ctx, parentID)
	if err != nil {
		return nil, err
	}

	row, err := tx.OrganisationNode.
		Create().
		SetName(name).
		SetParent(parent).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.insertSelfClosure(ctx, tx, row.ID); err != nil {
		return nil, err
	}
	if err := s.insertAncestorClosures(ctx, tx, parentID, row.ID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error) {
	row, err := s.cli.OrganisationNode.
		Query().
		Where(orgnode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (s *service) subtreeDepths(ctx context.Context, tx *ent.Tx, nodeID uuid.UUID) (map[uuid.UUID]int, []uuid.UUID, error) {
	rows, err := tx.OrganisationNodeClosure.
		Query().
		Where(orgclosure.AncestorIDEQ(nodeID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	depth := make(map[uuid.UUID]int, len(rows))
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		depth[r.DescendantID] = r.Depth
		ids = append(ids, r.DescendantID)
	}
	return depth, ids, nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID) (*entities.OrganisationNode, error) {
	tx, err := s.cli.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	cur, err := tx.OrganisationNode.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	upd := tx.OrganisationNode.UpdateOneID(id).SetName(name)

	var newParent *ent.OrganisationNode
	if parentID == nil {
		upd = upd.ClearParent()
	} else {
		if *parentID == id {
			return nil, fmt.Errorf("node cannot be its own parent")
		}
		newParent, err = tx.OrganisationNode.Get(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		upd = upd.SetParent(newParent)
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	// If parent unchanged, no closure updates needed.
	oldParentID, hadOld := cur.ParentID, *cur.ParentID != uuid.Nil
	newParentID, hadNew := uuid.Nil, parentID != nil
	if hadNew {
		newParentID = *parentID
	}
	if hadOld == hadNew && (!hadOld || *oldParentID == newParentID) {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return transform.ToEntityPtr[entities.OrganisationNode](row), nil
	}

	// Subtree nodes (descendants including self) + their depth from "id"
	subDepth, subIDs, err := s.subtreeDepths(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	// Cycle check: cannot move under own subtree
	if hadNew {
		if _, ok := subDepth[newParentID]; ok {
			return nil, fmt.Errorf("cannot move node under its own subtree")
		}
	}

	oldAncRows, err := tx.OrganisationNodeClosure.
		Query().
		Where(
			orgclosure.DescendantIDEQ(id),
			orgclosure.AncestorIDNotIn(subIDs...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	oldAncIDs := make([]uuid.UUID, 0, len(oldAncRows))
	for _, a := range oldAncRows {
		oldAncIDs = append(oldAncIDs, a.AncestorID)
	}

	// Remove old external links: (old external ancestors) -> (subtree descendants)
	if len(oldAncIDs) > 0 && len(subIDs) > 0 {
		if _, err := tx.OrganisationNodeClosure.
			Delete().
			Where(
				orgclosure.AncestorIDIn(oldAncIDs...),
				orgclosure.DescendantIDIn(subIDs...),
			).
			Exec(ctx); err != nil {
			return nil, err
		}
	}

	// If new parent is nil (becomes root), we’re done after deletion.
	if !hadNew {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return transform.ToEntityPtr[entities.OrganisationNode](row), nil
	}

	// New external ancestors: ancestors of new parent (including parent)
	newAncRows, err := tx.OrganisationNodeClosure.
		Query().
		Where(orgclosure.DescendantIDEQ(newParentID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Insert new external links: for each (newAncestor) and each (descendant in subtree)
	bulk := make([]*ent.OrganisationNodeClosureCreate, 0, len(newAncRows)*len(subIDs))
	for _, na := range newAncRows {
		for _, descID := range subIDs {
			d := subDepth[descID] // distance from "id" to this descendant
			bulk = append(bulk,
				tx.OrganisationNodeClosure.Create().
					SetAncestorID(na.AncestorID).
					SetDescendantID(descID).
					SetDepth(na.Depth+1+d),
			)
		}
	}
	if len(bulk) > 0 {
		_, err = tx.OrganisationNodeClosure.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (s *service) ListRoots(ctx context.Context) ([]entities.OrganisationNode, error) {
	rows, err := s.cli.OrganisationNode.
		Query().
		Where(orgnode.Not(orgnode.HasParent())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNode](rows), nil
}

func (s *service) ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error) {
	rows, err := s.cli.OrganisationNode.
		Query().
		Where(orgnode.HasParentWith(orgnode.IDEQ(parentID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNode](rows), nil
}
