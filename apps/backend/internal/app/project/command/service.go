package command

import (
	"context"
	"fmt"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
	"github.com/samber/lo"
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
}

func NewService(
	es eventstore.Store,
	evtSvc event.Service,
	pc cache.ProjectCache,
	ref cache.ProjectCacheRefresher,
	currentUser appauth.CurrentUserProvider,
	entClient EntClientProvider,
	roleSvc projectrole.Service,
) Service {
	s := &service{
		es:          es,
		evtSvc:      evtSvc,
		cache:       pc,
		refresher:   ref,
		currentUser: currentUser,
		entClient:   entClient,
		roleSvc:     roleSvc,
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

		// Check role-based permissions (EBAC)
		if cur != nil {
			userRole := s.getUserProjectRole(ctx, req.ProjectID, u.PersonID())
			if userRole != nil && !userRole.CanUseEventType(req.Type) {
				return nil, fmt.Errorf("your role does not allow executing %s events", req.Type)
			}
		}

		if !meta.IsAllowed(ctx, e, cli) {
			return nil, fmt.Errorf("not allowed to execute %s", req.Type)
		}

		_ = meta.NeedsApproval(ctx, e, cli)

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

func (s *service) onStatusChange(ctx context.Context, e events.Event) error {
	if s.refresher == nil {
		return nil
	}
	_, err := s.refresher.Refresh(ctx, e.AggregateID())
	return err
}
