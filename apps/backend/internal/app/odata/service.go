package odata

import (
	"context"
	"net/url"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// Service provides OData query functionality
type Service struct {
	repo   Repository
	parser QueryParser
}

// NewService creates a new OData service
func NewService(repo Repository, parser QueryParser) *Service {
	return &Service{
		repo:   repo,
		parser: parser,
	}
}

// GetBudgets returns budgets the user has access to
func (s *Service) GetBudgets(ctx context.Context, userID uuid.UUID, queryParams url.Values) (entities.ODataResult[BudgetODataDTO], error) {
	query, err := s.parser.Parse(queryParams)
	if err != nil {
		return entities.ODataResult[BudgetODataDTO]{}, err
	}

	// Validate fields against allowed list
	if err := validateQuery(query, entities.BudgetODataFields); err != nil {
		return entities.ODataResult[BudgetODataDTO]{}, err
	}

	return s.repo.QueryBudgets(ctx, userID, query)
}

// GetLineItems returns budget line items the user has access to
func (s *Service) GetLineItems(ctx context.Context, userID uuid.UUID, queryParams url.Values) (entities.ODataResult[LineItemODataDTO], error) {
	query, err := s.parser.Parse(queryParams)
	if err != nil {
		return entities.ODataResult[LineItemODataDTO]{}, err
	}

	if err := validateQuery(query, entities.BudgetLineItemODataFields); err != nil {
		return entities.ODataResult[LineItemODataDTO]{}, err
	}

	return s.repo.QueryLineItems(ctx, userID, query)
}

// GetActuals returns budget actuals the user has access to
func (s *Service) GetActuals(ctx context.Context, userID uuid.UUID, queryParams url.Values) (entities.ODataResult[ActualODataDTO], error) {
	query, err := s.parser.Parse(queryParams)
	if err != nil {
		return entities.ODataResult[ActualODataDTO]{}, err
	}

	if err := validateQuery(query, entities.BudgetActualODataFields); err != nil {
		return entities.ODataResult[ActualODataDTO]{}, err
	}

	return s.repo.QueryActuals(ctx, userID, query)
}

// GetAnalytics returns aggregated analytics the user has access to
func (s *Service) GetAnalytics(ctx context.Context, userID uuid.UUID, queryParams url.Values) (entities.ODataResult[AnalyticsODataDTO], error) {
	query, err := s.parser.Parse(queryParams)
	if err != nil {
		return entities.ODataResult[AnalyticsODataDTO]{}, err
	}

	return s.repo.QueryAnalytics(ctx, userID, query)
}

// validateQuery checks that all fields in the query are allowed
func validateQuery(query entities.ODataQuery, allowed entities.AllowedFields) error {
	for _, field := range query.Select {
		if !allowed.IsSelectAllowed(field) {
			return &InvalidFieldError{Field: field, Reason: "not selectable"}
		}
	}

	if query.Filter != nil {
		if err := validateFilter(query.Filter, allowed); err != nil {
			return err
		}
	}

	for _, order := range query.OrderBy {
		if !allowed.IsSortAllowed(order.Field) {
			return &InvalidFieldError{Field: order.Field, Reason: "not sortable"}
		}
	}

	for _, expand := range query.Expand {
		if !allowed.IsExpandAllowed(expand) {
			return &InvalidFieldError{Field: expand, Reason: "not expandable"}
		}
	}

	return nil
}

func validateFilter(filter *entities.ODataFilter, allowed entities.AllowedFields) error {
	if filter == nil {
		return nil
	}

	if filter.Field != "" && !allowed.IsFilterAllowed(filter.Field) {
		return &InvalidFieldError{Field: filter.Field, Reason: "not filterable"}
	}

	if filter.And != nil {
		if err := validateFilter(filter.And, allowed); err != nil {
			return err
		}
	}

	if filter.Or != nil {
		if err := validateFilter(filter.Or, allowed); err != nil {
			return err
		}
	}

	return nil
}

// InvalidFieldError indicates a field is not allowed in the query
type InvalidFieldError struct {
	Field  string
	Reason string
}

func (e *InvalidFieldError) Error() string {
	return "invalid field '" + e.Field + "': " + e.Reason
}
