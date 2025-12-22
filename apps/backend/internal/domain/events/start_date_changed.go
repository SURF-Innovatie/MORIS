package events

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

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

func init() {
	RegisterMeta(EventMeta{
		Type:         StartDateChangedType,
		FriendlyName: "Start Date Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &StartDateChanged{} })
}
