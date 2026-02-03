package role

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/role"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	role2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
)

type entRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) role.Repository {
	return &entRepo{cli: cli}
}

func (e *entRepo) Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*role2.ProjectRole, error) {
	r, err := e.cli.ProjectRole.Create().
		SetKey(key).
		SetName(name).
		SetOrganisationNodeID(orgNodeID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[role2.ProjectRole](r), nil
}

func (e *entRepo) GetByKeyAndOrg(ctx context.Context, key string, orgNodeID uuid.UUID) (*role2.ProjectRole, error) {
	r, err := e.cli.ProjectRole.Query().
		Where(
			entprojectrole.KeyEQ(key),
			entprojectrole.OrganisationNodeIDEQ(orgNodeID),
			entprojectrole.ArchivedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[role2.ProjectRole](r), nil
}

func (e *entRepo) Delete(ctx context.Context, id uuid.UUID, orgNodeID uuid.UUID) error {
	// Soft delete
	n, err := e.cli.ProjectRole.Update().
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
		return errors.New("projectrole not found")
	}
	return nil
}

func (e *entRepo) ListByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]role2.ProjectRole, error) {
	rows, err := e.cli.ProjectRole.Query().
		Where(
			entprojectrole.OrganisationNodeIDIn(orgIDs...),
			entprojectrole.ArchivedAtIsNil(),
		).
		Order(ent.Asc(entprojectrole.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[role2.ProjectRole](rows), nil
}

func (e *entRepo) Exists(ctx context.Context, key string, orgNodeID uuid.UUID) (bool, error) {
	return e.cli.ProjectRole.Query().
		Where(
			entprojectrole.KeyEQ(key),
			entprojectrole.OrganisationNodeIDEQ(orgNodeID),
		).
		Exist(ctx)
}

func (e *entRepo) Unarchive(ctx context.Context, key string, orgNodeID uuid.UUID) error {
	return e.cli.ProjectRole.Update().
		Where(
			entprojectrole.KeyEQ(key),
			entprojectrole.OrganisationNodeIDEQ(orgNodeID),
			entprojectrole.ArchivedAtNotNil(),
		).
		ClearArchivedAt().
		Exec(ctx)
}

func (e *entRepo) CreateOrRestore(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*role2.ProjectRole, error) {
	// First check if it exists (including archived)
	existing, err := e.cli.ProjectRole.Query().
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
			updated, err := e.cli.ProjectRole.UpdateOne(existing).
				ClearArchivedAt().
				SetName(name).
				Save(ctx)
			if err != nil {
				return nil, err
			}
			return transform.ToEntityPtr[role2.ProjectRole](updated), nil
		}
		return nil, fmt.Errorf("role with key '%s' already exists", key)
	}

	r, err := e.cli.ProjectRole.Create().
		SetKey(key).
		SetName(name).
		SetOrganisationNodeID(orgNodeID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[role2.ProjectRole](r), nil
}

func (e *entRepo) List(ctx context.Context) ([]role2.ProjectRole, error) {
	rows, err := e.cli.ProjectRole.Query().
		Where(entprojectrole.ArchivedAtIsNil()).
		Order(ent.Asc(entprojectrole.FieldKey)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[role2.ProjectRole](rows), nil
}

func (e *entRepo) CreateWithEventTypes(ctx context.Context, key, name string, orgNodeID uuid.UUID, allowedEventTypes []string) (*role2.ProjectRole, error) {
	r, err := e.cli.ProjectRole.Create().
		SetKey(key).
		SetName(name).
		SetOrganisationNodeID(orgNodeID).
		SetAllowedEventTypes(allowedEventTypes).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[role2.ProjectRole](r), nil
}

func (e *entRepo) GetByID(ctx context.Context, id uuid.UUID) (*role2.ProjectRole, error) {
	r, err := e.cli.ProjectRole.Query().
		Where(
			entprojectrole.ID(id),
			entprojectrole.ArchivedAtIsNil(),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[role2.ProjectRole](r), nil
}

func (e *entRepo) UpdateAllowedEventTypes(ctx context.Context, id uuid.UUID, eventTypes []string) (*role2.ProjectRole, error) {
	r, err := e.cli.ProjectRole.UpdateOneID(id).
		SetAllowedEventTypes(eventTypes).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntityPtr[role2.ProjectRole](r), nil
}
