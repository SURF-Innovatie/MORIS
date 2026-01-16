package odata

import (
	"net/http"

	appOdata "github.com/SURF-Innovatie/MORIS/internal/app/odata"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles OData HTTP requests for Power BI integration
type Handler struct {
	service *appOdata.Service
}

// NewHandler creates a new OData handler
func NewHandler(service *appOdata.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers OData routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/odata", func(r chi.Router) {
		// These routes support both JWT and API key authentication
		r.Get("/budgets", h.GetBudgets)
		r.Get("/budget-line-items", h.GetBudgetLineItems)
		r.Get("/budget-actuals", h.GetBudgetActuals)
		r.Get("/budget-analytics", h.GetBudgetAnalytics)
	})
}

// GetBudgets godoc
// @Summary Get budgets via OData
// @Description OData endpoint for Power BI - returns budgets with $select, $filter, $orderby, $top, $skip support
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param $select query string false "Fields to select (comma-separated)"
// @Param $filter query string false "OData filter expression"
// @Param $orderby query string false "Order by field (e.g., 'title desc')"
// @Param $top query integer false "Maximum number of results"
// @Param $skip query integer false "Number of results to skip"
// @Param $count query boolean false "Include total count in response"
// @Success 200 {object} map[string]any "OData response with @odata.count and value array"
// @Failure 400 {object} map[string]any "OData error response"
// @Failure 401 {string} string "unauthorized"
// @Router /odata/budgets [get]
func (h *Handler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetBudgets(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponse(w, result)
}

// GetBudgetLineItems godoc
// @Summary Get budget line items via OData
// @Description OData endpoint for Power BI - returns budget line items
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param $select query string false "Fields to select"
// @Param $filter query string false "OData filter expression"
// @Param $orderby query string false "Order by field"
// @Param $top query integer false "Maximum number of results"
// @Param $skip query integer false "Number of results to skip"
// @Success 200 {object} map[string]any "OData response"
// @Failure 400 {object} map[string]any "OData error response"
// @Router /odata/budget-line-items [get]
func (h *Handler) GetBudgetLineItems(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetLineItems(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponse(w, result)
}

// GetBudgetActuals godoc
// @Summary Get budget actuals via OData
// @Description OData endpoint for Power BI - returns recorded actuals
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param $select query string false "Fields to select"
// @Param $filter query string false "OData filter expression"
// @Param $top query integer false "Maximum number of results"
// @Param $skip query integer false "Number of results to skip"
// @Success 200 {object} map[string]any "OData response"
// @Failure 400 {object} map[string]any "OData error response"
// @Router /odata/budget-actuals [get]
func (h *Handler) GetBudgetActuals(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetActuals(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponse(w, result)
}

// GetBudgetAnalytics godoc
// @Summary Get budget analytics via OData
// @Description OData endpoint for Power BI - returns pre-computed analytics summary per project
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any "OData response with analytics data"
// @Failure 401 {string} string "unauthorized"
// @Router /odata/budget-analytics [get]
func (h *Handler) GetBudgetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetAnalytics(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponse(w, result)
}

// Helper functions

func getUserIDFromContext(r *http.Request) uuid.UUID {
	// TODO: Extract from JWT or API key middleware
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		return uuid.Nil
	}
	if id, ok := userIDVal.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

func writeODataResponse(w http.ResponseWriter, result any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("OData-Version", "4.0")
	_ = httputil.WriteJSON(w, http.StatusOK, result)
}

func writeODataError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("OData-Version", "4.0")
	_ = httputil.WriteJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
