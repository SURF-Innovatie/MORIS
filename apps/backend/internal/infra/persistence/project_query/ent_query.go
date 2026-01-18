package projectquery

import (
	"context"
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	organisationent "github.com/SURF-Innovatie/MORIS/ent/organisationnode"
	personent "github.com/SURF-Innovatie/MORIS/ent/person"
	productent "github.com/SURF-Innovatie/MORIS/ent/product"
	entprojectrole "github.com/SURF-Innovatie/MORIS/ent/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) PeopleByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.Person, error) {
	out := make(map[uuid.UUID]entities.Person)
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

	return lo.Associate(rows, func(p *ent.Person) (uuid.UUID, entities.Person) {
		return p.ID, transform.ToEntity[entities.Person](p)
	}), nil
}

func (r *EntRepo) ProjectRolesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.ProjectRole, error) {
	out := make(map[uuid.UUID]entities.ProjectRole)
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

	return lo.Associate(rows, func(pr *ent.ProjectRole) (uuid.UUID, entities.ProjectRole) {
		return pr.ID, transform.ToEntity[entities.ProjectRole](pr)
	}), nil
}

func (r *EntRepo) ProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]entities.Product, error) {
	if len(ids) == 0 {
		return []entities.Product{}, nil
	}

	rows, err := r.cli.Product.
		Query().
		Where(productent.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return transform.ToEntities[entities.Product](rows), nil
}

func (r *EntRepo) OrganisationNodeByID(ctx context.Context, id uuid.UUID) (entities.OrganisationNode, error) {
	row, err := r.cli.OrganisationNode.
		Query().
		Where(organisationent.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return entities.OrganisationNode{}, err
	}

	return transform.ToEntity[entities.OrganisationNode](row), nil
}

func (r *EntRepo) OrganisationNodesByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]entities.OrganisationNode, error) {
	out := make(map[uuid.UUID]entities.OrganisationNode)
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

	return lo.Associate(rows, func(o *ent.OrganisationNode) (uuid.UUID, entities.OrganisationNode) {
		return o.ID, transform.ToEntity[entities.OrganisationNode](o)
	}), nil
}

func (r *EntRepo) ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error) {
	evts, err := r.cli.Event.
		Query().
		Where(en.TypeEQ(events.ProjectRoleAssignedType)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	projectIDs := lo.FilterMap(evts, func(e *ent.Event, _ int) (uuid.UUID, bool) {
		b, _ := json.Marshal(e.Data)
		var payload events.ProjectRoleAssigned
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
		Where(en.TypeEQ(events.ProjectStartedType)).
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
