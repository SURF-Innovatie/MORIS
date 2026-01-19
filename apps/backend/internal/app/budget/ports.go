package budget

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Repository defines the contract for budget data access
type Repository interface {
	GetBudget(ctx context.Context, projectID uuid.UUID) (*entities.Budget, error)
	GetBudgetByID(ctx context.Context, budgetID uuid.UUID) (*entities.Budget, error)
	CreateBudget(ctx context.Context, projectID uuid.UUID, title, description string) (*entities.Budget, error)
	BudgetExists(ctx context.Context, projectID uuid.UUID) (bool, error)

	GetLineItems(ctx context.Context, budgetID uuid.UUID) ([]entities.BudgetLineItem, error)
	AddLineItem(ctx context.Context, budgetID uuid.UUID, item entities.BudgetLineItem) (*entities.BudgetLineItem, error)
	RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error

	GetActuals(ctx context.Context, budgetID uuid.UUID) ([]entities.BudgetActual, error)
	RecordActual(ctx context.Context, actual entities.BudgetActual) (*entities.BudgetActual, error)
}
