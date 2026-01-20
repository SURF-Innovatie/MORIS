package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
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

func (d BudgetResponse) FromEntity(e entities.Budget) BudgetResponse {
	summary := e.CalculateSummary()
	return BudgetResponse{
		ID:            e.ID,
		ProjectID:     e.ProjectID,
		Title:         e.Title,
		Description:   e.Description,
		Status:        e.Status,
		TotalAmount:   e.TotalAmount,
		Currency:      e.Currency,
		Version:       e.Version,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		LineItems:     transform.ToDTOs[BudgetLineItemResponse](e.LineItems),
		TotalBudgeted: summary.TotalBudgeted,
		TotalActuals:  summary.TotalActuals,
		Remaining:     summary.Remaining,
		BurnRate:      summary.BurnRate,
	}
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
	NWOGrantID     *string                 `json:"nwoGrantId,omitempty"` // The ID of the linked NWO grant

	// Computed fields
	TotalActuals float64 `json:"totalActuals"`
	Remaining    float64 `json:"remaining"`
}

func (d BudgetLineItemResponse) FromEntity(e entities.BudgetLineItem) BudgetLineItemResponse {
	totalActuals := e.CalculateActualTotal()
	return BudgetLineItemResponse{
		ID:             e.ID,
		BudgetID:       e.BudgetID,
		Category:       e.Category,
		Description:    e.Description,
		BudgetedAmount: e.BudgetedAmount,
		Year:           e.Year,
		FundingSource:  e.FundingSource,
		Actuals:        transform.ToDTOs[BudgetActualResponse](e.Actuals),
		TotalActuals:   totalActuals,
		NWOGrantID:     e.NWOGrantID,
		Remaining:      e.BudgetedAmount - totalActuals,
	}
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

func (d BudgetActualResponse) FromEntity(e entities.BudgetActual) BudgetActualResponse {
	return BudgetActualResponse{
		ID:           e.ID,
		LineItemID:   e.LineItemID,
		Amount:       e.Amount,
		Description:  e.Description,
		RecordedDate: e.RecordedDate,
		Source:       e.Source,
		ExternalRef:  e.ExternalRef,
	}
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
	NWOGrantID     *string                 `json:"nwoGrantId,omitempty"`
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
	BudgetID      uuid.UUID                    `json:"budgetId"`
	ProjectID     uuid.UUID                    `json:"projectId"`
	Title         string                       `json:"title"`
	Status        entities.BudgetStatus        `json:"status"`
	TotalBudgeted float64                      `json:"totalBudgeted"`
	TotalActuals  float64                      `json:"totalActuals"`
	Remaining     float64                      `json:"remaining"`
	BurnRate      float64                      `json:"burnRate"`
	ByCategory    []CategoryBreakdown          `json:"byCategory"`
	ByYear        []YearBreakdown              `json:"byYear"`
	ByFunding     []FundingBreakdown           `json:"byFunding"`
	BurnRateData  []entities.BurnRateDataPoint `json:"burnRateData"` // Using entities directly for now or transform if needed
}

func (d BudgetAnalyticsResponse) FromEntity(e entities.BudgetAnalytics) BudgetAnalyticsResponse {
	var byCategory []CategoryBreakdown
	for _, v := range e.CategoryMap {
		byCategory = append(byCategory, CategoryBreakdown{
			Category:  entities.BudgetCategory(v.Category),
			Budgeted:  v.Budgeted,
			Actuals:   v.Actuals,
			Remaining: v.Remaining,
		})
	}

	var byYear []YearBreakdown
	for _, v := range e.YearMap {
		byYear = append(byYear, YearBreakdown{
			Year:      v.Year,
			Budgeted:  v.Budgeted,
			Actuals:   v.Actuals,
			Remaining: v.Remaining,
		})
	}

	var byFunding []FundingBreakdown
	for _, v := range e.FundingMap {
		byFunding = append(byFunding, FundingBreakdown{
			FundingSource: entities.FundingSource(v.FundingSource),
			Budgeted:      v.Budgeted,
			Actuals:       v.Actuals,
			Remaining:     v.Remaining,
		})
	}

	return BudgetAnalyticsResponse{
		BudgetID:      e.BudgetID,
		ProjectID:     e.ProjectID,
		Title:         e.Title,
		Status:        e.Status,
		TotalBudgeted: e.TotalBudgeted,
		TotalActuals:  e.TotalActuals,
		Remaining:     e.Remaining,
		BurnRate:      e.BurnRate,
		ByCategory:    byCategory,
		ByYear:        byYear,
		ByFunding:     byFunding,
		BurnRateData:  []entities.BurnRateDataPoint{}, // TODO: Populate if available
	}
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
