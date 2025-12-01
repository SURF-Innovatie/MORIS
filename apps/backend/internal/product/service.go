package product

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/person"
	"github.com/SURF-Innovatie/MORIS/ent/product"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetAll(context.Context) ([]*entities.Product, error)
	GetAllForUser(ctx context.Context, personId uuid.UUID) ([]*entities.Product, error)
	Create(ctx context.Context, product entities.Product) (*entities.Product, error)
	Update(ctx context.Context, id uuid.UUID, product entities.Product) (*entities.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	cli *ent.Client
}

func NewService(cli *ent.Client) Service {
	return &service{cli: cli}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	row, err := s.cli.Product.
		Query().
		Where(product.IDEQ(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (s *service) GetAll(ctx context.Context) ([]*entities.Product, error) {
	rows, err := s.cli.Product.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	products := make([]*entities.Product, 0)
	for _, row := range rows {
		products = append(products, mapRow(row))
	}
	return products, nil
}

func (s *service) Create(ctx context.Context, product entities.Product) (*entities.Product, error) {
	row, err := s.cli.Product.
		Create().
		SetName(product.Name).
		SetNillableLanguage(&product.Language).
		SetType(int(product.Type)).
		SetNillableDoi(&product.DOI).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, p entities.Product) (*entities.Product, error) {
	row, err := s.cli.Product.
		UpdateOneID(id).
		SetName(p.Name).
		SetType(int(p.Type)).
		SetNillableLanguage(&p.Language).
		SetNillableDoi(&p.DOI).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return mapRow(row), nil
}

func (s *service) GetAllForUser(ctx context.Context, personId uuid.UUID) ([]*entities.Product, error) {
	rows, err := s.cli.Person.Query().
		Where(person.IDEQ(personId)).
		QueryProducts().
		All(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*entities.Product, 0)
	for _, row := range rows {
		products = append(products, mapRow(row))
	}

	return products, nil
}

func mapRow(row *ent.Product) *entities.Product {
	return &entities.Product{
		Id:       row.ID,
		Type:     entities.ProductType(row.Type),
		Language: *row.Language,
		Name:     row.Name,
		DOI:      *row.Doi,
	}
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.cli.Product.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
