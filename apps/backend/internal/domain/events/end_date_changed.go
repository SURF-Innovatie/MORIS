package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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

func (e *EndDateChanged) Apply(project *entities.Project) {
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
	cur *entities.Project,
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
	return &EndDateChanged{
		Base:    NewBase(projectID, actor, status),
		EndDate: in.EndDate,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         EndDateChangedType,
		FriendlyName: "End Date Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &EndDateChanged{} })

	RegisterDecider[EndDateChangedInput](EndDateChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in EndDateChangedInput, status Status) (Event, error) {
			return DecideEndDateChanged(projectID, actor, cur, in, status)
		})

	RegisterInputType(EndDateChangedType, EndDateChangedInput{})
}
