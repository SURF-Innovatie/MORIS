package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const BudgetLineItemAddedType = "budget.line_item_added"

// BudgetLineItemAdded represents the event when a line item is added to a budget
type BudgetLineItemAdded struct {
	Base
	LineItemID     uuid.UUID               `json:"lineItemId"`
	Category       entities.BudgetCategory `json:"category"`
	Description    string                  `json:"description"`
	BudgetedAmount float64                 `json:"budgetedAmount"`
	Year           int                     `json:"year"`
	FundingSource  entities.FundingSource  `json:"fundingSource"`
}

func (BudgetLineItemAdded) isEvent()     {}
func (BudgetLineItemAdded) Type() string { return BudgetLineItemAddedType }
func (e BudgetLineItemAdded) String() string {
	return fmt.Sprintf("Budget line item added: %s (€%.2f)", e.Description, e.BudgetedAmount)
}

func (e *BudgetLineItemAdded) NotificationMessage() string {
	return fmt.Sprintf("Budget line item '%s' (€%.2f) has been added.", e.Description, e.BudgetedAmount)
}

type BudgetLineItemAddedInput struct {
	Category       entities.BudgetCategory `json:"category"`
	Description    string                  `json:"description"`
	BudgetedAmount float64                 `json:"budgetedAmount"`
	Year           int                     `json:"year"`
	FundingSource  entities.FundingSource  `json:"fundingSource"`
}

func DecideBudgetLineItemAdded(
	projectID uuid.UUID,
	actor uuid.UUID,
	in BudgetLineItemAddedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.Description == "" {
		return nil, errors.New("description is required")
	}
	if in.BudgetedAmount <= 0 {
		return nil, errors.New("budgeted amount must be positive")
	}
	if in.Year < 2000 || in.Year > 2100 {
		return nil, errors.New("invalid year")
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = BudgetLineItemAddedMeta.FriendlyName

	return &BudgetLineItemAdded{
		Base:           base,
		LineItemID:     uuid.New(),
		Category:       in.Category,
		Description:    in.Description,
		BudgetedAmount: in.BudgetedAmount,
		Year:           in.Year,
		FundingSource:  in.FundingSource,
	}, nil
}

var BudgetLineItemAddedMeta = EventMeta{
	Type:         BudgetLineItemAddedType,
	FriendlyName: "Budget Line Item Added",
}

func init() {
	RegisterMeta(BudgetLineItemAddedMeta, func() Event {
		return &BudgetLineItemAdded{
			Base: Base{FriendlyNameStr: BudgetLineItemAddedMeta.FriendlyName},
		}
	})

	RegisterDecider[BudgetLineItemAddedInput](BudgetLineItemAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in BudgetLineItemAddedInput, status Status) (Event, error) {
			return DecideBudgetLineItemAdded(projectID, actor, in, status)
		})

	RegisterInputType(BudgetLineItemAddedType, BudgetLineItemAddedInput{})
}
