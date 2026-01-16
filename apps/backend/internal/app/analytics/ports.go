package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OrgAnalyticsSummary provides organization-level analytics
type OrgAnalyticsSummary struct {
	TotalProjects    int     `json:"totalProjects"`
	TotalBudgeted    float64 `json:"totalBudgeted"`
	TotalActuals     float64 `json:"totalActuals"`
	AverageSpendRate float64 `json:"averageSpendRate"`
	ProjectsAtRisk   int     `json:"projectsAtRisk"`
	ProjectsOnTrack  int     `json:"projectsOnTrack"`
}

// BurnRateDataPoint represents a point in burn rate time series
type BurnRateDataPoint struct {
	Date        time.Time `json:"date"`
	ProjectID   uuid.UUID `json:"projectId"`
	ProjectName string    `json:"projectName"`
	Budgeted    float64   `json:"budgeted"`
	Actual      float64   `json:"actual"`
	Projected   float64   `json:"projected"`
}

// CategoryBreakdown shows spending by category
type CategoryBreakdown struct {
	Category  string  `json:"category"`
	Budgeted  float64 `json:"budgeted"`
	Actuals   float64 `json:"actuals"`
	Remaining float64 `json:"remaining"`
}

// ProjectHealthSummary shows health indicators for a project
type ProjectHealthSummary struct {
	ProjectID   uuid.UUID `json:"projectId"`
	ProjectName string    `json:"projectName"`
	BudgetID    uuid.UUID `json:"budgetId"`
	Budgeted    float64   `json:"budgeted"`
	Spent       float64   `json:"spent"`
	Remaining   float64   `json:"remaining"`
	BurnRate    float64   `json:"burnRate"`
	Status      string    `json:"status"` // "on_track", "warning", "at_risk"
}

// FundingBreakdown shows spending by funding source
type FundingBreakdown struct {
	FundingSource string  `json:"fundingSource"`
	Budgeted      float64 `json:"budgeted"`
	Actuals       float64 `json:"actuals"`
	Remaining     float64 `json:"remaining"`
}

// DateRangeParams defines date filtering for analytics
type DateRangeParams struct {
	StartDate *time.Time
	EndDate   *time.Time
}

// Repository defines the contract for analytics data access
type Repository interface {
	GetOrgSummary(ctx context.Context, orgID uuid.UUID) (OrgAnalyticsSummary, error)
	GetBurnRateData(ctx context.Context, orgID uuid.UUID, params DateRangeParams) ([]BurnRateDataPoint, error)
	GetCategoryBreakdown(ctx context.Context, orgID uuid.UUID) ([]CategoryBreakdown, error)
	GetProjectHealth(ctx context.Context, orgID uuid.UUID) ([]ProjectHealthSummary, error)
	GetFundingBreakdown(ctx context.Context, orgID uuid.UUID) ([]FundingBreakdown, error)
}
