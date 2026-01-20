package entities

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/google/uuid"
)

// BudgetStatus represents the lifecycle state of a budget
type BudgetStatus string

const (
	BudgetStatusDraft     BudgetStatus = "draft"
	BudgetStatusSubmitted BudgetStatus = "submitted"
	BudgetStatusApproved  BudgetStatus = "approved"
	BudgetStatusLocked    BudgetStatus = "locked"
)

// BudgetCategory represents the type of budget line item
type BudgetCategory string

const (
	BudgetCategoryPersonnel  BudgetCategory = "personnel"
	BudgetCategoryMaterial   BudgetCategory = "material"
	BudgetCategoryInvestment BudgetCategory = "investment"
	BudgetCategoryTravel     BudgetCategory = "travel"
	BudgetCategoryManagement BudgetCategory = "management"
	BudgetCategoryGrant      BudgetCategory = "grant"
	BudgetCategoryOther      BudgetCategory = "other"
)

// FundingSource represents the source of funding for a budget line item
type FundingSource string

const (
	FundingSourceSubsidy           FundingSource = "subsidy"
	FundingSourceCofinancingCash   FundingSource = "cofinancing_cash"
	FundingSourceCofinancingInkind FundingSource = "cofinancing_inkind"
)

// Budget represents a project budget aggregate
type Budget struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Title       string
	Description string
	Status      BudgetStatus
	TotalAmount float64
	Currency    string
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LineItems   []BudgetLineItem
}

func (b *Budget) FromEnt(e *ent.Budget) *Budget {
	b.ID = e.ID
	b.ProjectID = e.ProjectID
	b.Title = e.Title
	b.Description = e.Description
	b.Status = BudgetStatus(e.Status)
	b.TotalAmount = e.TotalAmount
	b.Currency = e.Currency
	b.Version = e.Version
	b.CreatedAt = e.CreatedAt
	b.UpdatedAt = e.UpdatedAt

	for _, li := range e.Edges.LineItems {
		b.LineItems = append(b.LineItems, *(&BudgetLineItem{}).FromEnt(li))
	}

	return b
}

// BudgetLineItem represents a single line in a budget
type BudgetLineItem struct {
	ID             uuid.UUID
	BudgetID       uuid.UUID
	Category       BudgetCategory
	Description    string
	BudgetedAmount float64
	Year           int
	FundingSource  FundingSource
	NWOGrantID     *string
	Actuals        []BudgetActual
}

func (bli *BudgetLineItem) FromEnt(e *ent.BudgetLineItem) *BudgetLineItem {
	bli.ID = e.ID
	bli.BudgetID = e.BudgetID
	bli.Category = BudgetCategory(e.Category)
	bli.Description = e.Description
	bli.BudgetedAmount = e.BudgetedAmount
	bli.Year = e.Year
	bli.FundingSource = FundingSource(e.FundingSource)

	if e.NwoGrantID != nil {
		bli.NWOGrantID = e.NwoGrantID
	}

	for _, a := range e.Edges.Actuals {
		bli.Actuals = append(bli.Actuals, *(&BudgetActual{}).FromEnt(a))
	}

	return bli
}

// BudgetActual represents an actual expenditure recorded against a line item
type BudgetActual struct {
	ID           uuid.UUID
	LineItemID   uuid.UUID
	Amount       float64
	Description  string
	RecordedDate time.Time
	Source       string // "manual" | "erp_sync"
	ExternalRef  string // For ERP reconciliation
}

func (ba *BudgetActual) FromEnt(e *ent.BudgetActual) *BudgetActual {
	ba.ID = e.ID
	ba.LineItemID = e.LineItemID
	ba.Amount = e.Amount
	ba.Description = e.Description
	ba.RecordedDate = e.RecordedDate
	ba.Source = e.Source
	ba.ExternalRef = e.ExternalRef

	return ba
}

// BudgetSummary provides computed analytics for a budget
type BudgetSummary struct {
	TotalBudgeted float64
	TotalActuals  float64
	Remaining     float64
	BurnRate      float64 // Percentage of budget consumed
}

// YearBreakdown shows spending by year
type YearBreakdown struct {
	Year      int
	Budgeted  float64
	Actuals   float64
	Remaining float64
}

// BudgetAnalytics provides detailed analytics for a budget
type BudgetAnalytics struct {
	BudgetID      uuid.UUID
	ProjectID     uuid.UUID
	Title         string
	Status        BudgetStatus
	CalculatedAt  time.Time
	TotalBudgeted float64
	TotalActuals  float64
	Remaining     float64
	BurnRate      float64
	CategoryMap   map[BudgetCategory]CategoryBreakdown
	YearMap       map[int]YearBreakdown
	FundingMap    map[FundingSource]FundingBreakdown
}

// CalculateSummary computes the budget summary from line items
func (b *Budget) CalculateSummary() BudgetSummary {
	var totalBudgeted, totalActuals float64

	for _, item := range b.LineItems {
		totalBudgeted += item.BudgetedAmount
		for _, actual := range item.Actuals {
			totalActuals += actual.Amount
		}
	}

	remaining := totalBudgeted - totalActuals
	var burnRate float64
	if totalBudgeted > 0 {
		burnRate = (totalActuals / totalBudgeted) * 100
	}

	return BudgetSummary{
		TotalBudgeted: totalBudgeted,
		TotalActuals:  totalActuals,
		Remaining:     remaining,
		BurnRate:      burnRate,
	}
}

// GetLineItemByID returns a line item by ID
func (b *Budget) GetLineItemByID(id uuid.UUID) *BudgetLineItem {
	for i := range b.LineItems {
		if b.LineItems[i].ID == id {
			return &b.LineItems[i]
		}
	}
	return nil
}

// CalculateActualTotal returns the total actuals for a line item
func (li *BudgetLineItem) CalculateActualTotal() float64 {
	var total float64
	for _, actual := range li.Actuals {
		total += actual.Amount
	}
	return total
}

// RemainingAmount returns the remaining budget for a line item
func (li *BudgetLineItem) RemainingAmount() float64 {
	return li.BudgetedAmount - li.CalculateActualTotal()
}
