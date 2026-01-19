package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// OrgAnalyticsSummaryResponse provides organization-level analytics
type OrgAnalyticsSummaryResponse struct {
	TotalProjects    int     `json:"totalProjects"`
	TotalBudgeted    float64 `json:"totalBudgeted"`
	TotalActuals     float64 `json:"totalActuals"`
	AverageSpendRate float64 `json:"averageSpendRate"`
	ProjectsAtRisk   int     `json:"projectsAtRisk"`
	ProjectsOnTrack  int     `json:"projectsOnTrack"`
}

func (d OrgAnalyticsSummaryResponse) FromEntity(e entities.OrgAnalyticsSummary) OrgAnalyticsSummaryResponse {
	return OrgAnalyticsSummaryResponse{
		TotalProjects:    e.TotalProjects,
		TotalBudgeted:    e.TotalBudgeted,
		TotalActuals:     e.TotalActuals,
		AverageSpendRate: e.AverageSpendRate,
		ProjectsAtRisk:   e.ProjectsAtRisk,
		ProjectsOnTrack:  e.ProjectsOnTrack,
	}
}

// BurnRateDataPointResponse represents a point in burn rate time series
type BurnRateDataPointResponse struct {
	Date        time.Time `json:"date"`
	ProjectID   uuid.UUID `json:"projectId"`
	ProjectName string    `json:"projectName"`
	Budgeted    float64   `json:"budgeted"`
	Actual      float64   `json:"actual"`
	Projected   float64   `json:"projected"`
}

func (d BurnRateDataPointResponse) FromEntity(e entities.BurnRateDataPoint) BurnRateDataPointResponse {
	return BurnRateDataPointResponse{
		Date:        e.Date,
		ProjectID:   e.ProjectID,
		ProjectName: e.ProjectName,
		Budgeted:    e.Budgeted,
		Actual:      e.Actual,
		Projected:   e.Projected,
	}
}

// CategoryBreakdownResponse shows spending by category
type CategoryBreakdownResponse struct {
	Category  string  `json:"category"`
	Budgeted  float64 `json:"budgeted"`
	Actuals   float64 `json:"actuals"`
	Remaining float64 `json:"remaining"`
}

func (d CategoryBreakdownResponse) FromEntity(e entities.CategoryBreakdown) CategoryBreakdownResponse {
	return CategoryBreakdownResponse{
		Category:  e.Category,
		Budgeted:  e.Budgeted,
		Actuals:   e.Actuals,
		Remaining: e.Remaining,
	}
}

// ProjectHealthSummaryResponse shows health indicators for a project
type ProjectHealthSummaryResponse struct {
	ProjectID   uuid.UUID `json:"projectId"`
	ProjectName string    `json:"projectName"`
	BudgetID    uuid.UUID `json:"budgetId"`
	Budgeted    float64   `json:"budgeted"`
	Spent       float64   `json:"spent"`
	Remaining   float64   `json:"remaining"`
	BurnRate    float64   `json:"burnRate"`
	Status      string    `json:"status"` // "on_track", "warning", "at_risk"
}

func (d ProjectHealthSummaryResponse) FromEntity(e entities.ProjectHealthSummary) ProjectHealthSummaryResponse {
	return ProjectHealthSummaryResponse{
		ProjectID:   e.ProjectID,
		ProjectName: e.ProjectName,
		BudgetID:    e.BudgetID,
		Budgeted:    e.Budgeted,
		Spent:       e.Spent,
		Remaining:   e.Remaining,
		BurnRate:    e.BurnRate,
		Status:      e.Status,
	}
}

// FundingBreakdownResponse shows spending by funding source
type FundingBreakdownResponse struct {
	FundingSource string  `json:"fundingSource"`
	Budgeted      float64 `json:"budgeted"`
	Actuals       float64 `json:"actuals"`
	Remaining     float64 `json:"remaining"`
}

func (d FundingBreakdownResponse) FromEntity(e entities.FundingBreakdown) FundingBreakdownResponse {
	return FundingBreakdownResponse{
		FundingSource: e.FundingSource,
		Budgeted:      e.Budgeted,
		Actuals:       e.Actuals,
		Remaining:     e.Remaining,
	}
}
