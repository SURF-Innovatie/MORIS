package budget

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Repository defines the contract for budget data access
type Repository interface {
	GetBudget(ctx context.Context, projectID uuid.UUID) (*ent.Budget, error)
	GetBudgetByID(ctx context.Context, budgetID uuid.UUID) (*ent.Budget, error)
	CreateBudget(ctx context.Context, projectID uuid.UUID, title, description string) (*ent.Budget, error)
	BudgetExists(ctx context.Context, projectID uuid.UUID) (bool, error)

	GetLineItems(ctx context.Context, budgetID uuid.UUID) ([]*ent.BudgetLineItem, error)
	AddLineItem(ctx context.Context, budgetID uuid.UUID, item entities.BudgetLineItem) (*ent.BudgetLineItem, error)
	RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error

	GetActuals(ctx context.Context, budgetID uuid.UUID) ([]*ent.BudgetActual, error)
	RecordActual(ctx context.Context, actual entities.BudgetActual) (*ent.BudgetActual, error)
}
