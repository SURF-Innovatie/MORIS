package command

import (
	"context"
	"encoding/json"
	"fmt"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	rbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	orgrole "github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Service interface {
	ListAvailableEvents(ctx context.Context, projectID *uuid.UUID) ([]AvailableEvent, error)
	ExecuteEvent(ctx context.Context, req ExecuteEventRequest) (*entities.Project, error)
}

type service struct {
	es          eventstore.Store
	evtSvc      event.Service
	exec        *commandbus.Executor[entities.Project]
	cache       cache.ProjectCache
	refresher   cache.ProjectCacheRefresher
	currentUser appauth.CurrentUserProvider
	entClient   EntClientProvider
	roleSvc     projectrole.Service
	evaluator   eventpolicy.Evaluator
	orgSvc      organisation.Service
	rbacSvc     rbacsvc.Service
}

func NewService(
	es eventstore.Store,
	evtSvc event.Service,
	pc cache.ProjectCache,
	ref cache.ProjectCacheRefresher,
	currentUser appauth.CurrentUserProvider,
	entClient EntClientProvider,
	roleSvc projectrole.Service,
	evaluator eventpolicy.Evaluator,
	orgSvc organisation.Service,
	rbacSvc rbacsvc.Service,
) Service {
	s := &service{
		es:          es,
		evtSvc:      evtSvc,
		cache:       pc,
		refresher:   ref,
		currentUser: currentUser,
		entClient:   entClient,
		roleSvc:     roleSvc,
		evaluator:   evaluator,
		orgSvc:      orgSvc,
		rbacSvc:     rbacSvc,
		exec: commandbus.NewExecutor[entities.Project](
			es,
			evtSvc,
			Reducer{},
			NewReducer{},
		),
	}

	evtSvc.RegisterStatusChangeHandler(s.onStatusChange)
	return s
}

func (s *service) ListAvailableEvents(ctx context.Context, projectID *uuid.UUID) ([]AvailableEvent, error) {
	u, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}

	metas := events.GetAllMetas()

	// If projectID provided, get user's role and filter based on allowed events
	var userRole *entities.ProjectRole
	if projectID != nil && *projectID != uuid.Nil {
		userRole = s.getUserProjectRole(ctx, *projectID, u.PersonID())
	}

	out := lo.Map(metas, func(m events.EventMeta, _ int) AvailableEvent {
		allowed := true
		if userRole != nil {
			allowed = userRole.CanUseEventType(m.Type)
		}
		return AvailableEvent{
			Type:          m.Type,
			FriendlyName:  m.FriendlyName,
			NeedsApproval: false,
			Allowed:       allowed,
			InputSchema:   events.GetInputSchema(m.Type),
		}
	})

	return out, nil
}

