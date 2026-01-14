package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const TitleChangedType = "project.title_changed"

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

type TitleChangedInput struct {
	Title string `json:"title"`
}

func DecideTitleChanged(
	id uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in TitleChangedInput,
	status Status,
) (Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("project is nil")
	}

	// no-op rules
	if in.Title == "" {
		return nil, nil
	}
	if cur.Title == in.Title {
		return nil, nil
	}

	return &TitleChanged{
		Base:  NewBase(id, actor, status),
		Title: in.Title,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         TitleChangedType,
		FriendlyName: "Title Change",
	}, func() Event { return &TitleChanged{} })

	RegisterDecider[TitleChangedInput](TitleChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in TitleChangedInput, status Status) (Event, error) {
			return DecideTitleChanged(projectID, actor, cur, in, status)
		},
	)

	RegisterInputType(TitleChangedType, TitleChangedInput{})
}
