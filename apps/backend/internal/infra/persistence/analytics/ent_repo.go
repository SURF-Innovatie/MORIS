package analytics

import (
	"context"
	"encoding/json"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/budget"
	"github.com/SURF-Innovatie/MORIS/ent/budgetactual"
	"github.com/SURF-Innovatie/MORIS/ent/budgetlineitem"
	"github.com/SURF-Innovatie/MORIS/ent/event"
	"github.com/SURF-Innovatie/MORIS/ent/organisationnodeclosure"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/google/uuid"
)

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) *EntRepository {
	return &EntRepository{client: client}
}

// Helper to get all descendant organization IDs (including the root itself)
func (r *EntRepository) getDescendantOrgIDs(ctx context.Context, rootOrgID uuid.UUID) ([]uuid.UUID, error) {
	// fetching descendants from the closure table
	descendants, err := r.client.OrganisationNodeClosure.Query().
		Where(organisationnodeclosure.AncestorID(rootOrgID)).
		Select(organisationnodeclosure.FieldDescendantID).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, d := range descendants {
		ids = append(ids, d.DescendantID)
	}

	// Ensure the root itself is included if the closure table doesn't guarantee self-reference (it usually should, depth 0)
	// But let's be safe or just rely on the query result.
	// If the closure table is correctly maintained, it includes depth 0 (self).
	// If for some reason the list is empty (orphan node?), we should at least return the root.
	if len(ids) == 0 {
		ids = append(ids, rootOrgID)
	}

	return ids, nil
}

