package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrgAnalyticsSummary provides organization-level analytics
type OrgAnalyticsSummary struct {
	TotalProjects    int
	TotalBudgeted    float64
	TotalActuals     float64
	AverageSpendRate float64
	ProjectsAtRisk   int
	ProjectsOnTrack  int
}

// BurnRateDataPoint represents a point in burn rate time series
type BurnRateDataPoint struct {
	Date        time.Time
	ProjectID   uuid.UUID
	ProjectName string
	Budgeted    float64
	Actual      float64
	Projected   float64
}

// CategoryBreakdown shows spending by category
type CategoryBreakdown struct {
	Category  string
	Budgeted  float64
	Actuals   float64
	Remaining float64
}

// ProjectHealthSummary shows health indicators for a project
type ProjectHealthSummary struct {
	ProjectID   uuid.UUID
	ProjectName string
	BudgetID    uuid.UUID
	Budgeted    float64
	Spent       float64
	Remaining   float64
	BurnRate    float64
	Status      string // "on_track", "warning", "at_risk"
}

// FundingBreakdown shows spending by funding source
type FundingBreakdown struct {
	FundingSource string
	Budgeted      float64
	Actuals       float64
	Remaining     float64
}

// DateRangeParams defines date filtering for analytics
type DateRangeParams struct {
	StartDate *time.Time
	EndDate   *time.Time
}
