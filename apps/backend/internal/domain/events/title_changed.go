package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
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
	Title string
}

func DecideTitleChanged(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in TitleChangedInput,
	status Status,
) (*TitleChanged, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if in.Title == "" || cur.Title == in.Title {
		return nil, nil
	}

	return &TitleChanged{
		Base:  NewBase(projectID, actor, status),
		Title: in.Title,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         TitleChangedType,
		FriendlyName: "Title Change",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &TitleChanged{} })

	RegisterDecider[TitleChangedInput](TitleChangedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur any, in TitleChangedInput, status Status) (Event, error) {
			p, ok := cur.(*entities.Project)
			if !ok {
				return nil, fmt.Errorf("expected *entities.Project, got %T", cur)
			}
			return DecideTitleChanged(projectID, actor, p, in, status)
		},
	)
}
