package analytics

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Repository defines the contract for analytics data access
type Repository interface {
	GetOrgSummary(ctx context.Context, orgID uuid.UUID) (entities.OrgAnalyticsSummary, error)
	GetBurnRateData(ctx context.Context, orgID uuid.UUID, params entities.DateRangeParams) ([]entities.BurnRateDataPoint, error)
	GetCategoryBreakdown(ctx context.Context, orgID uuid.UUID) ([]entities.CategoryBreakdown, error)
	GetProjectHealth(ctx context.Context, orgID uuid.UUID) ([]entities.ProjectHealthSummary, error)
	GetFundingBreakdown(ctx context.Context, orgID uuid.UUID) ([]entities.FundingBreakdown, error)
}
