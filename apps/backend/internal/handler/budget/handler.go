package budget

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/budget"
	"github.com/SURF-Innovatie/MORIS/ent/budgetlineitem"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

// Handler handles budget-related HTTP requests
type Handler struct {
	client *ent.Client
}

// NewHandler creates a new budget handler
func NewHandler(client *ent.Client) *Handler {
	return &Handler{client: client}
}

// RegisterRoutes registers global budget routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/budgets/{budgetId}", func(r chi.Router) {
		r.Get("/", h.GetBudgetByID)
		r.Route("/line-items", func(r chi.Router) {
			r.Get("/", h.GetLineItems)
			r.Post("/", h.AddLineItem)
			r.Delete("/{lineItemId}", h.RemoveLineItem)
		})
		r.Route("/actuals", func(r chi.Router) {
			r.Get("/", h.GetActuals)
			r.Post("/", h.RecordActual)
		})
		r.Get("/analytics", h.GetAnalytics)
	})
}

// RegisterProjectRoutes registers project-scoped budget routes
func (h *Handler) RegisterProjectRoutes(r chi.Router) {
	r.Route("/{projectId}/budget", func(r chi.Router) {
		r.Get("/", h.GetBudget)
		r.Post("/", h.CreateBudget)
	})
}

// GetBudget godoc
// @Summary Get budget for a project
// @Description Retrieves the budget for a specific project
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID (UUID)"
// @Success 200 {object} dto.BudgetResponse
// @Failure 400 {string} string "invalid project id"
// @Failure 404 {string} string "budget not found"
// @Router /projects/{projectId}/budget [get]
func (h *Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "projectId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	b, err := h.client.Budget.
		Query().
		Where(budget.ProjectID(projectID)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		Only(r.Context())

	if err != nil {
		if ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "budget not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := mapBudgetToResponse(b)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetBudgetByID godoc
// @Summary Get budget by ID
// @Description Retrieves a specific budget by its ID
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Success 200 {object} dto.BudgetResponse
// @Failure 400 {string} string "invalid budget id"
// @Failure 404 {string} string "budget not found"
// @Router /budgets/{budgetId} [get]
func (h *Handler) GetBudgetByID(w http.ResponseWriter, r *http.Request) {
	budgetID, err := httputil.ParseUUIDParam(r, "budgetId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid budget id", nil)
		return
	}

	b, err := h.client.Budget.
		Query().
		Where(budget.ID(budgetID)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		Only(r.Context())

	if err != nil {
		if ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "budget not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := mapBudgetToResponse(b)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// CreateBudget godoc
// @Summary Create a budget for a project
// @Description Creates a new budget for the specified project
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID (UUID)"
// @Param body body dto.CreateBudgetRequest true "Budget creation request"
// @Success 201 {object} dto.BudgetResponse
// @Failure 400 {string} string "invalid request"
// @Failure 409 {string} string "budget already exists"
// @Router /projects/{projectId}/budget [post]
func (h *Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "projectId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	var req dto.CreateBudgetRequest
	if !httputil.ReadJSON(w, r, &req) {
		return // ReadJSON already wrote error
	}

	// Check if budget already exists
	exists, err := h.client.Budget.
		Query().
		Where(budget.ProjectID(projectID)).
		Exist(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if exists {
		httputil.WriteError(w, r, http.StatusConflict, "budget already exists for this project", nil)
		return
	}

	b, err := h.client.Budget.
		Create().
		SetProjectID(projectID).
		SetTitle(req.Title).
		SetDescription(req.Description).
		SetStatus(budget.StatusDraft).
		Save(r.Context())

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := mapBudgetToResponse(b)
	_ = httputil.WriteJSON(w, http.StatusCreated, resp)
}

// GetLineItems godoc
// @Summary Get line items for a budget
// @Description Retrieves all line items for a specific budget
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Success 200 {array} dto.BudgetLineItemResponse
// @Failure 400 {string} string "invalid budget id"
// @Router /budgets/{budgetId}/line-items [get]
func (h *Handler) GetLineItems(w http.ResponseWriter, r *http.Request) {
	budgetID, err := httputil.ParseUUIDParam(r, "budgetId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid budget id", nil)
		return
	}

	items, err := h.client.BudgetLineItem.
		Query().
		Where(budgetlineitem.BudgetID(budgetID)).
		WithActuals().
		All(r.Context())

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := mapLineItemsToResponse(items)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// AddLineItem godoc
// @Summary Add a line item to a budget
// @Description Adds a new line item to the specified budget
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Param body body dto.AddLineItemRequest true "Line item request"
// @Success 201 {object} dto.BudgetLineItemResponse
// @Failure 400 {string} string "invalid request"
// @Router /budgets/{budgetId}/line-items [post]
func (h *Handler) AddLineItem(w http.ResponseWriter, r *http.Request) {
	budgetID, err := httputil.ParseUUIDParam(r, "budgetId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid budget id", nil)
		return
	}

	var req dto.AddLineItemRequest
	if !httputil.ReadJSON(w, r, &req) {
		return // ReadJSON already wrote error
	}

	item, err := h.client.BudgetLineItem.
		Create().
		SetBudgetID(budgetID).
		SetCategory(budgetlineitem.Category(req.Category)).
		SetDescription(req.Description).
		SetBudgetedAmount(req.BudgetedAmount).
		SetYear(req.Year).
		SetFundingSource(budgetlineitem.FundingSource(req.FundingSource)).
		Save(r.Context())

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := mapLineItemToResponse(item)
	_ = httputil.WriteJSON(w, http.StatusCreated, resp)
}

// RemoveLineItem godoc
// @Summary Remove a line item from a budget
// @Description Removes a line item from the specified budget
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Param lineItemId path string true "Line Item ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "line item not found"
// @Router /budgets/{budgetId}/line-items/{lineItemId} [delete]
func (h *Handler) RemoveLineItem(w http.ResponseWriter, r *http.Request) {
	lineItemID, err := httputil.ParseUUIDParam(r, "lineItemId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid line item id", nil)
		return
	}

	err = h.client.BudgetLineItem.
		DeleteOneID(lineItemID).
		Exec(r.Context())

	if err != nil {
		if ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "line item not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetActuals godoc
// @Summary Get actuals for a budget
// @Description Retrieves all actual records for a budget
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Success 200 {array} dto.BudgetActualResponse
// @Failure 400 {string} string "invalid budget id"
// @Router /budgets/{budgetId}/actuals [get]
func (h *Handler) GetActuals(w http.ResponseWriter, r *http.Request) {
	budgetID, err := httputil.ParseUUIDParam(r, "budgetId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid budget id", nil)
		return
	}

	// Get all line items for this budget, then get their actuals
	items, err := h.client.BudgetLineItem.
		Query().
		Where(budgetlineitem.BudgetID(budgetID)).
		WithActuals().
		All(r.Context())

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var actuals []dto.BudgetActualResponse
	for _, item := range items {
		for _, actual := range item.Edges.Actuals {
			actuals = append(actuals, mapActualToResponse(actual))
		}
	}

	if actuals == nil {
		actuals = []dto.BudgetActualResponse{}
	}

	_ = httputil.WriteJSON(w, http.StatusOK, actuals)
}

// RecordActual godoc
// @Summary Record an actual expenditure
// @Description Records an actual expenditure against a budget line item
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Param body body dto.RecordActualRequest true "Actual record request"
// @Success 201 {object} dto.BudgetActualResponse
// @Failure 400 {string} string "invalid request"
// @Router /budgets/{budgetId}/actuals [post]
func (h *Handler) RecordActual(w http.ResponseWriter, r *http.Request) {
	var req dto.RecordActualRequest
	if !httputil.ReadJSON(w, r, &req) {
		return // ReadJSON already wrote error
	}

	builder := h.client.BudgetActual.
		Create().
		SetLineItemID(req.LineItemID).
		SetAmount(req.Amount).
		SetDescription(req.Description)

	if req.RecordedDate != nil {
		builder = builder.SetRecordedDate(*req.RecordedDate)
	}

	if req.Source != "" {
		builder = builder.SetSource(req.Source)
	}

	actual, err := builder.Save(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := mapActualToResponse(actual)
	_ = httputil.WriteJSON(w, http.StatusCreated, resp)
}

// GetAnalytics godoc
// @Summary Get budget analytics
// @Description Retrieves analytics data for a specific budget
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budgetId path string true "Budget ID (UUID)"
// @Success 200 {object} dto.BudgetAnalyticsResponse
// @Failure 400 {string} string "invalid budget id"
// @Failure 404 {string} string "budget not found"
// @Router /budgets/{budgetId}/analytics [get]
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	budgetID, err := httputil.ParseUUIDParam(r, "budgetId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid budget id", nil)
		return
	}

	b, err := h.client.Budget.
		Query().
		Where(budget.ID(budgetID)).
		WithLineItems(func(q *ent.BudgetLineItemQuery) {
			q.WithActuals()
		}).
		Only(r.Context())

	if err != nil {
		if ent.IsNotFound(err) {
			httputil.WriteError(w, r, http.StatusNotFound, "budget not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := computeAnalytics(b)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// Helper functions to map Ent entities to DTOs

func mapBudgetToResponse(b *ent.Budget) dto.BudgetResponse {
	var lineItems []dto.BudgetLineItemResponse
	var totalBudgeted, totalActuals float64

	for _, item := range b.Edges.LineItems {
		li := mapLineItemToResponse(item)
		lineItems = append(lineItems, li)
		totalBudgeted += item.BudgetedAmount
		for _, actual := range item.Edges.Actuals {
			totalActuals += actual.Amount
		}
	}

	if lineItems == nil {
		lineItems = []dto.BudgetLineItemResponse{}
	}

	remaining := totalBudgeted - totalActuals
	var burnRate float64
	if totalBudgeted > 0 {
		burnRate = (totalActuals / totalBudgeted) * 100
	}

	return dto.BudgetResponse{
		ID:            b.ID,
		ProjectID:     b.ProjectID,
		Title:         b.Title,
		Description:   b.Description,
		Status:        entities.BudgetStatus(b.Status),
		TotalAmount:   b.TotalAmount,
		Currency:      b.Currency,
		Version:       b.Version,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		LineItems:     lineItems,
		TotalBudgeted: totalBudgeted,
		TotalActuals:  totalActuals,
		Remaining:     remaining,
		BurnRate:      burnRate,
	}
}

func mapLineItemsToResponse(items []*ent.BudgetLineItem) []dto.BudgetLineItemResponse {
	resp := make([]dto.BudgetLineItemResponse, len(items))
	for i, item := range items {
		resp[i] = mapLineItemToResponse(item)
	}
	return resp
}

func mapLineItemToResponse(item *ent.BudgetLineItem) dto.BudgetLineItemResponse {
	var actuals []dto.BudgetActualResponse
	var totalActuals float64

	for _, actual := range item.Edges.Actuals {
		actuals = append(actuals, mapActualToResponse(actual))
		totalActuals += actual.Amount
	}

	if actuals == nil {
		actuals = []dto.BudgetActualResponse{}
	}

	return dto.BudgetLineItemResponse{
		ID:             item.ID,
		BudgetID:       item.BudgetID,
		Category:       entities.BudgetCategory(item.Category),
		Description:    item.Description,
		BudgetedAmount: item.BudgetedAmount,
		Year:           item.Year,
		FundingSource:  entities.FundingSource(item.FundingSource),
		Actuals:        actuals,
		TotalActuals:   totalActuals,
		Remaining:      item.BudgetedAmount - totalActuals,
	}
}

func mapActualToResponse(actual *ent.BudgetActual) dto.BudgetActualResponse {
	return dto.BudgetActualResponse{
		ID:           actual.ID,
		LineItemID:   actual.LineItemID,
		Amount:       actual.Amount,
		Description:  actual.Description,
		RecordedDate: actual.RecordedDate,
		Source:       actual.Source,
		ExternalRef:  actual.ExternalRef,
	}
}

func computeAnalytics(b *ent.Budget) dto.BudgetAnalyticsResponse {
	categoryMap := make(map[entities.BudgetCategory]*dto.CategoryBreakdown)
	yearMap := make(map[int]*dto.YearBreakdown)
	fundingMap := make(map[entities.FundingSource]*dto.FundingBreakdown)

	var totalBudgeted, totalActuals float64

	for _, item := range b.Edges.LineItems {
		category := entities.BudgetCategory(item.Category)
		funding := entities.FundingSource(item.FundingSource)

		var itemActuals float64
		for _, actual := range item.Edges.Actuals {
			itemActuals += actual.Amount
		}

		totalBudgeted += item.BudgetedAmount
		totalActuals += itemActuals

		// Category breakdown
		if _, ok := categoryMap[category]; !ok {
			categoryMap[category] = &dto.CategoryBreakdown{Category: category}
		}
		categoryMap[category].Budgeted += item.BudgetedAmount
		categoryMap[category].Actuals += itemActuals

		// Year breakdown
		if _, ok := yearMap[item.Year]; !ok {
			yearMap[item.Year] = &dto.YearBreakdown{Year: item.Year}
		}
		yearMap[item.Year].Budgeted += item.BudgetedAmount
		yearMap[item.Year].Actuals += itemActuals

		// Funding breakdown
		if _, ok := fundingMap[funding]; !ok {
			fundingMap[funding] = &dto.FundingBreakdown{FundingSource: funding}
		}
		fundingMap[funding].Budgeted += item.BudgetedAmount
		fundingMap[funding].Actuals += itemActuals
	}

	// Calculate remaining for each breakdown
	var byCategory []dto.CategoryBreakdown
	for _, cb := range categoryMap {
		cb.Remaining = cb.Budgeted - cb.Actuals
		byCategory = append(byCategory, *cb)
	}
	if byCategory == nil {
		byCategory = []dto.CategoryBreakdown{}
	}

	var byYear []dto.YearBreakdown
	for _, yb := range yearMap {
		yb.Remaining = yb.Budgeted - yb.Actuals
		byYear = append(byYear, *yb)
	}
	if byYear == nil {
		byYear = []dto.YearBreakdown{}
	}

	var byFunding []dto.FundingBreakdown
	for _, fb := range fundingMap {
		fb.Remaining = fb.Budgeted - fb.Actuals
		byFunding = append(byFunding, *fb)
	}
	if byFunding == nil {
		byFunding = []dto.FundingBreakdown{}
	}

	remaining := totalBudgeted - totalActuals
	var burnRate float64
	if totalBudgeted > 0 {
		burnRate = (totalActuals / totalBudgeted) * 100
	}

	return dto.BudgetAnalyticsResponse{
		BudgetID:      b.ID,
		ProjectID:     b.ProjectID,
		Title:         b.Title,
		Status:        entities.BudgetStatus(b.Status),
		TotalBudgeted: totalBudgeted,
		TotalActuals:  totalActuals,
		Remaining:     remaining,
		BurnRate:      burnRate,
		ByCategory:    byCategory,
		ByYear:        byYear,
		ByFunding:     byFunding,
		BurnRateData:  []dto.BurnRateDataPoint{}, // TODO: Compute time series
	}
}
