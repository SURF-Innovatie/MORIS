package projectrole

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	EnsureDefaults(ctx context.Context) error
	Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error)
	ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]entities.ProjectRole, error)
	Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) EnsureDefaults(ctx context.Context) error {
	roots, err := s.cli.OrganisationNode.Query().
		Where(organisationnode.ParentIDIsNil()).
		All(ctx)
	if err != nil {
		return fmt.Errorf("listing root nodes: %w", err)
	}

	defs := []struct {
		key  string
		name string
	}{
		{key: "contributor", name: "Contributor"},
		{key: "lead", name: "Project lead"},
	}

	for _, root := range roots {
		for _, d := range defs {
			exists, err := s.cli.ProjectRole.Query().
				Where(
					entprojectrole.KeyEQ(d.key),
					entprojectrole.OrganisationNodeIDEQ(root.ID),
				).
				Exist(ctx)
			if err != nil {
				return fmt.Errorf("checking existence of role %s for root %s: %w", d.key, root.ID, err)
			}

			if !exists {
				_, err := s.cli.ProjectRole.Create().
					SetKey(d.key).
					SetName(d.name).
					SetOrganisationNodeID(root.ID).
					Save(ctx)
				if err != nil {
					return fmt.Errorf("creating default role %s for root %s: %w", d.key, root.ID, err)
				}
			} else {
				err := s.cli.ProjectRole.Update().
					Where(
						entprojectrole.KeyEQ(d.key),
						entprojectrole.OrganisationNodeIDEQ(root.ID),
						entprojectrole.ArchivedAtNotNil(),
					).
					ClearArchivedAt().
					Exec(ctx)
				if err != nil {
					return fmt.Errorf("ensuring default role %s is unarchived for root %s: %w", d.key, root.ID, err)
				}
			}
		}
	}

	return nil
}

func (s *service) Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error) {
	existing, err := s.cli.ProjectRole.Query().
		Where(
			entprojectrole.KeyEQ(key),
			entprojectrole.OrganisationNodeIDEQ(orgNodeID),
		).
		Only(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	if existing != nil {
		if existing.ArchivedAt != nil {
			updated, err := s.cli.ProjectRole.UpdateOne(existing).
				ClearArchivedAt().
				SetName(name).
				Save(ctx)
			if err != nil {
				return nil, err
			}
			return &entities.ProjectRole{
				ID:                 updated.ID,
				Key:                updated.Key,
				Name:               updated.Name,
				OrganisationNodeID: updated.OrganisationNodeID,
			}, nil
		}
		return nil, fmt.Errorf("role with key '%s' already exists", key)
	}

	r, err := s.cli.ProjectRole.Create().
		SetKey(key).
		SetName(name).
		SetOrganisationNodeID(orgNodeID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectRole{ID: r.ID, Key: r.Key, Name: r.Name, OrganisationNodeID: r.OrganisationNodeID}, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error {
	n, err := s.cli.ProjectRole.Update().
		Where(
			entprojectrole.ID(id),
			entprojectrole.OrganisationNodeIDEQ(orgNodeID),
		).
		SetArchivedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("role not found or not owned by organisation")
	}
	return nil
}

func (s *service) ListAvailableForNode(ctx context.Context, orgNodeID uuid.UUID) ([]entities.ProjectRole, error) {
	ancestors, err := s.cli.OrganisationNodeClosure.Query().
		Where(organisationnodeclosure.DescendantIDEQ(orgNodeID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(ancestors))
	for _, a := range ancestors {
		ids = append(ids, a.AncestorID)
	}

	rows, err := s.cli.ProjectRole.Query().
		Where(
			entprojectrole.OrganisationNodeIDIn(ids...),
			entprojectrole.ArchivedAtIsNil(),
		).
		Order(ent.Asc(entprojectrole.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.ProjectRole, 0, len(rows))
	for _, r := range rows {
		out = append(out, entities.ProjectRole{
			ID:                 r.ID,
			Key:                r.Key,
			Name:               r.Name,
			OrganisationNodeID: r.OrganisationNodeID,
		})
	}
	return out, nil
}
