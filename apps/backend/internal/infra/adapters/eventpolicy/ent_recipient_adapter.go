package eventpolicy

import (
	"context"
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/ent/membership"
	"github.com/SURF-Innovatie/MORIS/ent/rolescope"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type RecipientAdapter struct{ client *ent.Client }

func NewRecipientAdapter(client *ent.Client) *RecipientAdapter {
	return &RecipientAdapter{client: client}
}

func (r *RecipientAdapter) ResolveUsers(ctx context.Context, personIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(personIDs) == 0 {
		return []uuid.UUID{}, nil
	}
	var userIDs []uuid.UUID
	if err := r.client.User.Query().
		Where(entuser.PersonIDIn(personIDs...)).
		Select(entuser.FieldID).
		Scan(ctx, &userIDs); err != nil {
		return nil, err
	}
	return lo.Uniq(userIDs), nil
}

func (r *RecipientAdapter) ResolveRole(ctx context.Context, roleID uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	personIDs, err := r.projectMembersWithRole(ctx, roleID, projectID)
	if err != nil {
		return nil, err
	}
	return r.ResolveUsers(ctx, personIDs)
}

func (r *RecipientAdapter) ResolveOrgRole(ctx context.Context, roleID uuid.UUID, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	// NOTE: currently ignores orgNodeID. Either:
	// (A) clarify semantics and filter by scopes rooted in orgNodeID+ancestors, or
	// (B) remove orgNodeID param from the interface if it’s not needed.
	memberships, err := r.client.Membership.Query().
		Where(membership.HasRoleScopeWith(rolescope.RoleIDEQ(roleID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	personIDs := lo.FilterMap(memberships, func(m *ent.Membership, _ int) (uuid.UUID, bool) {
		return m.PersonID, m.PersonID != uuid.Nil
	})
	return r.ResolveUsers(ctx, lo.Uniq(personIDs))
}

func (r *RecipientAdapter) ResolveDynamic(ctx context.Context, dynType string, projectID uuid.UUID, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	switch dynType {
	case "project_members":
		personIDs, err := r.projectMembers(ctx, projectID)
		if err != nil {
			log.Error().Err(err).Msg("ResolveDynamic project_members error")
			return nil, err
		}
		return r.ResolveUsers(ctx, personIDs)

	case "project_owner":
		createdEvt, err := r.client.Event.Query().
			Where(en.TypeEQ(events2.ProjectStartedType), en.ProjectIDEQ(projectID)).
			First(ctx)
		if err != nil {
			log.Error().Err(err).Msg("ResolveDynamic project_owner error")
			return nil, nil
		}
		return []uuid.UUID{createdEvt.CreatedBy}, nil

	case "org_admins":
		return []uuid.UUID{}, nil // TODO
	default:
		return []uuid.UUID{}, nil
	}
}

func (r *RecipientAdapter) projectMembers(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error) {
	evts, err := r.client.Event.Query().
		Where(en.TypeEQ(events2.ProjectRoleAssignedType), en.ProjectIDEQ(projectID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	personIDs := lo.FilterMap(evts, func(e *ent.Event, _ int) (uuid.UUID, bool) {
		var payload events2.ProjectRoleAssigned
		b, _ := json.Marshal(e.Data)
		if json.Unmarshal(b, &payload) == nil && payload.PersonID != uuid.Nil {
			return payload.PersonID, true
		}
		return uuid.Nil, false
	})
	return lo.Uniq(personIDs), nil
}

func (r *RecipientAdapter) projectMembersWithRole(ctx context.Context, roleID, projectID uuid.UUID) ([]uuid.UUID, error) {
	evts, err := r.client.Event.Query().
		Where(en.TypeEQ(events2.ProjectRoleAssignedType), en.ProjectIDEQ(projectID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	personIDs := lo.FilterMap(evts, func(e *ent.Event, _ int) (uuid.UUID, bool) {
		var payload events2.ProjectRoleAssigned
		b, _ := json.Marshal(e.Data)
		if json.Unmarshal(b, &payload) == nil && payload.ProjectRoleID == roleID && payload.PersonID != uuid.Nil {
			return payload.PersonID, true
		}
		return uuid.Nil, false
	})
	return lo.Uniq(personIDs), nil
}
