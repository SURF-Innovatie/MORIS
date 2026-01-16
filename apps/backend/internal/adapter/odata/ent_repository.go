package odata

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/budget"
	"github.com/SURF-Innovatie/MORIS/ent/budgetlineitem"
	appOdata "github.com/SURF-Innovatie/MORIS/internal/app/odata"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// EntRepository implements OData repository using Ent ORM
type EntRepository struct {
	client *ent.Client
}

// NewEntRepository creates a new Ent-based OData repository
func NewEntRepository(client *ent.Client) *EntRepository {
	return &EntRepository{client: client}
}

// QueryBudgets returns budgets filtered by user access and OData query
func (r *EntRepository) QueryBudgets(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[appOdata.BudgetODataDTO], error) {
	q := r.client.Budget.Query()

	// TODO: Apply RBAC filtering based on userID
	// For now, return all budgets (in production, filter by project access)

	// Apply OData $filter
	if query.Filter != nil {
		q = applyBudgetFilter(q, query.Filter)
	}

	// Apply $orderby
	for _, order := range query.OrderBy {
		q = applyBudgetOrder(q, order)
	}

	// Get total count if requested
	var count *int
	if query.Count {
		total, err := q.Clone().Count(ctx)
		if err != nil {
			return entities.ODataResult[appOdata.BudgetODataDTO]{}, err
		}
		count = &total
	}

	// Apply $skip
	if query.Skip != nil {
		q = q.Offset(*query.Skip)
	}

	// Apply $top with default limit
	limit := 100
	if query.Top != nil {
		limit = *query.Top
	}
	q = q.Limit(limit)

	// Load line items for computed fields
	q = q.WithLineItems(func(lq *ent.BudgetLineItemQuery) {
		lq.WithActuals()
	})

	budgets, err := q.All(ctx)
	if err != nil {
		return entities.ODataResult[appOdata.BudgetODataDTO]{}, err
	}

	// Map to DTOs
	dtos := make([]appOdata.BudgetODataDTO, len(budgets))
	for i, b := range budgets {
		dtos[i] = mapBudgetToODataDTO(b)
	}

	return entities.ODataResult[appOdata.BudgetODataDTO]{
		Value: dtos,
		Count: count,
	}, nil
}

// QueryLineItems returns line items filtered by user access and OData query
func (r *EntRepository) QueryLineItems(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[appOdata.LineItemODataDTO], error) {
	q := r.client.BudgetLineItem.Query()

	// Apply $filter
	if query.Filter != nil {
		q = applyLineItemFilter(q, query.Filter)
	}

	// Apply $orderby
	for _, order := range query.OrderBy {
		q = applyLineItemOrder(q, order)
	}

	// Get count
	var count *int
	if query.Count {
		total, err := q.Clone().Count(ctx)
		if err != nil {
			return entities.ODataResult[appOdata.LineItemODataDTO]{}, err
		}
		count = &total
	}

	// Apply pagination
	if query.Skip != nil {
		q = q.Offset(*query.Skip)
	}
	limit := 100
	if query.Top != nil {
		limit = *query.Top
	}
	q = q.Limit(limit)

	// Load actuals for computed fields
	q = q.WithActuals()

	items, err := q.All(ctx)
	if err != nil {
		return entities.ODataResult[appOdata.LineItemODataDTO]{}, err
	}

	dtos := make([]appOdata.LineItemODataDTO, len(items))
	for i, item := range items {
		dtos[i] = mapLineItemToODataDTO(item)
	}

	return entities.ODataResult[appOdata.LineItemODataDTO]{
		Value: dtos,
		Count: count,
	}, nil
}

// QueryActuals returns actuals filtered by user access and OData query
func (r *EntRepository) QueryActuals(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[appOdata.ActualODataDTO], error) {
	q := r.client.BudgetActual.Query()

	// Apply pagination
	if query.Skip != nil {
		q = q.Offset(*query.Skip)
	}
	limit := 100
	if query.Top != nil {
		limit = *query.Top
	}
	q = q.Limit(limit)

	var count *int
	if query.Count {
		total, err := q.Clone().Count(ctx)
		if err != nil {
			return entities.ODataResult[appOdata.ActualODataDTO]{}, err
		}
		count = &total
	}

	actuals, err := q.All(ctx)
	if err != nil {
		return entities.ODataResult[appOdata.ActualODataDTO]{}, err
	}

	dtos := make([]appOdata.ActualODataDTO, len(actuals))
	for i, actual := range actuals {
		dtos[i] = appOdata.ActualODataDTO{
			ID:           actual.ID,
			LineItemID:   actual.LineItemID,
			Amount:       actual.Amount,
			Description:  actual.Description,
			RecordedDate: actual.RecordedDate.Format("2006-01-02"),
			Source:       actual.Source,
		}
	}

	return entities.ODataResult[appOdata.ActualODataDTO]{
		Value: dtos,
		Count: count,
	}, nil
}

