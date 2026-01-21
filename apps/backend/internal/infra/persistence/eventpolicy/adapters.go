package eventpolicy

import (
	"context"
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	en "github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/ent/membership"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	"github.com/SURF-Innovatie/MORIS/ent/rolescope"
	entuser "github.com/SURF-Innovatie/MORIS/ent/user"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// OrgClosureAdapter implements eventpolicy.OrgClosureProvider using the org repository
type OrgClosureAdapter struct {
	orgRepo organisation.Repository
}

// NewOrgClosureAdapter creates a new OrgClosureAdapter
func NewOrgClosureAdapter(orgRepo organisation.Repository) *OrgClosureAdapter {
	return &OrgClosureAdapter{orgRepo: orgRepo}
}

// GetAncestorIDs returns all ancestor org node IDs for a given node (excluding self)
func (a *OrgClosureAdapter) GetAncestorIDs(ctx context.Context, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	closures, err := a.orgRepo.ListClosuresByDescendant(ctx, orgNodeID)
	if err != nil {
		return nil, err
	}

	// Filter out self-reference (depth 0) and return ancestor IDs
	return lo.FilterMap(closures, func(c entities.OrganisationNodeClosure, _ int) (uuid.UUID, bool) {
		if c.Depth == 0 {
			return uuid.Nil, false // Skip self
		}
		return c.AncestorID, true
	}), nil
}

// NotificationAdapter implements eventpolicy.NotificationSender
type NotificationAdapter struct {
	client *ent.Client
}

// NewNotificationAdapter creates a new NotificationAdapter
func NewNotificationAdapter(client *ent.Client) *NotificationAdapter {
	return &NotificationAdapter{client: client}
}

// SendNotification creates info notifications for users about an event
func (n *NotificationAdapter) SendNotification(ctx context.Context, userIDs []uuid.UUID, eventID uuid.UUID, message string) error {
	for _, userID := range userIDs {
		_, err := n.client.Notification.Create().
			SetUserID(userID).
			SetEventID(eventID).
			SetMessage(message).
			SetType(notification.TypeInfo).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendApprovalRequest creates approval request notifications for users
func (n *NotificationAdapter) SendApprovalRequest(ctx context.Context, userIDs []uuid.UUID, eventID uuid.UUID, message string) error {
	for _, userID := range userIDs {
		_, err := n.client.Notification.Create().
			SetUserID(userID).
			SetEventID(eventID).
			SetMessage(message).
			SetType(notification.TypeApprovalRequest).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// RecipientAdapter implements eventpolicy.RecipientResolver
type RecipientAdapter struct {
	client *ent.Client
}

// NewRecipientAdapter creates a new RecipientAdapter
func NewRecipientAdapter(client *ent.Client) *RecipientAdapter {
	return &RecipientAdapter{client: client}
}

// ResolveUsers converts person IDs to user IDs (since policies store person IDs as "user IDs")
func (r *RecipientAdapter) ResolveUsers(ctx context.Context, personIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(personIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	var userIDs []uuid.UUID
	err := r.client.User.Query().
		Where(entuser.PersonIDIn(personIDs...)).
		Select(entuser.FieldID).
		Scan(ctx, &userIDs)
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

// ResolveRole returns user IDs for users with a given role
func (r *RecipientAdapter) ResolveRole(ctx context.Context, roleID uuid.UUID, projectID uuid.UUID) ([]uuid.UUID, error) {
	// Query events to find users who were assigned this role in this project
	evts, err := r.client.Event.Query().
		Where(
			en.TypeEQ(events.ProjectRoleAssignedType),
			en.ProjectIDEQ(projectID),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	personIDs := lo.FilterMap(evts, func(e *ent.Event, _ int) (uuid.UUID, bool) {
		b, _ := json.Marshal(e.Data)
		var payload events.ProjectRoleAssigned
		if err := json.Unmarshal(b, &payload); err == nil {
			// Check against ProjectRoleID
			if payload.ProjectRoleID == roleID {
				return payload.PersonID, true
			}
		}
		return uuid.Nil, false
	})

	personIDs = lo.Uniq(personIDs)
	if len(personIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	// Resolve PersonIDs to UserIDs
	var userIDs []uuid.UUID
	err = r.client.User.Query().
		Where(entuser.PersonIDIn(personIDs...)).
		Select(entuser.FieldID).
		Scan(ctx, &userIDs)
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

// ResolveOrgRole returns user IDs for users with a given organisation role
func (r *RecipientAdapter) ResolveOrgRole(ctx context.Context, roleID uuid.UUID, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	// Query memberships through role_scope to find users with this org role
	// OrganisationRole → RoleScope → Membership → Person → user_id
	memberships, err := r.client.Membership.Query().
		Where(
			membership.HasRoleScopeWith(
				rolescope.RoleIDEQ(roleID),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	personIDs := lo.FilterMap(memberships, func(m *ent.Membership, _ int) (uuid.UUID, bool) {
		if m.PersonID != uuid.Nil {
			return m.PersonID, true
		}
		return uuid.Nil, false
	})

	personIDs = lo.Uniq(personIDs)
	if len(personIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	var userIDs []uuid.UUID
	err = r.client.User.Query().
		Where(entuser.PersonIDIn(personIDs...)).
		Select(entuser.FieldID).
		Scan(ctx, &userIDs)
	if err != nil {
		return nil, err
	}

	return lo.Uniq(userIDs), nil
}

// ResolveDynamic returns user IDs for dynamic recipient types
func (r *RecipientAdapter) ResolveDynamic(ctx context.Context, dynType string, projectID uuid.UUID, orgNodeID uuid.UUID) ([]uuid.UUID, error) {
	switch dynType {
	case "project_members":
		// Get all users who are members of the project via events
		evts, err := r.client.Event.Query().
			Where(
				en.TypeEQ(events.ProjectRoleAssignedType),
				en.ProjectIDEQ(projectID),
			).
			All(ctx)
		if err != nil {
			log.Error().Err(err).Msg("ResolveDynamic project_members error")
			return nil, err
		}

		personIDs := lo.FilterMap(evts, func(e *ent.Event, _ int) (uuid.UUID, bool) {
			b, _ := json.Marshal(e.Data)
			var payload events.ProjectRoleAssigned
			if err := json.Unmarshal(b, &payload); err == nil {
				return payload.PersonID, true
			}
			return uuid.Nil, false
		})

		personIDs = lo.Uniq(personIDs)
		if len(personIDs) == 0 {
			return []uuid.UUID{}, nil
		}

		var userIDs []uuid.UUID
		err = r.client.User.Query().
			Where(entuser.PersonIDIn(personIDs...)).
			Select(entuser.FieldID).
			Scan(ctx, &userIDs)
		if err != nil {
			return nil, err
		}
		return userIDs, nil

	case "project_owner":
		// CreatedBy in ProjectStarted event is the user ID of the creator (actor)
		createdEvt, err := r.client.Event.Query().
			Where(
				en.TypeEQ(events.ProjectStartedType),
				en.ProjectIDEQ(projectID),
			).
			First(ctx)
		if err != nil {
			log.Error().Err(err).Msg("ResolveDynamic project_owner error")
			return nil, nil
		}
		// CreatedBy is the actor UUID (a User ID)
		return []uuid.UUID{createdEvt.CreatedBy}, nil

	case "org_admins":
		// TODO: Cross-module call to Organisation Svc to find admin users of orgNodeID
		return []uuid.UUID{}, nil
	default:
		return []uuid.UUID{}, nil
	}
}
