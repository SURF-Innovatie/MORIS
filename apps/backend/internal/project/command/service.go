package command

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/commandbus"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/SURF-Innovatie/MORIS/internal/project"
	"github.com/google/uuid"
)

type Service interface {
	ListAvailableEvents(ctx context.Context, projectID *uuid.UUID) ([]dto.AvailableEvent, error)
	ExecuteEvent(ctx context.Context, req dto.ExecuteEventRequest) (*entities.Project, error)
}

type service struct {
	cli       *ent.Client
	es        eventstore.Store
	evtSvc    event.Service
	exec      *commandbus.Executor[entities.Project]
	cache     cache.ProjectCache
	refresher cache.ProjectCacheRefresher
}

func NewService(es eventstore.Store, cli *ent.Client, evtSvc event.Service, pc cache.ProjectCache, ref cache.ProjectCacheRefresher) Service {
	s := &service{
		cli:       cli,
		es:        es,
		evtSvc:    evtSvc,
		cache:     pc,
		refresher: ref,
		exec: commandbus.NewExecutor[entities.Project](
			es,
			evtSvc,
			project.Reducer{},
			project.NewReducer{},
		),
	}

	evtSvc.RegisterStatusChangeHandler(s.onStatusChange)
	return s
}

func (s *service) ListAvailableEvents(ctx context.Context, projectID *uuid.UUID) ([]dto.AvailableEvent, error) {
	_, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
	}

	metas := events.GetAllMetas()
	out := make([]dto.AvailableEvent, 0, len(metas))

	for _, m := range metas {
		allowed := true
		needsApproval := false
		schema := events.GetInputSchema(m.Type)

		out = append(out, dto.AvailableEvent{
			Type:          m.Type,
			FriendlyName:  m.FriendlyName,
			NeedsApproval: needsApproval,
			Allowed:       allowed,
			InputSchema:   schema,
		})
	}

	return out, nil
}

func (s *service) ExecuteEvent(ctx context.Context, req dto.ExecuteEventRequest) (*entities.Project, error) {
	if req.ProjectID == uuid.Nil {
		return nil, fmt.Errorf("projectId is required")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	user, err := currentUser(ctx, s.cli)
	if err != nil {
		return nil, err
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

		e, err := decider(ctx, req.ProjectID, user.ID, cur, req.Input, status)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, nil
		}

		if !meta.IsAllowed(ctx, e, s.cli) {
			return nil, fmt.Errorf("not allowed to execute %s", req.Type)
		}

		_ = meta.NeedsApproval(ctx, e, s.cli)

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

func currentUser(ctx context.Context, cli *ent.Client) (*ent.User, error) {
	authUser, ok := httputil.GetUserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no authenticated user in context")
	}
	return cli.User.Get(ctx, authUser.User.ID)
}
