package projectmembership

import (
	"context"
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	ene "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) ProjectIDsForPerson(ctx context.Context, personID uuid.UUID) ([]uuid.UUID, error) {
	evts, err := r.cli.Event.Query().
		Where(ene.TypeEQ(events.ProjectRoleAssignedType)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	unique := make(map[uuid.UUID]struct{})
	for _, e := range evts {
		b, _ := json.Marshal(e.Data)
		var payload events.ProjectRoleAssigned
		if err := json.Unmarshal(b, &payload); err == nil && payload.PersonID == personID {
			unique[e.ProjectID] = struct{}{}
		}
	}

	out := make([]uuid.UUID, 0, len(unique))
	for id := range unique {
		out = append(out, id)
	}
	return out, nil
}

func (r *EntRepo) PersonIDsForProjects(ctx context.Context, projectIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(projectIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	evts, err := r.cli.Event.Query().
		Where(
			ene.TypeEQ(events.ProjectRoleAssignedType),
			ene.ProjectIDIn(projectIDs...),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	unique := make(map[uuid.UUID]struct{})
	for _, e := range evts {
		b, _ := json.Marshal(e.Data)
		var payload events.ProjectRoleAssigned
		if err := json.Unmarshal(b, &payload); err == nil {
			unique[payload.PersonID] = struct{}{}
		}
	}

	out := make([]uuid.UUID, 0, len(unique))
	for id := range unique {
		out = append(out, id)
	}
	return out, nil
}
