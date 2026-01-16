package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

// BudgetResponse represents a full budget with line items and computed fields
type BudgetResponse struct {
	ID          uuid.UUID                `json:"id"`
	ProjectID   uuid.UUID                `json:"projectId"`
	Title       string                   `json:"title"`
	Description string                   `json:"description,omitempty"`
	Status      entities.BudgetStatus    `json:"status"`
	TotalAmount float64                  `json:"totalAmount"`
	Currency    string                   `json:"currency"`
	Version     int                      `json:"version"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
	LineItems   []BudgetLineItemResponse `json:"lineItems,omitempty"`

	// Computed fields
	TotalBudgeted float64 `json:"totalBudgeted"`
	TotalActuals  float64 `json:"totalActuals"`
	Remaining     float64 `json:"remaining"`
	BurnRate      float64 `json:"burnRate"` // Percentage
}

// BudgetSummaryResponse is a lightweight list view of budgets
type BudgetSummaryResponse struct {
	ID            uuid.UUID             `json:"id"`
	ProjectID     uuid.UUID             `json:"projectId"`
	Title         string                `json:"title"`
	Status        entities.BudgetStatus `json:"status"`
	TotalBudgeted float64               `json:"totalBudgeted"`
	TotalActuals  float64               `json:"totalActuals"`
	BurnRate      float64               `json:"burnRate"`
}

// BudgetLineItemResponse represents a single budget line item
type BudgetLineItemResponse struct {
	ID             uuid.UUID               `json:"id"`
	BudgetID       uuid.UUID               `json:"budgetId"`
	Category       entities.BudgetCategory `json:"category"`
	Description    string                  `json:"description"`
	BudgetedAmount float64                 `json:"budgetedAmount"`
	Year           int                     `json:"year"`
	FundingSource  entities.FundingSource  `json:"fundingSource"`
	Actuals        []BudgetActualResponse  `json:"actuals,omitempty"`

	// Computed fields
	TotalActuals float64 `json:"totalActuals"`
	Remaining    float64 `json:"remaining"`
}

// BudgetActualResponse represents an actual expenditure
type BudgetActualResponse struct {
	ID           uuid.UUID `json:"id"`
	LineItemID   uuid.UUID `json:"lineItemId"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description,omitempty"`
	RecordedDate time.Time `json:"recordedDate"`
	Source       string    `json:"source"`
	ExternalRef  string    `json:"externalRef,omitempty"`
}

// CreateBudgetRequest is the input for creating a budget
type CreateBudgetRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// AddLineItemRequest is the input for adding a budget line item
type AddLineItemRequest struct {
	Category       entities.BudgetCategory `json:"category" binding:"required"`
	Description    string                  `json:"description" binding:"required"`
	BudgetedAmount float64                 `json:"budgetedAmount" binding:"required,gt=0"`
	Year           int                     `json:"year" binding:"required,gte=2000,lte=2100"`
	FundingSource  entities.FundingSource  `json:"fundingSource" binding:"required"`
}

// RecordActualRequest is the input for recording an actual expenditure
type RecordActualRequest struct {
	LineItemID   uuid.UUID  `json:"lineItemId" binding:"required"`
	Amount       float64    `json:"amount" binding:"required,gt=0"`
	Description  string     `json:"description"`
	RecordedDate *time.Time `json:"recordedDate"`
	Source       string     `json:"source"` // defaults to "manual"
}

// BudgetAnalyticsResponse provides analytics data for a budget
type BudgetAnalyticsResponse struct {
	BudgetID      uuid.UUID             `json:"budgetId"`
	ProjectID     uuid.UUID             `json:"projectId"`
	Title         string                `json:"title"`
	Status        entities.BudgetStatus `json:"status"`
	TotalBudgeted float64               `json:"totalBudgeted"`
	TotalActuals  float64               `json:"totalActuals"`
	Remaining     float64               `json:"remaining"`
	BurnRate      float64               `json:"burnRate"`
	ByCategory    []CategoryBreakdown   `json:"byCategory"`
	ByYear        []YearBreakdown       `json:"byYear"`
	ByFunding     []FundingBreakdown    `json:"byFunding"`
	BurnRateData  []BurnRateDataPoint   `json:"burnRateData"`
}

// CategoryBreakdown shows budgeted vs actuals by category
type CategoryBreakdown struct {
	Category  entities.BudgetCategory `json:"category"`
	Budgeted  float64                 `json:"budgeted"`
	Actuals   float64                 `json:"actuals"`
	Remaining float64                 `json:"remaining"`
}

// YearBreakdown shows budgeted vs actuals by year
type YearBreakdown struct {
	Year      int     `json:"year"`
	Budgeted  float64 `json:"budgeted"`
	Actuals   float64 `json:"actuals"`
	Remaining float64 `json:"remaining"`
}

// FundingBreakdown shows budgeted vs actuals by funding source
type FundingBreakdown struct {
	FundingSource entities.FundingSource `json:"fundingSource"`
	Budgeted      float64                `json:"budgeted"`
	Actuals       float64                `json:"actuals"`
	Remaining     float64                `json:"remaining"`
}

// BurnRateDataPoint represents a point in the burn rate time series
type BurnRateDataPoint struct {
	Date             time.Time `json:"date"`
	CumulativeActual float64   `json:"cumulativeActual"`
	IdealBurn        float64   `json:"idealBurn"`
	ProjectedEnd     float64   `json:"projectedEnd,omitempty"`
}
