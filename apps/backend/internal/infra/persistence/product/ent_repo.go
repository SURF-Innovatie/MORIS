package product

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entperson "github.com/SURF-Innovatie/MORIS/ent/person"
	entproduct "github.com/SURF-Innovatie/MORIS/ent/product"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	row, err := r.cli.Product.Query().
		Where(entproduct.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Product](row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]*entities.Product, error) {
	rows, err := r.cli.Product.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[entities.Product](rows), nil
}

func (r *EntRepo) ListByAuthorPersonID(ctx context.Context, personID uuid.UUID) ([]*entities.Product, error) {
	rows, err := r.cli.Person.Query().
		Where(entperson.IDEQ(personID)).
		QueryProducts().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[entities.Product](rows), nil
}

func (r *EntRepo) Create(ctx context.Context, p entities.Product) (*entities.Product, error) {
	row, err := r.cli.Product.Create().
		SetName(p.Name).
		SetType(int(p.Type)).
		SetNillableLanguage(&p.Language).
		SetNillableDoi(&p.DOI).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Product](row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, p entities.Product) (*entities.Product, error) {
	row, err := r.cli.Product.UpdateOneID(id).
		SetName(p.Name).
		SetType(int(p.Type)).
		SetNillableLanguage(&p.Language).
		SetNillableDoi(&p.DOI).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Product](row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.Product.DeleteOneID(id).Exec(ctx)
}
