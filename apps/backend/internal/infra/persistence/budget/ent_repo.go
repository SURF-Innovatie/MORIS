package budget

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	entbudget "github.com/SURF-Innovatie/MORIS/ent/budget"
	"github.com/SURF-Innovatie/MORIS/ent/budgetactual"
	"github.com/SURF-Innovatie/MORIS/ent/budgetlineitem"
	"github.com/SURF-Innovatie/MORIS/internal/app/budget"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) *EntRepository {
	return &EntRepository{client: client}
}

// Ensure EntRepository implements budget.Repository
var _ budget.Repository = (*EntRepository)(nil)

func (r *EntRepository) GetBudget(ctx context.Context, projectID uuid.UUID) (*ent.Budget, error) {
	b, err := r.client.Budget.
		Query().
		Where(entbudget.ProjectID(projectID)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, budget.ErrBudgetNotFound
		}
		return nil, err
	}

	return b, nil
}

func (r *EntRepository) GetBudgetByID(ctx context.Context, budgetID uuid.UUID) (*ent.Budget, error) {
	b, err := r.client.Budget.
		Query().
		Where(entbudget.ID(budgetID)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, budget.ErrBudgetNotFound
		}
		return nil, err
	}

	return b, nil
}

func (r *EntRepository) CreateBudget(ctx context.Context, projectID uuid.UUID, title, description string) (*ent.Budget, error) {
	return r.client.Budget.
		Create().
		SetProjectID(projectID).
		SetTitle(title).
		SetDescription(description).
		SetStatus(entbudget.StatusDraft).
		Save(ctx)
}

func (r *EntRepository) BudgetExists(ctx context.Context, projectID uuid.UUID) (bool, error) {
	return r.client.Budget.
		Query().
		Where(entbudget.ProjectID(projectID)).
		Exist(ctx)
}

func (r *EntRepository) GetLineItems(ctx context.Context, budgetID uuid.UUID) ([]*ent.BudgetLineItem, error) {
	return r.client.BudgetLineItem.
		Query().
		Where(budgetlineitem.BudgetID(budgetID)).
		WithActuals().
		All(ctx)
}

func (r *EntRepository) AddLineItem(ctx context.Context, budgetID uuid.UUID, item entities.BudgetLineItem) (*ent.BudgetLineItem, error) {
	builder := r.client.BudgetLineItem.
		Create().
		SetBudgetID(budgetID).
		SetCategory(budgetlineitem.Category(item.Category)).
		SetDescription(item.Description).
		SetBudgetedAmount(item.BudgetedAmount).
		SetYear(item.Year).
		SetFundingSource(budgetlineitem.FundingSource(item.FundingSource))

	if item.NWOGrantID != nil {
		builder.SetNwoGrantID(*item.NWOGrantID)
	}

	return builder.Save(ctx)
}

func (r *EntRepository) RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	err := r.client.BudgetLineItem.
		DeleteOneID(lineItemID).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return budget.ErrLineItemNotFound
		}
		return err
	}
	return nil
}

func (r *EntRepository) GetActuals(ctx context.Context, budgetID uuid.UUID) ([]*ent.BudgetActual, error) {
	return r.client.BudgetActual.
		Query().
		Where(budgetactual.HasLineItemWith(budgetlineitem.BudgetID(budgetID))).
		All(ctx)
}

func (r *EntRepository) RecordActual(ctx context.Context, actual entities.BudgetActual) (*ent.BudgetActual, error) {
	builder := r.client.BudgetActual.
		Create().
		SetLineItemID(actual.LineItemID).
		SetAmount(actual.Amount).
		SetDescription(actual.Description).
		SetRecordedDate(actual.RecordedDate)

	if actual.Source != "" {
		builder.SetSource(actual.Source)
	}
	if actual.ExternalRef != "" {
		builder.SetExternalRef(actual.ExternalRef)
	}

	return builder.Save(ctx)
}
