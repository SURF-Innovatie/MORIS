package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const BudgetCreatedType = "budget.created"

// BudgetCreated represents the event when a budget is created for a project
type BudgetCreated struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (BudgetCreated) isEvent()     {}
func (BudgetCreated) Type() string { return BudgetCreatedType }
func (e BudgetCreated) String() string {
	return fmt.Sprintf("Budget created: %s", e.Title)
}

func (e *BudgetCreated) NotificationMessage() string {
	return fmt.Sprintf("Budget '%s' has been created.", e.Title)
}

type BudgetCreatedInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func DecideBudgetCreated(
	projectID uuid.UUID,
	actor uuid.UUID,
	in BudgetCreatedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.Title == "" {
		return nil, errors.New("title is required")
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = BudgetCreatedMeta.FriendlyName

	return &BudgetCreated{
		Base:        base,
		Title:       in.Title,
		Description: in.Description,
	}, nil
}

var BudgetCreatedMeta = EventMeta{
	Type:         BudgetCreatedType,
	FriendlyName: "Budget Created",
}

func init() {
	RegisterMeta(BudgetCreatedMeta, func() Event {
		return &BudgetCreated{
			Base: Base{FriendlyNameStr: BudgetCreatedMeta.FriendlyName},
		}
	})

	RegisterDecider[BudgetCreatedInput](BudgetCreatedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in BudgetCreatedInput, status Status) (Event, error) {
			return DecideBudgetCreated(projectID, actor, in, status)
		})

	RegisterInputType(BudgetCreatedType, BudgetCreatedInput{})
}
