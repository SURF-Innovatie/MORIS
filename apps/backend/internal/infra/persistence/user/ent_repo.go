package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	row, err := r.cli.User.Query().Where(entuser.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.User{}).FromEnt(row), nil
}

func (r *EntRepo) GetByPersonID(ctx context.Context, personID uuid.UUID) (*entities.User, error) {
	row, err := r.cli.User.Query().Where(entuser.PersonIDEQ(personID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.User{}).FromEnt(row), nil
}

func (r *EntRepo) Create(ctx context.Context, u entities.User) (*entities.User, error) {
	row, err := r.cli.User.Create().
		SetPersonID(u.PersonID).
		SetPassword(u.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.User{}).FromEnt(row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, u entities.User) (*entities.User, error) {
	row, err := r.cli.User.UpdateOneID(id).
		SetPersonID(u.PersonID).
		SetPassword(u.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return (&entities.User{}).FromEnt(row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.User.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return r.cli.User.UpdateOneID(id).SetIsActive(isActive).Exec(ctx)
}

func (r *EntRepo) ListUsers(ctx context.Context, limit, offset int) ([]entities.User, int, error) {
	total, err := r.cli.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.cli.User.Query().Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}
	out := make([]entities.User, 0, len(rows))
	for _, row := range rows {
		out = append(out, *(&entities.User{}).FromEnt(row))
	}
	return out, total, nil
}
