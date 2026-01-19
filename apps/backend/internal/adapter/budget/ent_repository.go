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

func (r *EntRepository) GetBudget(ctx context.Context, projectID uuid.UUID) (*entities.Budget, error) {
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

	return mapEntToBudget(b), nil
}

func (r *EntRepository) GetBudgetByID(ctx context.Context, budgetID uuid.UUID) (*entities.Budget, error) {
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

	return mapEntToBudget(b), nil
}

func (r *EntRepository) CreateBudget(ctx context.Context, projectID uuid.UUID, title, description string) (*entities.Budget, error) {
	b, err := r.client.Budget.
		Create().
		SetProjectID(projectID).
		SetTitle(title).
		SetDescription(description).
		SetStatus(entbudget.StatusDraft).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	return mapEntToBudget(b), nil
}

func (r *EntRepository) BudgetExists(ctx context.Context, projectID uuid.UUID) (bool, error) {
	return r.client.Budget.
		Query().
		Where(entbudget.ProjectID(projectID)).
		Exist(ctx)
}

func (r *EntRepository) GetLineItems(ctx context.Context, budgetID uuid.UUID) ([]entities.BudgetLineItem, error) {
	items, err := r.client.BudgetLineItem.
		Query().
		Where(budgetlineitem.BudgetID(budgetID)).
		WithActuals().
		All(ctx)

	if err != nil {
		return nil, err
	}

	var result []entities.BudgetLineItem
	for _, item := range items {
		result = append(result, mapEntToLineItem(item))
	}
	return result, nil
}

func (r *EntRepository) AddLineItem(ctx context.Context, budgetID uuid.UUID, item entities.BudgetLineItem) (*entities.BudgetLineItem, error) {
	created, err := r.client.BudgetLineItem.
		Create().
		SetBudgetID(budgetID).
		SetCategory(budgetlineitem.Category(item.Category)).
		SetDescription(item.Description).
		SetBudgetedAmount(item.BudgetedAmount).
		SetYear(item.Year).
		SetFundingSource(budgetlineitem.FundingSource(item.FundingSource)).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	// We need to fetch it back to get edges safely or just map what we have (no actuals yet)
	res := mapEntToLineItem(created)
	return &res, nil
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

func (r *EntRepository) GetActuals(ctx context.Context, budgetID uuid.UUID) ([]entities.BudgetActual, error) {
	actuals, err := r.client.BudgetActual.
		Query().
		Where(budgetactual.HasLineItemWith(budgetlineitem.BudgetID(budgetID))).
		All(ctx)

	if err != nil {
		return nil, err
	}

	var result []entities.BudgetActual
	for _, a := range actuals {
		result = append(result, mapEntToActual(a))
	}
	return result, nil
}

func (r *EntRepository) RecordActual(ctx context.Context, actual entities.BudgetActual) (*entities.BudgetActual, error) {
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

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	res := mapEntToActual(created)
	return &res, nil
}

// Mappers

func mapEntToBudget(b *ent.Budget) *entities.Budget {
	var items []entities.BudgetLineItem
	for _, item := range b.Edges.LineItems {
		items = append(items, mapEntToLineItem(item))
	}

	return &entities.Budget{
		ID:          b.ID,
		ProjectID:   b.ProjectID,
		Title:       b.Title,
		Description: b.Description,
		Status:      entities.BudgetStatus(b.Status),
		TotalAmount: b.TotalAmount,
		Currency:    b.Currency,
		Version:     b.Version,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
		LineItems:   items,
	}
}

func mapEntToLineItem(li *ent.BudgetLineItem) entities.BudgetLineItem {
	var actuals []entities.BudgetActual
	for _, a := range li.Edges.Actuals {
		actuals = append(actuals, mapEntToActual(a))
	}

	return entities.BudgetLineItem{
		ID:             li.ID,
		BudgetID:       li.BudgetID,
		Category:       entities.BudgetCategory(li.Category),
		Description:    li.Description,
		BudgetedAmount: li.BudgetedAmount,
		Year:           li.Year,
		FundingSource:  entities.FundingSource(li.FundingSource),
		Actuals:        actuals,
	}
}

func mapEntToActual(a *ent.BudgetActual) entities.BudgetActual {
	return entities.BudgetActual{
		ID:           a.ID,
		LineItemID:   a.LineItemID,
		Amount:       a.Amount,
		Description:  a.Description,
		RecordedDate: a.RecordedDate,
		Source:       a.Source,
		ExternalRef:  a.ExternalRef,
	}
}
