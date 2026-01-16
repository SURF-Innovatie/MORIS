package analytics

import (
	"context"

	"github.com/google/uuid"
)

// Service provides organization-level analytics
type Service struct {
	repo Repository
}

// NewService creates a new analytics service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetOrgSummary returns overall statistics for an organization
func (s *Service) GetOrgSummary(ctx context.Context, orgID uuid.UUID) (OrgAnalyticsSummary, error) {
	return s.repo.GetOrgSummary(ctx, orgID)
}

// GetBurnRateData returns time-series burn rate data
func (s *Service) GetBurnRateData(ctx context.Context, orgID uuid.UUID, params DateRangeParams) ([]BurnRateDataPoint, error) {
	return s.repo.GetBurnRateData(ctx, orgID, params)
}

// GetCategoryBreakdown returns spending grouped by category
func (s *Service) GetCategoryBreakdown(ctx context.Context, orgID uuid.UUID) ([]CategoryBreakdown, error) {
	return s.repo.GetCategoryBreakdown(ctx, orgID)
}

// GetProjectHealth returns health summary for all projects
func (s *Service) GetProjectHealth(ctx context.Context, orgID uuid.UUID) ([]ProjectHealthSummary, error) {
	return s.repo.GetProjectHealth(ctx, orgID)
}

// GetFundingBreakdown returns spending grouped by funding source
func (s *Service) GetFundingBreakdown(ctx context.Context, orgID uuid.UUID) ([]FundingBreakdown, error) {
	return s.repo.GetFundingBreakdown(ctx, orgID)
}
