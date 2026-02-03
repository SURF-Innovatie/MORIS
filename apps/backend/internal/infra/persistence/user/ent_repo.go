package user

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	row, err := r.cli.User.Query().Where(entuser.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.User](row), nil
}

func (r *EntRepo) GetByPersonID(ctx context.Context, personID uuid.UUID) (*identity.User, error) {
	row, err := r.cli.User.Query().Where(entuser.PersonIDEQ(personID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.User](row), nil
}

func (r *EntRepo) Create(ctx context.Context, u identity.User) (*identity.User, error) {
	builder := r.cli.User.Create().
		SetPersonID(u.PersonID).
		SetIsSysAdmin(u.IsSysAdmin)

	// Only set password if provided (OAuth-only users don't have passwords)
	if u.Password != "" {
		builder.SetPassword(u.Password)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.User](row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, u identity.User) (*identity.User, error) {
	row, err := r.cli.User.UpdateOneID(id).
		SetPersonID(u.PersonID).
		SetPassword(u.Password).
		SetIsSysAdmin(u.IsSysAdmin).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.User](row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.User.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return r.cli.User.UpdateOneID(id).SetIsActive(isActive).Exec(ctx)
}

func (r *EntRepo) ListUsers(ctx context.Context, limit, offset int) ([]identity.User, int, error) {
	total, err := r.cli.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.cli.User.Query().Limit(limit).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return transform.ToEntities[identity.User](rows), total, nil
}

func (r *EntRepo) SetZenodoTokens(ctx context.Context, userID uuid.UUID, access, refresh string) error {
	upd := r.cli.User.UpdateOneID(userID).SetZenodoAccessToken(access)
	if refresh != "" {
		upd = upd.SetZenodoRefreshToken(refresh)
	}
	_, err := upd.Save(ctx)
	return err
}

func (r *EntRepo) ClearZenodoTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.cli.User.UpdateOneID(userID).
		ClearZenodoAccessToken().
		ClearZenodoRefreshToken().
		Save(ctx)
	return err
}