// QueryAnalytics returns aggregated analytics
func (r *EntRepository) QueryAnalytics(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[appOdata.AnalyticsODataDTO], error) {
	// Get all budgets with line items and actuals
	budgets, err := r.client.Budget.
		Query().
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		All(ctx)

	if err != nil {
		return entities.ODataResult[appOdata.AnalyticsODataDTO]{}, err
	}

	dtos := make([]appOdata.AnalyticsODataDTO, len(budgets))
	for i, b := range budgets {
		var totalBudgeted, totalActuals float64
		for _, item := range b.Edges.LineItems {
			totalBudgeted += item.BudgetedAmount
			for _, actual := range item.Edges.Actuals {
				totalActuals += actual.Amount
			}
		}

		var burnRate float64
		if totalBudgeted > 0 {
			burnRate = (totalActuals / totalBudgeted) * 100
		}

		dtos[i] = appOdata.AnalyticsODataDTO{
			ProjectID:     b.ProjectID,
			BudgetID:      b.ID,
			TotalBudgeted: totalBudgeted,
			TotalActuals:  totalActuals,
			Remaining:     totalBudgeted - totalActuals,
			BurnRate:      burnRate,
			Status:        string(b.Status),
		}
	}

	count := len(dtos)
	return entities.ODataResult[appOdata.AnalyticsODataDTO]{
		Value: dtos,
		Count: &count,
	}, nil
}

// Helper functions

func mapBudgetToODataDTO(b *ent.Budget) appOdata.BudgetODataDTO {
	var totalBudgeted, totalActuals float64
	for _, item := range b.Edges.LineItems {
		totalBudgeted += item.BudgetedAmount
		for _, actual := range item.Edges.Actuals {
			totalActuals += actual.Amount
		}
	}

	var burnRate float64
	if totalBudgeted > 0 {
		burnRate = (totalActuals / totalBudgeted) * 100
	}

	return appOdata.BudgetODataDTO{
		ID:            b.ID,
		ProjectID:     b.ProjectID,
		Title:         b.Title,
		Status:        string(b.Status),
		TotalAmount:   b.TotalAmount,
		TotalBudgeted: totalBudgeted,
		TotalActuals:  totalActuals,
		BurnRate:      burnRate,
		Currency:      b.Currency,
	}
}

func mapLineItemToODataDTO(item *ent.BudgetLineItem) appOdata.LineItemODataDTO {
	var totalActuals float64
	for _, actual := range item.Edges.Actuals {
		totalActuals += actual.Amount
	}

	return appOdata.LineItemODataDTO{
		ID:             item.ID,
		BudgetID:       item.BudgetID,
		Category:       string(item.Category),
		Description:    item.Description,
		BudgetedAmount: item.BudgetedAmount,
		Year:           item.Year,
		FundingSource:  string(item.FundingSource),
		TotalActuals:   totalActuals,
	}
}

func applyBudgetFilter(q *ent.BudgetQuery, filter *entities.ODataFilter) *ent.BudgetQuery {
	if filter == nil {
		return q
	}

	switch filter.Field {
	case "projectId":
		if id, ok := filter.Value.(string); ok {
			if uid, err := uuid.Parse(id); err == nil {
				q = q.Where(budget.ProjectID(uid))
			}
		}
	case "status":
		if status, ok := filter.Value.(string); ok {
			q = q.Where(budget.StatusEQ(budget.Status(status)))
		}
	}

	if filter.And != nil {
		q = applyBudgetFilter(q, filter.And)
	}

	return q
}

func applyBudgetOrder(q *ent.BudgetQuery, order entities.ODataOrderBy) *ent.BudgetQuery {
	switch order.Field {
	case "title":
		if order.Desc {
			return q.Order(ent.Desc(budget.FieldTitle))
		}
		return q.Order(ent.Asc(budget.FieldTitle))
	case "createdAt":
		if order.Desc {
			return q.Order(ent.Desc(budget.FieldCreatedAt))
		}
		return q.Order(ent.Asc(budget.FieldCreatedAt))
	case "totalAmount":
		if order.Desc {
			return q.Order(ent.Desc(budget.FieldTotalAmount))
		}
		return q.Order(ent.Asc(budget.FieldTotalAmount))
	}
	return q
}

func applyLineItemFilter(q *ent.BudgetLineItemQuery, filter *entities.ODataFilter) *ent.BudgetLineItemQuery {
	if filter == nil {
		return q
	}

	switch filter.Field {
	case "budgetId":
		if id, ok := filter.Value.(string); ok {
			if uid, err := uuid.Parse(id); err == nil {
				q = q.Where(budgetlineitem.BudgetID(uid))
			}
		}
	case "category":
		if cat, ok := filter.Value.(string); ok {
			q = q.Where(budgetlineitem.CategoryEQ(budgetlineitem.Category(cat)))
		}
	case "year":
		if year, ok := filter.Value.(int); ok {
			q = q.Where(budgetlineitem.Year(year))
		}
	case "fundingSource":
		if fs, ok := filter.Value.(string); ok {
			q = q.Where(budgetlineitem.FundingSourceEQ(budgetlineitem.FundingSource(fs)))
		}
	}

	if filter.And != nil {
		q = applyLineItemFilter(q, filter.And)
	}

	return q
}

func applyLineItemOrder(q *ent.BudgetLineItemQuery, order entities.ODataOrderBy) *ent.BudgetLineItemQuery {
	switch order.Field {
	case "category":
		if order.Desc {
			return q.Order(ent.Desc(budgetlineitem.FieldCategory))
		}
		return q.Order(ent.Asc(budgetlineitem.FieldCategory))
	case "budgetedAmount":
		if order.Desc {
			return q.Order(ent.Desc(budgetlineitem.FieldBudgetedAmount))
		}
		return q.Order(ent.Asc(budgetlineitem.FieldBudgetedAmount))
	case "year":
		if order.Desc {
			return q.Order(ent.Desc(budgetlineitem.FieldYear))
		}
		return q.Order(ent.Asc(budgetlineitem.FieldYear))
	}
	return q
}
