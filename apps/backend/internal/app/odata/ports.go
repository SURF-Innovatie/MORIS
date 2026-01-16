package odata

import (
	"context"
	"net/url"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// BudgetODataDTO represents a budget in OData responses
type BudgetODataDTO struct {
	ID            uuid.UUID `json:"id"`
	ProjectID     uuid.UUID `json:"projectId"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	TotalAmount   float64   `json:"totalAmount"`
	TotalBudgeted float64   `json:"totalBudgeted"`
	TotalActuals  float64   `json:"totalActuals"`
	BurnRate      float64   `json:"burnRate"`
	Currency      string    `json:"currency"`
}

// LineItemODataDTO represents a line item in OData responses
type LineItemODataDTO struct {
	ID             uuid.UUID `json:"id"`
	BudgetID       uuid.UUID `json:"budgetId"`
	Category       string    `json:"category"`
	Description    string    `json:"description"`
	BudgetedAmount float64   `json:"budgetedAmount"`
	Year           int       `json:"year"`
	FundingSource  string    `json:"fundingSource"`
	TotalActuals   float64   `json:"totalActuals"`
}

// ActualODataDTO represents an actual in OData responses
type ActualODataDTO struct {
	ID           uuid.UUID `json:"id"`
	LineItemID   uuid.UUID `json:"lineItemId"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	RecordedDate string    `json:"recordedDate"`
	Source       string    `json:"source"`
}

// AnalyticsODataDTO represents analytics summary in OData responses
type AnalyticsODataDTO struct {
	ProjectID     uuid.UUID `json:"projectId"`
	ProjectTitle  string    `json:"projectTitle"`
	BudgetID      uuid.UUID `json:"budgetId"`
	TotalBudgeted float64   `json:"totalBudgeted"`
	TotalActuals  float64   `json:"totalActuals"`
	Remaining     float64   `json:"remaining"`
	BurnRate      float64   `json:"burnRate"`
	Status        string    `json:"status"`
}

// Repository defines the contract for OData data access
type Repository interface {
	QueryBudgets(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[BudgetODataDTO], error)
	QueryLineItems(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[LineItemODataDTO], error)
	QueryActuals(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[ActualODataDTO], error)
	QueryAnalytics(ctx context.Context, userID uuid.UUID, query entities.ODataQuery) (entities.ODataResult[AnalyticsODataDTO], error)
}

// QueryParser parses OData query strings
type QueryParser interface {
	Parse(queryParams url.Values) (entities.ODataQuery, error)
}
