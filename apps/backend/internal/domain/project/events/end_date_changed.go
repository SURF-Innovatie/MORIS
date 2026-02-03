package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const EndDateChangedType = "project.end_date_changed"

type EndDateChanged struct {
	Base
	EndDate time.Time `json:"endDate"`
}

func (EndDateChanged) isEvent()     {}
func (EndDateChanged) Type() string { return EndDateChangedType }
func (e EndDateChanged) String() string {
	return fmt.Sprintf("End date changed: %s", e.EndDate.Format("2006-01-02"))
}

func (e *EndDateChanged) Apply(project *project.Project) {
	project.EndDate = e.EndDate
}

func (e *EndDateChanged) NotificationMessage() string {
	return fmt.Sprintf("Project end date changed to %s.", e.EndDate.Format("2006-01-02"))
}

type EndDateChangedInput struct {
	EndDate time.Time `json:"end_date"`
}

func DecideEndDateChanged(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in EndDateChangedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if cur.EndDate.Equal(in.EndDate) {
		return nil, nil
	}
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = EndDateChangedMeta.FriendlyName

	return &EndDateChanged{
		Base:    base,
		EndDate: in.EndDate,
	}, nil
}

var EndDateChangedMeta = EventMeta{
	Type:         EndDateChangedType,
	FriendlyName: "End Date Change",
}

func init() {
	RegisterMeta(EndDateChangedMeta, func() Event {
		return &EndDateChanged{
			Base: Base{FriendlyNameStr: EndDateChangedMeta.FriendlyName},
		}
	})

	RegisterDecider[EndDateChangedInput](EndDateChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *project.Project, in EndDateChangedInput, status Status) (Event, error) {
			return DecideEndDateChanged(projectID, actor, cur, in, status)
		})

	RegisterInputType(EndDateChangedType, EndDateChangedInput{})
}
