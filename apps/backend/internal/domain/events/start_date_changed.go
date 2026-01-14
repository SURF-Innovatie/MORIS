package events

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	StartDate time.Time `json:"start_date"`
}

func DecideStartDateChanged(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in StartDateChangedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if cur.StartDate.Equal(in.StartDate) {
		return nil, nil
	}
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = StartDateChangedMeta.FriendlyName

	return &StartDateChanged{
		Base:      base,
		StartDate: in.StartDate,
	}, nil
}

var StartDateChangedMeta = EventMeta{
	Type:         StartDateChangedType,
	FriendlyName: "Start Date Change",
}

func init() {
	RegisterMeta(StartDateChangedMeta, func() Event {
		return &StartDateChanged{
			Base: Base{FriendlyNameStr: StartDateChangedMeta.FriendlyName},
		}
	})

	RegisterDecider[StartDateChangedInput](StartDateChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in StartDateChangedInput, status Status) (Event, error) {
			return DecideStartDateChanged(projectID, actor, cur, in, status)
		})

	RegisterInputType(StartDateChangedType, StartDateChangedInput{})
}
