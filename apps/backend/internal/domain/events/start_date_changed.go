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

const StartDateChangedType = "project.start_date_changed"

type StartDateChanged struct {
	Base
	StartDate time.Time `json:"startDate"`
}

func (StartDateChanged) isEvent()     {}
func (StartDateChanged) Type() string { return StartDateChangedType }
func (e StartDateChanged) String() string {
	return fmt.Sprintf("Start date changed: %s", e.StartDate.Format("2006-01-02"))
}

func (e *StartDateChanged) Apply(project *entities.Project) {
	project.StartDate = e.StartDate
}

func (e *StartDateChanged) NotificationMessage() string {
	return fmt.Sprintf("Project start date changed to %s.", e.StartDate.Format("2006-01-02"))
}

type StartDateChangedInput struct {
	StartDate time.Time
}

func DecideStartDateChanged(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in StartDateChangedInput,
	status Status,
) (*StartDateChanged, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if cur.StartDate.Equal(in.StartDate) {
		return nil, nil
	}
	return &StartDateChanged{
		Base:      NewBase(projectID, actor, status),
		StartDate: in.StartDate,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         StartDateChangedType,
		FriendlyName: "Start Date Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &StartDateChanged{} })

	RegisterDecider[StartDateChangedInput](StartDateChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur any, in StartDateChangedInput, status Status) (Event, error) {
			p := cur.(*entities.Project)
			return DecideStartDateChanged(projectID, actor, p, in, status)
		})
}
