package projectrole

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type Repository interface {
	Upsert(ctx context.Context, key, name string) error
	List(ctx context.Context) ([]entities.ProjectRole, error)
	GetByKey(ctx context.Context, key string) (*entities.ProjectRole, error)
}

type entRepo struct {
	cli *ent.Client
}

func NewRepository(cli *ent.Client) Repository {
	return &entRepo{cli: cli}
}

func (e *entRepo) Upsert(ctx context.Context, key, name string) error {
	// If exists, update name; else create
	row, err := e.cli.ProjectRole.Query().Where(entprojectrole.KeyEQ(key)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			_, err := e.cli.ProjectRole.Create().SetKey(key).SetName(name).Save(ctx)
			return err
		}
		return err
	}
	_, err = e.cli.ProjectRole.UpdateOneID(row.ID).SetName(name).Save(ctx)
	return err
}

func (e *entRepo) List(ctx context.Context) ([]entities.ProjectRole, error) {
	rows, err := e.cli.ProjectRole.Query().Order(ent.Asc(entprojectrole.FieldKey)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.ProjectRole, 0, len(rows))
	for _, r := range rows {
		out = append(out, entities.ProjectRole{ID: r.ID, Key: r.Key, Name: r.Name})
	}
	return out, nil
}

func (e *entRepo) GetByKey(ctx context.Context, key string) (*entities.ProjectRole, error) {
	r, err := e.cli.ProjectRole.Query().Where(entprojectrole.KeyEQ(key)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return &entities.ProjectRole{ID: r.ID, Key: r.Key, Name: r.Name}, nil
}
