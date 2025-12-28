package events

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const DescriptionChangedType = "project.description_changed"

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

type DescriptionChangedInput struct {
	Description string `json:"description"`
}

func DecideDescriptionChanged(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in DescriptionChangedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if cur.Description == in.Description {
		return nil, nil
	}
	return &DescriptionChanged{
		Base:        NewBase(projectID, actor, status),
		Description: in.Description,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         DescriptionChangedType,
		FriendlyName: "Description Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &DescriptionChanged{} })

	RegisterDecider[DescriptionChangedInput](DescriptionChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in DescriptionChangedInput, status Status) (Event, error) {
			return DecideDescriptionChanged(projectID, actor, cur, in, status)
		},
	)

	RegisterInputType(DescriptionChangedType, DescriptionChangedInput{})
}
