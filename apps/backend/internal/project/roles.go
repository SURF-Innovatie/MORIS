package project

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type RoleService interface {
	EnsureDefaults(ctx context.Context) error
	List(ctx context.Context) ([]entities.ProjectRole, error)
	GetByKey(ctx context.Context, key string) (*entities.ProjectRole, error)
}

type roleService struct {
	cli *ent.Client
}

func NewRoleService(cli *ent.Client) RoleService {
	return &roleService{cli: cli}
}

func (s *roleService) EnsureDefaults(ctx context.Context) error {
	type def struct {
		key  string
		name string
	}
	defs := []def{
		{key: "contributor", name: "Contributor"},
		{key: "lead", name: "Project lead"},
	}

	for _, d := range defs {
		// try to fetch first
		_, err := s.cli.ProjectRole.Query().
			Where(entprojectrole.KeyEQ(d.key)).
			Only(ctx)
		if err != nil {
			// if not found, create it
			if ent.IsNotFound(err) {
				_, createErr := s.cli.ProjectRole.
					Create().
					SetKey(d.key).
					SetName(d.name).
					Save(ctx)
				if createErr != nil {
					return fmt.Errorf("create project role %q: %w", d.key, createErr)
				}
				continue
			}
			return fmt.Errorf("query project role %q: %w", d.key, err)
		}
		// optional: update name if changed
		_, _ = s.cli.ProjectRole.
			Update().
			Where(entprojectrole.KeyEQ(d.key)).
			SetName(d.name).
			Save(ctx)
	}

	return nil
}

func (s *roleService) List(ctx context.Context) ([]entities.ProjectRole, error) {
	rows, err := s.cli.ProjectRole.Query().Order(ent.Asc(entprojectrole.FieldKey)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.ProjectRole, 0, len(rows))
	for _, r := range rows {
		out = append(out, entities.ProjectRole{ID: r.ID, Key: r.Key, Name: r.Name})
	}
	return out, nil
}

func (s *roleService) GetByKey(ctx context.Context, key string) (*entities.ProjectRole, error) {
	r, err := s.cli.ProjectRole.Query().Where(entprojectrole.KeyEQ(key)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectRole{ID: r.ID, Key: r.Key, Name: r.Name}, nil
}
