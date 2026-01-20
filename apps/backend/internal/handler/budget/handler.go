package budget

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/budget"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

// Handler handles budget-related HTTP requests
type Handler struct {
	service *budget.Service
}

// NewHandler creates a new budget handler
func NewHandler(service *budget.Service) *Handler {
	return &Handler{service: service}
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

	b, err := h.service.GetBudget(r.Context(), projectID)
	if err != nil {
		if err == budget.ErrBudgetNotFound {
			httputil.WriteError(w, r, http.StatusNotFound, "budget not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOItem[dto.BudgetResponse](*b)
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

	b, err := h.service.GetBudgetByID(r.Context(), budgetID)
	if err != nil {
		if err == budget.ErrBudgetNotFound {
			httputil.WriteError(w, r, http.StatusNotFound, "budget not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOItem[dto.BudgetResponse](*b)
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

	b, err := h.service.CreateBudget(r.Context(), projectID, req.Title, req.Description)
	if err != nil {
		if err == budget.ErrBudgetAlreadyExists {
			httputil.WriteError(w, r, http.StatusConflict, err.Error(), nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOItem[dto.BudgetResponse](*b)
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

	items, err := h.service.GetLineItems(r.Context(), budgetID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOs[dto.BudgetLineItemResponse](items)
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

	item := entities.BudgetLineItem{
		Category:       req.Category,
		Description:    req.Description,
		BudgetedAmount: req.BudgetedAmount,
		Year:           req.Year,
		FundingSource:  req.FundingSource,
		NWOGrantID:     req.NWOGrantID,
	}

	created, err := h.service.AddLineItem(r.Context(), budgetID, item)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOItem[dto.BudgetLineItemResponse](*created)
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

	err = h.service.RemoveLineItem(r.Context(), lineItemID)
	if err != nil {
		if err == budget.ErrLineItemNotFound {
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

	actuals, err := h.service.GetActuals(r.Context(), budgetID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOs[dto.BudgetActualResponse](actuals)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
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

	actual := entities.BudgetActual{
		LineItemID:  req.LineItemID,
		Amount:      req.Amount,
		Description: req.Description,
		Source:      req.Source,
	}
	if req.RecordedDate != nil {
		actual.RecordedDate = *req.RecordedDate
	}

	created, err := h.service.RecordActual(r.Context(), actual)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOItem[dto.BudgetActualResponse](*created)
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

	analytics, err := h.service.GetAnalytics(r.Context(), budgetID)
	if err != nil {
		if err == budget.ErrBudgetNotFound {
			httputil.WriteError(w, r, http.StatusNotFound, "budget not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := transform.ToDTOItem[dto.BudgetAnalyticsResponse](*analytics)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}
