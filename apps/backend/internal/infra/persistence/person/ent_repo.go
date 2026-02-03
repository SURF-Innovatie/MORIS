package person

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	pe "github.com/SURF-Innovatie/MORIS/ent/person"
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

func (r *EntRepo) Create(ctx context.Context, p identity.Person) (*identity.Person, error) {
	if p.ORCiD != nil && *p.ORCiD == "" {
		p.ORCiD = nil
	}
	row, err := r.cli.Person.
		Create().
		SetName(p.Name).
		SetNillableGivenName(p.GivenName).
		SetNillableFamilyName(p.FamilyName).
		SetEmail(p.Email).
		SetNillableAvatarURL(p.AvatarUrl).
		SetNillableDescription(p.Description).
		SetNillableOrcidID(p.ORCiD).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.Person](row), nil
}

func (r *EntRepo) Get(ctx context.Context, id uuid.UUID) (*identity.Person, error) {
	row, err := r.cli.Person.
		Query().
		Where(pe.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.Person](row), nil
}

func (r *EntRepo) Update(ctx context.Context, id uuid.UUID, p identity.Person) (*identity.Person, error) {
	q := r.cli.Person.
		UpdateOneID(id).
		SetName(p.Name).
		SetNillableGivenName(p.GivenName).
		SetNillableFamilyName(p.FamilyName).
		SetEmail(p.Email).
		SetNillableAvatarURL(p.AvatarUrl).
		SetNillableDescription(p.Description).
		SetOrgCustomFields(p.OrgCustomFields)

	if p.ORCiD != nil && *p.ORCiD == "" {
		q.ClearOrcidID()
	} else {
		q.SetNillableOrcidID(p.ORCiD)
	}

	row, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.Person](row), nil
}

func (r *EntRepo) List(ctx context.Context) ([]*identity.Person, error) {
	rows, err := r.cli.Person.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntitiesPtr[identity.Person](rows), nil
}

func (r *EntRepo) GetByEmail(ctx context.Context, email string) (*identity.Person, error) {
	row, err := r.cli.Person.
		Query().
		Where(pe.EmailEQ(email)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[identity.Person](row), nil
}

func (r *EntRepo) Search(ctx context.Context, query string, limit int) ([]identity.Person, error) {
	rows, err := r.cli.Person.
		Query().
		Where(pe.Or(
			pe.NameContainsFold(query),
			pe.EmailContainsFold(query),
		)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntities[identity.Person](rows), nil
}

func (r *EntRepo) SetORCID(ctx context.Context, personID uuid.UUID, orcidID string) error {
	_, err := r.cli.Person.
		UpdateOneID(personID).
		SetOrcidID(orcidID).
		Save(ctx)
	return err
}

func (r *EntRepo) ClearORCID(ctx context.Context, personID uuid.UUID) error {
	_, err := r.cli.Person.
		UpdateOneID(personID).
		ClearOrcidID().
		Save(ctx)
	return err
}

func (r *EntRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]identity.Person, error) {
	if len(ids) == 0 {
		return []identity.Person{}, nil
	}
	rows, err := r.cli.Person.Query().
		Where(pe.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transform.ToEntities[identity.Person](rows), nil
}
