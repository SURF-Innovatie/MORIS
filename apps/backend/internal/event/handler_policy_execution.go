package event

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PolicyExecutionHandler executes event policies for occurred events
type PolicyExecutionHandler struct {
	Evaluator  eventpolicy.Evaluator
	ProjectSvc queries.Service
}

func (h *PolicyExecutionHandler) Handle(ctx context.Context, event events.Event) error {
	// Skip policy evaluation for rejected events
	// Approved events are allowed through so notification policies can trigger
	status := event.GetStatus()
	if status == "rejected" {
		return nil
	}

	// Policies apply to projects. We need the project context.
	// Currently policies are linked to ProjectID (aggregate ID).
	projectID := event.AggregateID()
	if projectID == uuid.Nil {
		return nil
	}

	// Fetch project to context (Evaluator needs it for conditions)
	// Uses GetProject to retrieve the full project aggregate state
	details, err := h.ProjectSvc.GetProject(ctx, projectID)
	if err != nil {
		if err == queries.ErrNotFound {
			return nil
		}
		logrus.Errorf("PolicyExecutionHandler: failed to get project %s: %v", projectID, err)
		return err
	}
	if details == nil {
		return nil
	}

	return h.Evaluator.EvaluateAndExecute(ctx, event, &details.Project)
}
