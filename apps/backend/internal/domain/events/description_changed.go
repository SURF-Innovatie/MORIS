package events

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type DescriptionChanged struct {
	Base
	Description string `json:"description"`
}

func (DescriptionChanged) isEvent()     {}
func (DescriptionChanged) Type() string { return DescriptionChangedType }
func (e DescriptionChanged) String() string {
	return "Description changed"
}

func (e *DescriptionChanged) Apply(project *entities.Project) {
	project.Description = e.Description
}

func (e *DescriptionChanged) NotificationMessage() string {
	return "Project description has been updated."
}

func init() {
	RegisterMeta(EventMeta{
		Type:         DescriptionChangedType,
		FriendlyName: "Description Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &DescriptionChanged{} })
}
