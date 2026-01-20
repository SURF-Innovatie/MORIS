package budget

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

var (
	ErrBudgetNotFound      = errors.New("budget not found")
	ErrBudgetAlreadyExists = errors.New("budget already exists for this project")
	ErrLineItemNotFound    = errors.New("line item not found")
)

// Service provides budget management operations
type Service struct {
	repo Repository
}

// NewService creates a new budget service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBudget(ctx context.Context, projectID uuid.UUID) (*entities.Budget, error) {
	b, err := s.repo.GetBudget(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Budget](b), nil
}

func (s *Service) GetBudgetByID(ctx context.Context, budgetID uuid.UUID) (*entities.Budget, error) {
	b, err := s.repo.GetBudgetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Budget](b), nil
}

func (s *Service) CreateBudget(ctx context.Context, projectID uuid.UUID, title, description string) (*entities.Budget, error) {
	exists, err := s.repo.BudgetExists(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrBudgetAlreadyExists
	}
	b, err := s.repo.CreateBudget(ctx, projectID, title, description)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.Budget](b), nil
}

func (s *Service) GetLineItems(ctx context.Context, budgetID uuid.UUID) ([]entities.BudgetLineItem, error) {
	items, err := s.repo.GetLineItems(ctx, budgetID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntities[entities.BudgetLineItem](items), nil
}

func (s *Service) AddLineItem(ctx context.Context, budgetID uuid.UUID, item entities.BudgetLineItem) (*entities.BudgetLineItem, error) {
	li, err := s.repo.AddLineItem(ctx, budgetID, item)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.BudgetLineItem](li), nil
}

func (s *Service) RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	return s.repo.RemoveLineItem(ctx, lineItemID)
}

func (s *Service) GetActuals(ctx context.Context, budgetID uuid.UUID) ([]entities.BudgetActual, error) {
	actuals, err := s.repo.GetActuals(ctx, budgetID)
	if err != nil {
		return nil, err
	}
	return transform.ToEntities[entities.BudgetActual](actuals), nil
}

func (s *Service) RecordActual(ctx context.Context, actual entities.BudgetActual) (*entities.BudgetActual, error) {
	a, err := s.repo.RecordActual(ctx, actual)
	if err != nil {
		return nil, err
	}
	return transform.ToEntityPtr[entities.BudgetActual](a), nil
}

// GetAnalytics computes analytics for a budget
func (s *Service) GetAnalytics(ctx context.Context, budgetID uuid.UUID) (*entities.BudgetAnalytics, error) {
	bEnt, err := s.repo.GetBudgetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}
	b := transform.ToEntityPtr[entities.Budget](bEnt)

	analytics := &entities.BudgetAnalytics{
		BudgetID:     b.ID,
		ProjectID:    b.ProjectID,
		Title:        b.Title,
		Status:       b.Status,
		CalculatedAt: b.UpdatedAt, // Approximation
		CategoryMap:  make(map[entities.BudgetCategory]entities.CategoryBreakdown),
		YearMap:      make(map[int]entities.YearBreakdown),
		FundingMap:   make(map[entities.FundingSource]entities.FundingBreakdown),
	}

	var totalBudgeted, totalActuals float64

	for _, item := range b.LineItems {
		var itemActuals float64
		for _, actual := range item.Actuals {
			itemActuals += actual.Amount
		}

		totalBudgeted += item.BudgetedAmount
		totalActuals += itemActuals

		// Category breakdown
		catStats := analytics.CategoryMap[item.Category]
		catStats.Category = string(item.Category) // Ensure name is set
		catStats.Budgeted += item.BudgetedAmount
		catStats.Actuals += itemActuals
		catStats.Remaining = catStats.Budgeted - catStats.Actuals
		analytics.CategoryMap[item.Category] = catStats

		// Year breakdown
		yearStats := analytics.YearMap[item.Year]
		yearStats.Year = item.Year
		yearStats.Budgeted += item.BudgetedAmount
		yearStats.Actuals += itemActuals
		yearStats.Remaining = yearStats.Budgeted - yearStats.Actuals
		analytics.YearMap[item.Year] = yearStats

		// Funding breakdown
		fundStats := analytics.FundingMap[item.FundingSource]
		fundStats.FundingSource = string(item.FundingSource)
		fundStats.Budgeted += item.BudgetedAmount
		fundStats.Actuals += itemActuals
		fundStats.Remaining = fundStats.Budgeted - fundStats.Actuals
		analytics.FundingMap[item.FundingSource] = fundStats
	}

	analytics.TotalBudgeted = totalBudgeted
	analytics.TotalActuals = totalActuals
	analytics.Remaining = totalBudgeted - totalActuals
	if totalBudgeted > 0 {
		analytics.BurnRate = (totalActuals / totalBudgeted) * 100
	}

	return analytics, nil
}
