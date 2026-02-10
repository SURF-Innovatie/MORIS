package product

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entperson "github.com/SURF-Innovatie/MORIS/ent/person"
	entproduct "github.com/SURF-Innovatie/MORIS/ent/product"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	row, err := r.cli.Product.Query().
		Where(entproduct.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[product.Product](row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]*product.Product, error) {
	rows, err := r.cli.Product.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[product.Product](rows), nil
}

func (r *EntRepo) ListByAuthorPersonID(ctx context.Context, personID uuid.UUID) ([]*product.Product, error) {
	rows, err := r.cli.Person.Query().
		Where(entperson.IDEQ(personID)).
		QueryProducts().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[product.Product](rows), nil
}

func (r *EntRepo) Create(ctx context.Context, p product.Product) (*product.Product, error) {
	builder := r.cli.Product.Create().
		SetName(p.Name).
		SetType(int(p.Type)).
		SetNillableLanguage(&p.Language).
		SetNillableDoi(&p.DOI).
		SetNillableZenodoDepositionID(&p.ZenodoDepositionID)

	// Link the product to its author if provided
	if p.AuthorPersonID != uuid.Nil {
		builder = builder.AddAuthorIDs(p.AuthorPersonID)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[product.Product](row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error) {
	row, err := r.cli.Product.UpdateOneID(id).
		SetName(p.Name).
		SetType(int(p.Type)).
		SetNillableLanguage(&p.Language).
		SetNillableDoi(&p.DOI).
		SetNillableZenodoDepositionID(&p.ZenodoDepositionID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[product.Product](row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.Product.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) GetByDOI(ctx context.Context, doi string) (*product.Product, error) {
	row, err := r.cli.Product.
		Query().
		Where(entproduct.DoiEQ(doi)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return (&product.Product{}).FromEnt(row), nil
}
