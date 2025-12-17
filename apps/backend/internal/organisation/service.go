package organisation

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	orgnode "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
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

func (s *service) CreateRoot(ctx context.Context, name string) (*entities.OrganisationNode, error) {
	row, err := s.cli.OrganisationNode.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.OrganisationNode{}).FromEnt(row), nil
}

func (s *service) CreateChild(ctx context.Context, parentID uuid.UUID, name string) (*entities.OrganisationNode, error) {
	parent, err := s.cli.OrganisationNode.Get(ctx, parentID)
	if err != nil {
		return nil, err
	}

	row, err := s.cli.OrganisationNode.
		Create().
		SetName(name).
		SetParent(parent).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.OrganisationNode{}).FromEnt(row), nil
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.OrganisationNode, error) {
	row, err := s.cli.OrganisationNode.
		Query().
		Where(orgnode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.OrganisationNode{}).FromEnt(row), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, name string, parentID *uuid.UUID) (*entities.OrganisationNode, error) {
	upd := s.cli.OrganisationNode.UpdateOneID(id).SetName(name)

	// parent change support:
	// - parentID == nil => make it a root (clear parent)
	// - parentID != nil => set new parent
	if parentID == nil {
		upd = upd.ClearParent()
	} else {
		if *parentID == id {
			return nil, fmt.Errorf("node cannot be its own parent")
		}
		parent, err := s.cli.OrganisationNode.Get(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		upd = upd.SetParent(parent)
	}

	row, err := upd.Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.OrganisationNode{}).FromEnt(row), nil
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
