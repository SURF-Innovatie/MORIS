package project

import (
	"context"
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	affiliatedorgent "github.com/SURF-Innovatie/MORIS/ent/affiliatedorganisation"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	organisationent "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	personent "github.com/SURF-Innovatie/MORIS/ent/person"
	productent "github.com/SURF-Innovatie/MORIS/ent/product"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/product"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) PeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error) {
	out := make(map[uuid.UUID]identity.Person)
	if len(ids) == 0 {
		return out, nil
	}

	rows, err := r.cli.Person.
		Query().
		Where(personent.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Associate(rows, func(p *ent.Person) (uuid.UUID, identity.Person) {
		return p.ID, transform.ToEntity[identity.Person](p)
	}), nil
}

func (r *EntRepo) ProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]role.ProjectRole, error) {
	out := make(map[uuid.UUID]role.ProjectRole)
	if len(ids) == 0 {
		return out, nil
	}

	rows, err := r.cli.ProjectRole.
		Query().
		Where(entprojectrole.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Associate(rows, func(pr *ent.ProjectRole) (uuid.UUID, role.ProjectRole) {
		return pr.ID, transform.ToEntity[role.ProjectRole](pr)
	}), nil
}

func (r *EntRepo) ProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]product.Product, error) {
	if len(ids) == 0 {
		return []product.Product{}, nil
	}

	rows, err := r.cli.Product.
		Query().
		Where(productent.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[product.Product](rows), nil
}

func (r *EntRepo) OrganisationNodeByID(ctx context.Context, id uuid.UUID) (organisation.OrganisationNode, error) {
	row, err := r.cli.OrganisationNode.
		Query().
		Where(organisationent.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return organisation.OrganisationNode{}, err
	}

	return transform.ToEntity[organisation.OrganisationNode](row), nil
}

func (r *EntRepo) OrganisationNodesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]organisation.OrganisationNode, error) {
	out := make(map[uuid.UUID]organisation.OrganisationNode)
	if len(ids) == 0 {
		return out, nil
	}

	rows, err := r.cli.OrganisationNode.
		Query().
		Where(organisationent.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Associate(rows, func(o *ent.OrganisationNode) (uuid.UUID, organisation.OrganisationNode) {
		return o.ID, transform.ToEntity[organisation.OrganisationNode](o)
	}), nil
}

// GetPeopleByIDs is an alias for PeopleByIDs to satisfy the hydrator.PersonLoader interface
func (r *EntRepo) GetPeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]identity.Person, error) {
	return r.PeopleByIDs(ctx, ids)
}

func (r *EntRepo) ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error) {
	evts, err := r.cli.Event.
		Query().
		Where(en.TypeEQ(events2.ProjectRoleAssignedType)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	projectIDs := lo.FilterMap(evts, func(e *ent.Event, _ int) (uuid.UUID, bool) {
		b, _ := json.Marshal(e.Data)
		var payload events2.ProjectRoleAssigned
		if err := json.Unmarshal(b, &payload); err == nil {
			if payload.PersonID == personID {
				return e.ProjectID, true
			}
		}
		return uuid.Nil, false
	})

	return lo.Uniq(projectIDs), nil
}

func (r *EntRepo) ProjectIDsStarted(ctx context.Context) ([]uuid.UUID, error) {
	var projectIDs []uuid.UUID
	if err := r.cli.Event.Query().
		Where(en.TypeEQ(events2.ProjectStartedType)).
		Select(en.FieldProjectID).
		Scan(ctx, &projectIDs); err != nil {
		return nil, err
	}
	return projectIDs, nil
}

func (r *EntRepo) ListAncestors(ctx context.Context, orgID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.cli.OrganisationNodeClosure.Query().
		Where().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Find all ancestors of the given org node
	return lo.FilterMap(rows, func(row *ent.OrganisationNodeClosure, _ int) (uuid.UUID, bool) {
		if row.DescendantID == orgID {
			return row.AncestorID, true
		}
		return uuid.Nil, false
	}), nil
}

func (r *EntRepo) GetAffiliatedOrganisationsByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation, error) {
	out := make(map[uuid.UUID]affiliatedorganisation.AffiliatedOrganisation)
	if len(ids) == 0 {
		return out, nil
	}

	rows, err := r.cli.AffiliatedOrganisation.
		Query().
		Where(affiliatedorgent.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Associate(rows, func(a *ent.AffiliatedOrganisation) (uuid.UUID, affiliatedorganisation.AffiliatedOrganisation) {
		var entity affiliatedorganisation.AffiliatedOrganisation
		return a.ID, *entity.FromEnt(a)
	}), nil
}
