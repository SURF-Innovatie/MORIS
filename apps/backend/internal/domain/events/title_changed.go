package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type TitleChanged struct {
	Base
	Title string `json:"title"`
}

func (TitleChanged) isEvent()     {}
func (TitleChanged) Type() string { return TitleChangedType }
func (e TitleChanged) String() string {
	return fmt.Sprintf("Title changed: %s", e.Title)
}

func (e *TitleChanged) Apply(project *entities.Project) {
	project.Title = e.Title
}

func (e *TitleChanged) NotificationMessage() string {
	return fmt.Sprintf("Project title changed to '%s'.", e.Title)
}

func init() {
	RegisterMeta(EventMeta{
		Type:         TitleChangedType,
		FriendlyName: "Title Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &TitleChanged{} })
}
