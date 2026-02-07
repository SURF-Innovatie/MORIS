package command

import (
	"context"
	"encoding/json"
	"fmt"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	rbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/role"
	"github.com/SURF-Innovatie/MORIS/internal/domain/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	role2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/role"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type Service interface {
	ListAvailableEvents(ctx context.Context, projectID *uuid.UUID) ([]AvailableEvent, error)
	ExecuteEvent(ctx context.Context, req ExecuteEventRequest) (*project.Project, error)
}

type service struct {
	evtSvc      event.Service
	exec        *commandbus.Executor[project.Project]
	cache       cache.ProjectCache
	currentUser appauth.CurrentUserProvider
	entClient   EntClientProvider
	roleSvc     role.Service
	evaluator   eventpolicy.Evaluator
	orgSvc      organisation.Service
	rbacSvc     rbacsvc.Service
	repo        queries.ProjectReadRepository
}

func NewService(
	evtSvc event.Service,
	pc cache.ProjectCache,
	currentUser appauth.CurrentUserProvider,
	entClient EntClientProvider,
	roleSvc role.Service,
	evaluator eventpolicy.Evaluator,
	orgSvc organisation.Service,
	rbacSvc rbacsvc.Service,
	evtPub event.Publisher,
	repo queries.ProjectReadRepository,
) Service {
	return &service{
		evtSvc:      evtSvc,
		cache:       pc,
		currentUser: currentUser,
		entClient:   entClient,
		roleSvc:     roleSvc,
		evaluator:   evaluator,
		orgSvc:      orgSvc,
		rbacSvc:     rbacSvc,
		repo:        repo,
		exec: commandbus.NewExecutor[project.Project](
			evtSvc,
			evtPub,
			Reducer{},
			NewReducer{},
		),
	}
}

func (s *service) ListAvailableEvents(ctx context.Context, projectID *uuid.UUID) ([]AvailableEvent, error) {
	u, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}

	metas := events2.GetAllMetas()

	// If projectID provided, get user's role and filter based on allowed events
	var userRole *role2.ProjectRole
	if projectID != nil && *projectID != uuid.Nil {
		userRole = s.getUserProjectRole(ctx, *projectID, u.PersonID)
	}

	out := lo.Map(metas, func(m events2.EventMeta, _ int) AvailableEvent {
		allowed := true
		if userRole != nil {
			allowed = userRole.CanUseEventType(m.Type)
		}
		return AvailableEvent{
			Type:          m.Type,
			FriendlyName:  m.FriendlyName,
			NeedsApproval: false,
			Allowed:       allowed,
			InputSchema:   events2.GetInputSchema(m.Type),
		}
	})

	return out, nil
}

func (s *service) ExecuteEvent(ctx context.Context, req ExecuteEventRequest) (*project.Project, error) {
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

	if req.Type == events2.ProjectStartedType {
		// Parse input safely using json encoding
		var inputMap map[string]interface{}
		if err := json.Unmarshal(req.Input, &inputMap); err == nil {
			if strID, ok := inputMap["owning_org_node_id"].(string); ok {
				orgID, err := uuid.Parse(strID)
				if err == nil {
					has, err := s.rbacSvc.HasPermission(ctx, u.PersonID, orgID, rbac.PermissionCreateProject)
					if err != nil {
						return nil, err
					}
					if !has {
						return nil, fmt.Errorf("you do not have permission to create projects in this organisation")
					}
				}
			}

			// Unique slug check
			if slug, ok := inputMap["slug"].(string); ok {
				existingID, err := s.repo.ProjectIDBySlug(ctx, slug)
				if err != nil {
					return nil, err
				}
				if existingID != uuid.Nil {
					return nil, fmt.Errorf("slug '%s' is already in use", slug)
				}
			}
		}
	}

	decider, ok := events2.GetDecider(req.Type)
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", req.Type)
	}

	status := req.Status
	if status == "" {
		status = events2.StatusApproved
	}

	proj, err := s.exec.Execute(ctx, req.ProjectID, func(ctx context.Context, cur *project.Project) ([]events2.Event, error) {
		meta := events2.GetMeta(req.Type)

		e, err := decider(ctx, req.ProjectID, u.UserID, cur, req.Input, status)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, nil
		}

		// Auto-role assignment for ProjectStarted
		if e.Type() == events2.ProjectStartedType {
			if started, ok := e.(*events2.ProjectStarted); ok {
				// Find a role that allows all events
				role, err := s.findPermissiveRole(ctx, started.OwningOrgNodeID)
				if err != nil {
					return nil, fmt.Errorf("failed to assign initial role: %w", err)
				}

				// Manually construct ProjectRoleAssigned event since we are in genesis block
				// and don't have a valid 'cur' project state for the standard decider.
				assignEvt := &events2.ProjectRoleAssigned{
					Base:          events2.NewBase(started.ProjectID, u.UserID, events2.StatusApproved),
					PersonID:      u.PersonID,
					ProjectRoleID: role.ID,
				}

				return []events2.Event{e, assignEvt}, nil
			}
		}

		// Check role-based permissions (EBAC)
		if cur != nil {
			userRole := s.getUserProjectRole(ctx, req.ProjectID, u.PersonID)
			if userRole != nil && !userRole.CanUseEventType(req.Type) && !u.IsSysAdmin {
				return nil, fmt.Errorf("your role does not allow executing %s events", req.Type)
			}
		}

		if !meta.IsAllowed(ctx, e, cli) {
			return nil, fmt.Errorf("not allowed to execute %s", req.Type)
		}

		// Check if any policy requires approval
		needsApproval, err := s.evaluator.CheckApprovalRequired(ctx, e, cur)
		if err != nil {
			log.Error().Err(err).Msg("error checking approval policy")
			// Decide if error should block or assume needed/not needed.
			// Safe default: log and proceed (unless policy evaluation is strict requirement).
		}

		if needsApproval {
			// Update the event status to Pending
			base := events2.Base{
				ID:              e.GetID(),
				ProjectID:       e.AggregateID(),
				At:              e.OccurredAt(),
				CreatedBy:       e.CreatedByID(),
				Status:          events2.StatusPending,
				FriendlyNameStr: e.FriendlyName(),
			}
			e.SetBase(base)
		} else {
			// Check legacy approval mechanism if no policy triggered it
			if meta.NeedsApproval(ctx, e, cli) {
				base := events2.Base{
					ID:              e.GetID(),
					ProjectID:       e.AggregateID(),
					At:              e.OccurredAt(),
					CreatedBy:       e.CreatedByID(),
					Status:          events2.StatusPending,
					FriendlyNameStr: e.FriendlyName(),
				}
				e.SetBase(base)
			}
		}

		return []events2.Event{e}, nil
	})

	if err != nil {
		return nil, err
	}

	_ = s.cache.SetProject(ctx, proj)

	return proj, nil
}

// getUserProjectRole returns the user's role on the project, if any
func (s *service) getUserProjectRole(ctx context.Context, projectID, personID uuid.UUID) *role2.ProjectRole {
	// Get project from cache to find user's role
	proj, err := s.cache.GetProject(ctx, projectID)
	if err != nil || proj == nil {
		return nil
	}

	// Find user's membership and role
	m, ok := lo.Find(proj.Members, func(m project.Member) bool {
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

func (s *service) findPermissiveRole(ctx context.Context, orgID uuid.UUID) (*role2.ProjectRole, error) {
	roles, err := s.roleSvc.ListAvailableForNode(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Strategy: Check if role has all known event types.
	allEvents := events2.GetRegisteredEventTypes()

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
