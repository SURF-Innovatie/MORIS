package affiliatedorganisation

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entaffiliatedorganisation "github.com/SURF-Innovatie/MORIS/ent/affiliatedorganisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// EntRepo is the ent-based repository for AffiliatedOrganisation.
type EntRepo struct {
	cli *ent.Client
}

// NewEntRepo creates a new ent repository.
func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*affiliatedorganisation.AffiliatedOrganisation, error) {
	row, err := r.cli.AffiliatedOrganisation.Query().
		Where(entaffiliatedorganisation.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	// Manual transform since generic might be tricky with new package if not exact match or if I just want to be safe
	// or assumes transform works if structure matches.
	// Actually, transform.ToEntityPtr[T] relies on T having FromEnt method.
	// So transform.ToEntityPtr[affiliatedorganisation.AffiliatedOrganisation](row) should work.
	var entity affiliatedorganisation.AffiliatedOrganisation
	return entity.FromEnt(row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]*affiliatedorganisation.AffiliatedOrganisation, error) {
	rows, err := r.cli.AffiliatedOrganisation.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(rows, func(row *ent.AffiliatedOrganisation, _ int) *affiliatedorganisation.AffiliatedOrganisation {
		var entity affiliatedorganisation.AffiliatedOrganisation
		return entity.FromEnt(row)
	}), nil
}

func (r *EntRepo) Create(ctx context.Context, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error) {
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
	var entity affiliatedorganisation.AffiliatedOrganisation
	return entity.FromEnt(row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, org affiliatedorganisation.AffiliatedOrganisation) (*affiliatedorganisation.AffiliatedOrganisation, error) {
	u := r.cli.AffiliatedOrganisation.UpdateOneID(id).
		SetName(org.Name).
		SetKvkNumber(org.KvkNumber).
		SetRorID(org.RorID).
		SetVatNumber(org.VatNumber).
		SetCity(org.City).
		SetCountry(org.Country)

	row, err := u.Save(ctx)
	if err != nil {
		return nil, err
	}
	var entity affiliatedorganisation.AffiliatedOrganisation
	return entity.FromEnt(row), nil
}

func (r *EntRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.cli.AffiliatedOrganisation.DeleteOneID(id).Exec(ctx)
}

func (r *EntRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation{}, nil
	}
	rows, err := r.cli.AffiliatedOrganisation.Query().
		Where(entaffiliatedorganisation.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation, len(rows))
	for _, row := range rows {
		var entity affiliatedorganisation.AffiliatedOrganisation
		e := entity.FromEnt(row)
		result[e.ID] = *e
	}
	return result, nil
}
