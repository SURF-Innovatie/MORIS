package organisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entorgnode "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	entclosure "github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	"github.com/google/uuid"
	"github.com/samber/lo"
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

func (r *EntRepo) node(ctx context.Context) *ent.OrganisationNodeClient {
	if tx, ok := enttx.TxFromContext(ctx); ok {
		return tx.OrganisationNode
	}
	return r.cli.OrganisationNode
}

func (r *EntRepo) closure(ctx context.Context) *ent.OrganisationNodeClosureClient {
	if tx, ok := enttx.TxFromContext(ctx); ok {
		return tx.OrganisationNodeClosure
	}
	return r.cli.OrganisationNodeClosure
}

func (r *EntRepo) CreateNode(ctx context.Context, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error) {
	create := r.node(ctx).Create().SetName(name).SetNillableRorID(rorID).SetNillableDescription(description).SetNillableAvatarURL(avatarURL)

	if parentID != nil {
		parent, err := r.node(ctx).Get(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		create = create.SetParent(parent)
	}

	row, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (r *EntRepo) GetNode(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error) {
	row, err := r.node(ctx).
		Query().
		Where(entorgnode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (r *EntRepo) UpdateNode(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID, rorID *string, description *string, avatarURL *string) (*entities.OrganisationNode, error) {
	upd := r.node(ctx).UpdateOneID(id).SetName(name).SetNillableRorID(rorID).SetNillableDescription(description).SetNillableAvatarURL(avatarURL)

	if parentID == nil {
		upd = upd.ClearParent()
	} else {
		parent, err := r.node(ctx).Get(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		upd = upd.SetParent(parent)
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[entities.OrganisationNode](row), nil
}

func (r *EntRepo) ListRoots(ctx context.Context) ([]entities.OrganisationNode, error) {
	rows, err := r.node(ctx).
		Query().
		Where(entorgnode.Not(entorgnode.HasParent())).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNode](rows), nil
}

func (r *EntRepo) ListChildren(ctx context.Context, parentID uuid.UUID) ([]entities.OrganisationNode, error) {
	rows, err := r.node(ctx).
		Query().
		Where(entorgnode.HasParentWith(entorgnode.IDEQ(parentID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNode](rows), nil
}

func (r *EntRepo) ListAll(ctx context.Context) ([]entities.OrganisationNode, error) {
	rows, err := r.node(ctx).
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNode](rows), nil
}

func (r *EntRepo) InsertClosure(ctx context.Context, ancestorID, descendantID uuid.UUID, depth int) error {
	_, err := r.closure(ctx).
		Create().
		SetAncestorID(ancestorID).
		SetDescendantID(descendantID).
		SetDepth(depth).
		Save(ctx)
	return err
}

func (r *EntRepo) ListClosuresByDescendant(ctx context.Context, descendantID uuid.UUID) ([]entities.OrganisationNodeClosure, error) {
	rows, err := r.closure(ctx).
		Query().
		Where(entclosure.DescendantIDEQ(descendantID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNodeClosure](rows), nil
}

func (r *EntRepo) ListClosuresByAncestor(ctx context.Context, ancestorID uuid.UUID) ([]entities.OrganisationNodeClosure, error) {
	rows, err := r.closure(ctx).
		Query().
		Where(entclosure.AncestorIDEQ(ancestorID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.OrganisationNodeClosure](rows), nil
}

func (r *EntRepo) DeleteClosures(ctx context.Context, ancestorIDs, descendantIDs []uuid.UUID) error {
	if len(ancestorIDs) == 0 || len(descendantIDs) == 0 {
		return nil
	}

	_, err := r.closure(ctx).
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

	bulk := lo.Map(rows, func(c entities.OrganisationNodeClosure, _ int) *ent.OrganisationNodeClosureCreate {
		return r.closure(ctx).
			Create().
			SetAncestorID(c.AncestorID).
			SetDescendantID(c.DescendantID).
			SetDepth(c.Depth)
	})
	_, err := r.closure(ctx).CreateBulk(bulk...).Save(ctx)
	return err
}

func (r *EntRepo) Search(ctx context.Context, query string, limit int) ([]entities.OrganisationNode, error) {
	rows, err := r.node(ctx).
		Query().
		Where(entorgnode.NameContainsFold(query)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntities[entities.OrganisationNode](rows), nil
}
