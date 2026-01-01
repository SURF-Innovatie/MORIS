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
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
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

	for _, p := range rows {
		out[p.ID] = entities.Person{
			ID:          p.ID,
			UserID:      p.UserID,
			Name:        p.Name,
			GivenName:   p.GivenName,
			FamilyName:  p.FamilyName,
			Email:       p.Email,
			ORCiD:       &p.OrcidID,
			AvatarUrl:   p.AvatarURL,
			Description: p.Description,
		}
	}
	return out, nil
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

	for _, pr := range rows {
		out[pr.ID] = entities.ProjectRole{
			ID:   pr.ID,
			Key:  pr.Key,
			Name: pr.Name,
		}
	}
	return out, nil
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

	out := make([]entities.Product, 0, len(rows))
	for _, p := range rows {
		out = append(out, entities.Product{
			Id:       p.ID,
			Name:     p.Name,
			Language: *p.Language,
			Type:     entities.ProductType(p.Type),
			DOI:      *p.Doi,
		})
	}
	return out, nil
}

func (r *EntRepo) OrganisationNodeByID(ctx context.Context, id uuid.UUID) (entities.OrganisationNode, error) {
	row, err := r.cli.OrganisationNode.
		Query().
		Where(organisationent.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return entities.OrganisationNode{}, err
	}

	return entities.OrganisationNode{
		ID:       row.ID,
		ParentID: row.ParentID,
		Name:     row.Name,
	}, nil
}

func (r *EntRepo) ListProjectRoles(ctx context.Context) ([]entities.ProjectRole, error) {
	rows, err := r.cli.ProjectRole.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]entities.ProjectRole, 0, len(rows))
	for _, pr := range rows {
		out = append(out, entities.ProjectRole{
			ID:   pr.ID,
			Key:  pr.Key,
			Name: pr.Name,
		})
	}
	return out, nil
}

func (r *EntRepo) ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error) {
	evts, err := r.cli.Event.
		Query().
		Where(en.TypeEQ(events.ProjectRoleAssignedType)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	unique := make(map[uuid.UUID]struct{})
	for _, e := range evts {
		b, _ := json.Marshal(e.Data)

		var payload events.ProjectRoleAssigned
		if err := json.Unmarshal(b, &payload); err == nil {
			if payload.PersonID == personID {
				unique[e.ProjectID] = struct{}{}
			}
		}
	}

	out := make([]uuid.UUID, 0, len(unique))
	for id := range unique {
		out = append(out, id)
	}
	return out, nil
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
