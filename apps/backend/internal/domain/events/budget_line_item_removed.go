package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const BudgetLineItemRemovedType = "budget.line_item_removed"

// BudgetLineItemRemoved represents the event when a line item is removed from a budget
type BudgetLineItemRemoved struct {
	Base
	LineItemID  uuid.UUID `json:"lineItemId"`
	Description string    `json:"description"` // For audit trail
}

func (BudgetLineItemRemoved) isEvent()     {}
func (BudgetLineItemRemoved) Type() string { return BudgetLineItemRemovedType }
func (e BudgetLineItemRemoved) String() string {
	return fmt.Sprintf("Budget line item removed: %s", e.Description)
}

func (e *BudgetLineItemRemoved) NotificationMessage() string {
	return fmt.Sprintf("Budget line item '%s' has been removed.", e.Description)
}

type BudgetLineItemRemovedInput struct {
	LineItemID  uuid.UUID `json:"lineItemId"`
	Description string    `json:"description"` // Captured for audit trail
}

func DecideBudgetLineItemRemoved(
	projectID uuid.UUID,
	actor uuid.UUID,
	in BudgetLineItemRemovedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.LineItemID == uuid.Nil {
		return nil, errors.New("line item id is required")
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = BudgetLineItemRemovedMeta.FriendlyName

	return &BudgetLineItemRemoved{
		Base:        base,
		LineItemID:  in.LineItemID,
		Description: in.Description,
	}, nil
}

var BudgetLineItemRemovedMeta = EventMeta{
	Type:         BudgetLineItemRemovedType,
	FriendlyName: "Budget Line Item Removed",
}

func init() {
	RegisterMeta(BudgetLineItemRemovedMeta, func() Event {
		return &BudgetLineItemRemoved{
			Base: Base{FriendlyNameStr: BudgetLineItemRemovedMeta.FriendlyName},
		}
	})

	RegisterDecider[BudgetLineItemRemovedInput](BudgetLineItemRemovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in BudgetLineItemRemovedInput, status Status) (Event, error) {
			return DecideBudgetLineItemRemoved(projectID, actor, in, status)
		})

	RegisterInputType(BudgetLineItemRemovedType, BudgetLineItemRemovedInput{})
}
