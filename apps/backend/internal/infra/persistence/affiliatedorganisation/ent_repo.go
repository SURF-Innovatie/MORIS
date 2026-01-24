package affiliatedorganisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entaffiliatedorganisation "github.com/SURF-Innovatie/MORIS/ent/affiliatedorganisation"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// EntRepo is the ent-based repository for AffiliatedOrganisation.
type EntRepo struct {
	cli *ent.Client
}

// NewEntRepo creates a new ent repository.
func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*entities.AffiliatedOrganisation, error) {
	row, err := r.cli.AffiliatedOrganisation.Query().
		Where(entaffiliatedorganisation.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.AffiliatedOrganisation](row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]*entities.AffiliatedOrganisation, error) {
	rows, err := r.cli.AffiliatedOrganisation.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[entities.AffiliatedOrganisation](rows), nil
}

func (r *EntRepo) Create(ctx context.Context, org entities.AffiliatedOrganisation) (*entities.AffiliatedOrganisation, error) {
	row, err := r.cli.AffiliatedOrganisation.Create().
		SetName(org.Name).
		SetKvkNumber(org.KvkNumber).
		SetRorID(org.RorID).
		SetVatNumber(org.VatNumber).
		SetCity(org.City).
		SetCountry(org.Country).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.AffiliatedOrganisation](row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, org entities.AffiliatedOrganisation) (*entities.AffiliatedOrganisation, error) {
	row, err := r.cli.AffiliatedOrganisation.UpdateOneID(id).
		SetName(org.Name).
		SetKvkNumber(org.KvkNumber).
		SetRorID(org.RorID).
		SetVatNumber(org.VatNumber).
		SetCity(org.City).
		SetCountry(org.Country).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.AffiliatedOrganisation](row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.AffiliatedOrganisation.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.AffiliatedOrganisation, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]entities.AffiliatedOrganisation{}, nil
	}
	rows, err := r.cli.AffiliatedOrganisation.Query().
		Where(entaffiliatedorganisation.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[uuid.UUID]entities.AffiliatedOrganisation, len(rows))
	for _, row := range rows {
		e := transform.ToEntityPtr[entities.AffiliatedOrganisation](row)
		result[e.ID] = *e
	}
	return result, nil
}
