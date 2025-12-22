package events

import (
	"context"
	"fmt"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

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

func init() {
	RegisterMeta(EventMeta{
		Type:         EndDateChangedType,
		FriendlyName: "End Date Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &EndDateChanged{} })
}
