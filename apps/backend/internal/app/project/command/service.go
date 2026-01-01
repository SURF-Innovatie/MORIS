package command

import (
	"context"
	"fmt"

	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
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
}

func NewService(
	es eventstore.Store,
	evtSvc event.Service,
	pc cache.ProjectCache,
	ref cache.ProjectCacheRefresher,
	currentUser appauth.CurrentUserProvider,
	entClient EntClientProvider,
) Service {
	s := &service{
		es:          es,
		evtSvc:      evtSvc,
		cache:       pc,
		refresher:   ref,
		currentUser: currentUser,
		entClient:   entClient,
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
	_, err := s.currentUser.Current(ctx)
	if err != nil {
		return nil, err
	}

	metas := events.GetAllMetas()
	out := make([]AvailableEvent, 0, len(metas))

	for _, m := range metas {
		out = append(out, AvailableEvent{
			Type:          m.Type,
			FriendlyName:  m.FriendlyName,
			NeedsApproval: false,
			Allowed:       true,
			InputSchema:   events.GetInputSchema(m.Type),
		})
	}

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

func (s *service) onStatusChange(ctx context.Context, e events.Event) error {
	if s.refresher == nil {
		return nil
	}
	_, err := s.refresher.Refresh(ctx, e.AggregateID())
	return err
}
