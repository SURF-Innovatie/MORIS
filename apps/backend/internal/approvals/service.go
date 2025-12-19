package approvals

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/approvalrequest"
	"github.com/SURF-Innovatie/MORIS/internal/domain/approvals"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type Options struct {
	StrictMissingRequest bool
}

type Service interface {
	FilterEffective(ctx context.Context, projectID uuid.UUID, evts []events.Event) ([]events.Event, error)
	Pending(ctx context.Context, projectID uuid.UUID, evts []events.Event) ([]events.Event, error)
}

type service struct {
	cli    *ent.Client
	policy approvals.Policy
	opts   Options
}

func NewService(cli *ent.Client, policy approvals.Policy, opts Options) Service {
	return &service{cli: cli, policy: policy, opts: opts}
}

func (s *service) FilterEffective(ctx context.Context, projectID uuid.UUID, evts []events.Event) ([]events.Event, error) {
	if len(evts) == 0 {
		return evts, nil
	}

	needsApprovalIDs := make([]uuid.UUID, 0)
	needsApprovalSet := make(map[uuid.UUID]struct{})

	for _, e := range evts {
		if _, ok := s.policy.ForEventType(e.Type()); ok {
			id := e.GetID()
			needsApprovalIDs = append(needsApprovalIDs, id)
			needsApprovalSet[id] = struct{}{}
		}
	}

	if len(needsApprovalIDs) == 0 {
		return evts, nil
	}

	type row struct {
		EventID uuid.UUID `json:"event_id"`
		Status  string    `json:"status"`
	}
	rows := []row{}

	if err := s.cli.ApprovalRequest.
		Query().
		Where(
			approvalrequest.ProjectIDEQ(projectID),
			approvalrequest.EventIDIn(needsApprovalIDs...),
		).
		Select(approvalrequest.FieldEventID, approvalrequest.FieldStatus).
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	approved := make(map[uuid.UUID]struct{}, len(rows))
	for _, r := range rows {
		if r.Status == string(approvalrequest.StatusApproved) {
			approved[r.EventID] = struct{}{}
		}
	}

	out := make([]events.Event, 0, len(evts))
	for _, e := range evts {
		if _, ok := s.policy.ForEventType(e.Type()); !ok {
			out = append(out, e)
			continue
		}

		if _, ok := approved[e.GetID()]; ok {
			out = append(out, e)
			continue
		}

		if s.opts.StrictMissingRequest {
			return nil, fmt.Errorf("event %s (%s) requires approval but has no approved request", e.GetID(), e.Type())
		}
	}

	return out, nil
}

func (s *service) Pending(ctx context.Context, projectID uuid.UUID, evts []events.Event) ([]events.Event, error) {
	if len(evts) == 0 {
		return []events.Event{}, nil
	}

	needsApprovalIDs := make([]uuid.UUID, 0, len(evts))
	for _, e := range evts {
		if _, ok := s.policy.ForEventType(e.Type()); ok {
			needsApprovalIDs = append(needsApprovalIDs, e.GetID())
		}
	}

	if len(needsApprovalIDs) == 0 {
		return []events.Event{}, nil
	}

	var rows []struct {
		EventID uuid.UUID `json:"event_id"`
		Status  string    `json:"status"`
	}

	if err := s.cli.ApprovalRequest.
		Query().
		Where(
			approvalrequest.ProjectIDEQ(projectID),
			approvalrequest.EventIDIn(needsApprovalIDs...),
		).
		Select(approvalrequest.FieldEventID, approvalrequest.FieldStatus).
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	open := make(map[uuid.UUID]struct{}, len(rows))
	for _, r := range rows {
		if r.Status == string(approvalrequest.StatusOpen) {
			open[r.EventID] = struct{}{}
		}
	}

	out := make([]events.Event, 0)
	for _, e := range evts {
		if _, ok := open[e.GetID()]; ok {
			out = append(out, e)
			continue
		}

		if s.opts.StrictMissingRequest {
			if _, ok := s.policy.ForEventType(e.Type()); ok {
				return nil, fmt.Errorf("event %s (%s) requires approval but has no approval request", e.GetID(), e.Type())
			}
		}
	}

	return out, nil
}