// Helper to get project IDs for an organization by scanning events
func (r *EntRepository) getProjectIDsForOrg(ctx context.Context, orgID uuid.UUID) ([]uuid.UUID, error) {
	// 1. Get all relevant org IDs (self + descendants)
	orgIDs, err := r.getDescendantOrgIDs(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Create a map for O(1) lookups
	orgIDMap := make(map[uuid.UUID]bool)
	for _, id := range orgIDs {
		orgIDMap[id] = true
	}

	// 2. Fetch all relevant events (Started + OwningOrgNodeChanged).
	// We need to replay them to know the *current* owning org.
	evts, err := r.client.Event.Query().
		Where(
			event.Or(
				event.TypeEQ(events.ProjectStartedType),
				event.TypeEQ(events.OwningOrgNodeChangedType),
			),
		).
		Order(ent.Asc(event.FieldOccurredAt)). // Ensure chronological order
		All(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Replay events to build current state of project ownership
	projectOrgMap := make(map[uuid.UUID]uuid.UUID)

	for _, e := range evts {
		var data struct {
			OwningOrgNodeID uuid.UUID `json:"owning_org_node_id"`
		}

		// Unmarshal event data to extract org node ID
		b, _ := json.Marshal(e.Data)
		if err := json.Unmarshal(b, &data); err != nil {
			continue // Should not happen if data is valid
		}

		// Update the mapping for this project
		if data.OwningOrgNodeID != uuid.Nil {
			projectOrgMap[e.ProjectID] = data.OwningOrgNodeID
		}
	}

	// 4. Filter projects that belong to one of the target organizations
	var projectIDs []uuid.UUID
	for projID, ownerID := range projectOrgMap {
		if orgIDMap[ownerID] {
			projectIDs = append(projectIDs, projID)
		}
	}
	return projectIDs, nil
}

func (r *EntRepository) GetOrgSummary(ctx context.Context, orgID uuid.UUID) (entities.OrgAnalyticsSummary, error) {
	projectIDs, err := r.getProjectIDsForOrg(ctx, orgID)
	if err != nil {
		return entities.OrgAnalyticsSummary{}, err
	}
	if len(projectIDs) == 0 {
		return entities.OrgAnalyticsSummary{}, nil
	}

	budgets, err := r.client.Budget.Query().
		Where(budget.ProjectIDIn(projectIDs...)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		All(ctx)

	if err != nil {
		return entities.OrgAnalyticsSummary{}, err
	}

	var summary entities.OrgAnalyticsSummary
	summary.TotalProjects = len(budgets)

	for _, b := range budgets {
		var budgetTotal float64
		// Calculate budgeted and actuals by iterating line items
		for _, li := range b.Edges.LineItems {
			budgetTotal += li.BudgetedAmount
			for _, act := range li.Edges.Actuals {
				summary.TotalActuals += act.Amount
			}
		}
		summary.TotalBudgeted += budgetTotal
	}

	if summary.TotalProjects > 0 {
		for _, b := range budgets {
			var actuals float64
			var budgetAmount float64

			for _, li := range b.Edges.LineItems {
				budgetAmount += li.BudgetedAmount
				for _, act := range li.Edges.Actuals {
					actuals += act.Amount
				}
			}

			if budgetAmount > 0 {
				ratio := actuals / budgetAmount
				if ratio > 1.0 {
					summary.ProjectsAtRisk++
				} else {
					summary.ProjectsOnTrack++
				}

				summary.AverageSpendRate += ratio
			}
		}
		summary.AverageSpendRate = summary.AverageSpendRate / float64(summary.TotalProjects) * 100
	}

	return summary, nil
}

func (r *EntRepository) GetBurnRateData(ctx context.Context, orgID uuid.UUID, params entities.DateRangeParams) ([]entities.BurnRateDataPoint, error) {
	projectIDs, err := r.getProjectIDsForOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if len(projectIDs) == 0 {
		return []entities.BurnRateDataPoint{}, nil
	}

	query := r.client.BudgetActual.Query().
		Where(budgetactual.HasLineItemWith(budgetlineitem.HasBudgetWith(budget.ProjectIDIn(projectIDs...))))

	if params.StartDate != nil {
		query.Where(budgetactual.RecordedDateGTE(*params.StartDate))
	}
	if params.EndDate != nil {
		query.Where(budgetactual.RecordedDateLTE(*params.EndDate))
	}

	actuals, err := query.Order(ent.Asc(budgetactual.FieldRecordedDate)).All(ctx)
	if err != nil {
		return nil, err
	}

	var points []entities.BurnRateDataPoint
	var runningTotal float64

	for _, act := range actuals {
		runningTotal += act.Amount
		points = append(points, entities.BurnRateDataPoint{
			Date:     act.RecordedDate,
			Actual:   runningTotal,
			Budgeted: 0,
		})
	}

	return points, nil
}

func (r *EntRepository) GetCategoryBreakdown(ctx context.Context, orgID uuid.UUID) ([]entities.CategoryBreakdown, error) {
	projectIDs, err := r.getProjectIDsForOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if len(projectIDs) == 0 {
		return []entities.CategoryBreakdown{}, nil
	}

	items, err := r.client.BudgetLineItem.Query().
		Where(budgetlineitem.HasBudgetWith(budget.ProjectIDIn(projectIDs...))).
		WithActuals().
		All(ctx)

	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]*entities.CategoryBreakdown)

	for _, item := range items {
		cat := item.Category.String()
		if _, exists := categoryMap[cat]; !exists {
			categoryMap[cat] = &entities.CategoryBreakdown{Category: cat}
		}

		budgeted := item.BudgetedAmount
		categoryMap[cat].Budgeted += budgeted

		var actuals float64
		for _, act := range item.Edges.Actuals {
			actuals += act.Amount
		}
		categoryMap[cat].Actuals += actuals
		categoryMap[cat].Remaining += (budgeted - actuals)
	}

	var result []entities.CategoryBreakdown
	for _, v := range categoryMap {
		result = append(result, *v)
	}

	return result, nil
}

func (r *EntRepository) GetProjectHealth(ctx context.Context, orgID uuid.UUID) ([]entities.ProjectHealthSummary, error) {
	projectIDs, err := r.getProjectIDsForOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if len(projectIDs) == 0 {
		return []entities.ProjectHealthSummary{}, nil
	}

	budgets, err := r.client.Budget.Query().
		Where(budget.ProjectIDIn(projectIDs...)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		All(ctx)

	if err != nil {
		return nil, err
	}

	var healths []entities.ProjectHealthSummary

	for _, b := range budgets {
		var actuals float64
		var budgeted float64
		for _, li := range b.Edges.LineItems {
			budgeted += li.BudgetedAmount
			for _, act := range li.Edges.Actuals {
				actuals += act.Amount
			}
		}

		var burnRate float64
		if budgeted > 0 {
			burnRate = (actuals / budgeted) * 100
		}

		status := "on_track"
		if burnRate > 100 {
			status = "at_risk"
		} else if burnRate > 80 {
			status = "warning"
		}

		h := entities.ProjectHealthSummary{
			ProjectID:   b.ProjectID,
			BudgetID:    b.ID,
			Budgeted:    budgeted,
			Spent:       actuals,
			Remaining:   budgeted - actuals,
			BurnRate:    burnRate,
			Status:      status,
			ProjectName: b.Title, // Using Budget Title as proxy for Project Name since Project entity missing
		}

		healths = append(healths, h)
	}

	return healths, nil
}

func (r *EntRepository) GetFundingBreakdown(ctx context.Context, orgID uuid.UUID) ([]entities.FundingBreakdown, error) {
	projectIDs, err := r.getProjectIDsForOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if len(projectIDs) == 0 {
		return []entities.FundingBreakdown{}, nil
	}

	items, err := r.client.BudgetLineItem.Query().
		Where(budgetlineitem.HasBudgetWith(budget.ProjectIDIn(projectIDs...))).
		WithActuals().
		All(ctx)

	if err != nil {
		return nil, err
	}

	fundingMap := make(map[string]*entities.FundingBreakdown)

	for _, item := range items {
		source := item.FundingSource.String()
		if source == "" {
			source = "Unspecified"
		}

		if _, exists := fundingMap[source]; !exists {
			fundingMap[source] = &entities.FundingBreakdown{FundingSource: source}
		}

		budgeted := item.BudgetedAmount
		fundingMap[source].Budgeted += budgeted

		var actuals float64
		for _, act := range item.Edges.Actuals {
			actuals += act.Amount
		}
		fundingMap[source].Actuals += actuals
		fundingMap[source].Remaining += (budgeted - actuals)
	}

	var result []entities.FundingBreakdown
	for _, v := range fundingMap {
		result = append(result, *v)
	}

	return result, nil
}
