package events

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	events2 "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project/projection"
	eventrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// PolicyExecutionHandler executes event policies for occurred events.
type PolicyExecutionHandler struct {
	evaluator eventpolicy.Evaluator
	eventRepo *eventrepo.EntRepo
}

func NewPolicyExecutionHandler(evaluator eventpolicy.Evaluator, eventRepo *eventrepo.EntRepo) *PolicyExecutionHandler {
	return &PolicyExecutionHandler{evaluator: evaluator, eventRepo: eventRepo}
}

func (h *PolicyExecutionHandler) Handle(ctx context.Context, evt events2.Event) error {
	// If you want this rule, it belongs here (or inside evaluator). Keeping it here is fine.
	if evt.GetStatus() == events2.StatusRejected {
		return nil
	}

	projectID := evt.AggregateID()
	if projectID == uuid.Nil {
		return nil
	}

	project, err := h.loadProject(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return nil
	}

	return h.evaluator.EvaluateAndExecute(ctx, evt, project)
}

func (h *PolicyExecutionHandler) loadProject(ctx context.Context, projectID uuid.UUID) (*project.Project, error) {
	evts, version, err := h.eventRepo.Load(ctx, projectID)
	if err != nil {
		log.Error().Err(err).Msgf("PolicyExecutionHandler: failed to load events for project %s", projectID)
		return nil, err
	}
	if len(evts) == 0 {
		return nil, nil
	}

	p := projection.Reduce(projectID, evts)
	if p == nil {
		return nil, nil
	}
	p.Version = version
	return p, nil
}
