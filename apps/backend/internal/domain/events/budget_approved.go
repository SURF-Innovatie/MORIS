package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const BudgetApprovedType = "budget.approved"

// BudgetApproved represents the event when a budget is approved and locked
type BudgetApproved struct {
	Base
	Title       string  `json:"title"`
	TotalAmount float64 `json:"totalAmount"`
}

func (BudgetApproved) isEvent()     {}
func (BudgetApproved) Type() string { return BudgetApprovedType }
func (e BudgetApproved) String() string {
	return fmt.Sprintf("Budget approved: %s (€%.2f)", e.Title, e.TotalAmount)
}

func (e *BudgetApproved) NotificationMessage() string {
	return fmt.Sprintf("Budget '%s' (€%.2f) has been approved and locked.", e.Title, e.TotalAmount)
}

type BudgetApprovedInput struct {
	Title       string  `json:"title"`
	TotalAmount float64 `json:"totalAmount"`
}

func DecideBudgetApproved(
	projectID uuid.UUID,
	actor uuid.UUID,
	in BudgetApprovedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = BudgetApprovedMeta.FriendlyName

	return &BudgetApproved{
		Base:        base,
		Title:       in.Title,
		TotalAmount: in.TotalAmount,
	}, nil
}

var BudgetApprovedMeta = EventMeta{
	Type:         BudgetApprovedType,
	FriendlyName: "Budget Approved",
}

func init() {
	RegisterMeta(BudgetApprovedMeta, func() Event {
		return &BudgetApproved{
			Base: Base{FriendlyNameStr: BudgetApprovedMeta.FriendlyName},
		}
	})

	RegisterDecider[BudgetApprovedInput](BudgetApprovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in BudgetApprovedInput, status Status) (Event, error) {
			return DecideBudgetApproved(projectID, actor, in, status)
		})

	RegisterInputType(BudgetApprovedType, BudgetApprovedInput{})
}
