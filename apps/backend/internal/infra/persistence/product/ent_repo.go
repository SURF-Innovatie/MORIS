package product

import (
	"context"
	"strings"

	"github.com/SURF-Innovatie/MORIS/ent"
	entperson "github.com/SURF-Innovatie/MORIS/ent/person"
	entproduct "github.com/SURF-Innovatie/MORIS/ent/product"
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
		WithAuthors().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return (&product.Product{}).FromEnt(row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]*product.Product, error) {
	rows, err := r.cli.Product.Query().
		WithAuthors().
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*product.Product, 0, len(rows))
	for _, row := range rows {
		out = append(out, (&product.Product{}).FromEnt(row))
	}
	return out, nil
}

func (r *EntRepo) ListByAuthorPersonID(ctx context.Context, personID uuid.UUID) ([]*product.Product, error) {
	rows, err := r.cli.Person.Query().
		Where(entperson.IDEQ(personID)).
		QueryAuthoredProducts().
		WithAuthors().
		All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*product.Product, 0, len(rows))
	for _, row := range rows {
		out = append(out, (&product.Product{}).FromEnt(row))
	}
	return out, nil
}

func (r *EntRepo) Create(ctx context.Context, p product.Product) (*product.Product, error) {
	lang := strings.TrimSpace(p.Language)
	doi := strings.TrimSpace(p.DOI)

	builder := r.cli.Product.Create().
		SetName(p.Name).
		SetType(int(p.Type))

	if lang != "" {
		builder = builder.SetLanguage(lang)
	}
	if doi != "" {
		builder = builder.SetDoi(doi)
	}
	if p.ZenodoDepositionID != 0 {
		builder = builder.SetZenodoDepositionID(p.ZenodoDepositionID)
	}

	if len(p.AuthorPersonIDs) > 0 {
		builder = builder.AddAuthorIDs(p.AuthorPersonIDs...)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	row, err = r.cli.Product.Query().
		Where(entproduct.IDEQ(row.ID)).
		WithAuthors().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return (&product.Product{}).FromEnt(row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, p product.Product) (*product.Product, error) {
	lang := strings.TrimSpace(p.Language)
	doi := strings.TrimSpace(p.DOI)

	upd := r.cli.Product.UpdateOneID(id).
		SetName(p.Name).
		SetType(int(p.Type)).
		SetNillableLanguage(nilIfEmpty(lang)).
		SetNillableDoi(nilIfEmpty(doi)).
		SetNillableZenodoDepositionID(nilIfZero(p.ZenodoDepositionID))

	// Replace authors only if field is provided (nil = leave unchanged)
	if p.AuthorPersonIDs != nil {
		upd = upd.ClearAuthors()
		if len(p.AuthorPersonIDs) > 0 {
			upd = upd.AddAuthorIDs(p.AuthorPersonIDs...)
		}
	}

	if _, err := upd.Save(ctx); err != nil {
		return nil, err
	}

	row, err := r.cli.Product.Query().
		Where(entproduct.IDEQ(id)).
		WithAuthors().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return (&product.Product{}).FromEnt(row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.Product.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) GetByDOI(ctx context.Context, doi string) (*product.Product, error) {
	doi = strings.TrimSpace(doi)
	if doi == "" {
		return nil, nil
	}

	row, err := r.cli.Product.Query().
		Where(entproduct.DoiEQ(doi)).
		WithAuthors().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return (&product.Product{}).FromEnt(row), nil
}

func nilIfEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	v := strings.TrimSpace(s)
	return &v
}

func nilIfZero(i int) *int {
	if i == 0 {
		return nil
	}
	v := i
	return &v
}