func (s *service) ExecuteEvent(ctx context.Context, req ExecuteEventRequest) (*entities.Project, error) {
	if req.ProjectID == uuid.Nil {
		return nil, fmt.Errorf("projectId is required")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	u, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}
	cli := s.entClient.Client()

	if req.Type == events.ProjectStartedType {
		// Parse input safely using json encoding
		var inputMap map[string]interface{}
		if err := json.Unmarshal(req.Input, &inputMap); err == nil {
			if strID, ok := inputMap["owning_org_node_id"].(string); ok {
				orgID, err := uuid.Parse(strID)
				if err == nil {
					has, err := s.rbacSvc.HasPermission(ctx, u.PersonID(), orgID, orgrole.PermissionCreateProject)
					if err != nil {
						return nil, err
					}
					if !has {
						return nil, fmt.Errorf("you do not have permission to create projects in this organisation")
					}
				}
			}
		}
	}

	decider, ok := events.GetDecider(req.Type)
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", req.Type)
	}

	status := req.Status
	if status == "" {
		status = events.StatusApproved
	}

	proj, err := s.exec.Execute(ctx, req.ProjectID, func(ctx context.Context, cur *entities.Project) ([]events.Event, error) {
		meta := events.GetMeta(req.Type)

		e, err := decider(ctx, req.ProjectID, u.UserID(), cur, req.Input, status)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, nil
		}

		// Auto-role assignment for ProjectStarted
		if e.Type() == events.ProjectStartedType {
			if started, ok := e.(*events.ProjectStarted); ok {
				// Find a role that allows all events
				role, err := s.findPermissiveRole(ctx, started.OwningOrgNodeID)
				if err != nil {
					return nil, fmt.Errorf("failed to assign initial role: %w", err)
				}

				// Manually construct ProjectRoleAssigned event since we are in genesis block
				// and don't have a valid 'cur' project state for the standard decider.
				assignEvt := &events.ProjectRoleAssigned{
					Base:          events.NewBase(started.ProjectID, u.UserID(), events.StatusApproved),
					PersonID:      u.PersonID(),
					ProjectRoleID: role.ID,
				}

				return []events.Event{e, assignEvt}, nil
			}
		}

		// Check role-based permissions (EBAC)
		if cur != nil {
			userRole := s.getUserProjectRole(ctx, req.ProjectID, u.PersonID())
			if userRole != nil && !userRole.CanUseEventType(req.Type) && !u.IsSysAdmin() {
				return nil, fmt.Errorf("your role does not allow executing %s events", req.Type)
			}
		}

		if !meta.IsAllowed(ctx, e, cli) {
			return nil, fmt.Errorf("not allowed to execute %s", req.Type)
		}

		// Check if any policy requires approval
		needsApproval, err := s.evaluator.CheckApprovalRequired(ctx, e, cur)
		if err != nil {
			logrus.Infof("error checking approval policy: %v", err)
			// Decide if error should block or assume needed/not needed.
			// Safe default: log and proceed (unless policy evaluation is strict requirement).
		}

		if needsApproval {
			// Update the event status to Pending
			base := events.Base{
				ID:        e.GetID(),
				ProjectID: e.AggregateID(),
				At:        e.OccurredAt(),
				CreatedBy: e.CreatedByID(),
				Status:    events.StatusPending,
			}
			e.SetBase(base)
		} else {
			// Check legacy approval mechanism if no policy triggered it
			if meta.NeedsApproval(ctx, e, cli) {
				base := events.Base{
					ID:        e.GetID(),
					ProjectID: e.AggregateID(),
					At:        e.OccurredAt(),
					CreatedBy: e.CreatedByID(),
					Status:    events.StatusPending,
				}
				e.SetBase(base)
			}
		}

		return []events.Event{e}, nil
	})

	if err != nil {
		return nil, err
	}

	_ = s.cache.SetProject(ctx, proj)

	return proj, nil
}

// getUserProjectRole returns the user's role on the project, if any
func (s *service) getUserProjectRole(ctx context.Context, projectID, personID uuid.UUID) *entities.ProjectRole {
	// Get project from cache to find user's role
	proj, err := s.cache.GetProject(ctx, projectID)
	if err != nil || proj == nil {
		return nil
	}

	// Find user's membership and role
	m, ok := lo.Find(proj.Members, func(m entities.ProjectMember) bool {
		return m.PersonID == personID
	})
	if ok {
		// Get the role details
		role, err := s.roleSvc.GetByID(ctx, m.ProjectRoleID)
		if err == nil {
			return role
		}
	}
	return nil
}

func (s *service) findPermissiveRole(ctx context.Context, orgID uuid.UUID) (*entities.ProjectRole, error) {
	roles, err := s.roleSvc.ListAvailableForNode(ctx, orgID)
	if err != nil {
		return nil, err
	}


	// Strategy: Check if role has all known event types.
	allEvents := events.GetRegisteredEventTypes()

	for _, r := range roles {

		// Let's check if it has all events.
		missing := lo.Filter(allEvents, func(t string, _ int) bool {
			return !r.CanUseEventType(t)
		})

		if len(missing) == 0 {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no permissive role found for organisation %s", orgID)
}

func (s *service) onStatusChange(ctx context.Context, e events.Event) error {
	if s.refresher == nil {
		return nil
	}
	_, err := s.refresher.Refresh(ctx, e.AggregateID())
	return err
}
