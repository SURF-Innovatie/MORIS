package projectrole

import (
	"context"
	"errors"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type entRepo struct {
	cli *ent.Client
}

func NewRepository(cli *ent.Client) projectrole.Repository {
	return &entRepo{cli: cli}
}

func (e *entRepo) Create(ctx context.Context, key, name string, orgNodeID uuid.UUID) (*entities.ProjectRole, error) {
	r, err := e.cli.ProjectRole.Create().
		SetKey(key).
		SetName(name).
		SetOrganisationNodeID(orgNodeID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectRole{
		ID:                 r.ID,
		Key:                r.Key,
		Name:               r.Name,
		OrganisationNodeID: r.OrganisationNodeID,
	}, nil
}

func (e *entRepo) GetByKeyAndOrg(ctx context.Context, key string, orgNodeID uuid.UUID) (*entities.ProjectRole, error) {
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
	return &entities.ProjectRole{
		ID:                 r.ID,
		Key:                r.Key,
		Name:               r.Name,
		OrganisationNodeID: r.OrganisationNodeID,
	}, nil
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

func (e *entRepo) ListByOrgIDs(ctx context.Context, orgIDs []uuid.UUID) ([]entities.ProjectRole, error) {
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

func (e *entRepo) List(ctx context.Context) ([]entities.ProjectRole, error) {
	rows, err := e.cli.ProjectRole.Query().
		Where(entprojectrole.ArchivedAtIsNil()).
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

