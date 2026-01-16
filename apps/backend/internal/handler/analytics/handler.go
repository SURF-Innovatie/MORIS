package analytics

import (
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/app/analytics"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles organization analytics HTTP requests
type Handler struct {
	service *analytics.Service
}

// NewHandler creates a new analytics handler
func NewHandler(service *analytics.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers analytics routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/organisations/{orgId}/analytics", func(r chi.Router) {
		r.Get("/summary", h.GetSummary)
		r.Get("/burn-rate", h.GetBurnRate)
		r.Get("/by-category", h.GetByCategory)
		r.Get("/by-project", h.GetByProject)
		r.Get("/by-funding", h.GetByFunding)
	})
}

// GetSummary godoc
// @Summary Get organization analytics summary
// @Description Returns overall statistics for all projects in an organization
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {object} analytics.OrgAnalyticsSummary
// @Failure 400 {string} string "invalid org id"
// @Failure 403 {string} string "forbidden"
// @Router /organisations/{orgId}/analytics/summary [get]
func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "orgId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	// TODO: Check RBAC - user must be org admin
	summary, err := h.service.GetOrgSummary(r.Context(), orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, summary)
}

// GetBurnRate godoc
// @Summary Get burn rate time-series data
// @Description Returns cumulative spending over time for projects in an organization
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} analytics.BurnRateDataPoint
// @Failure 400 {string} string "invalid parameters"
// @Router /organisations/{orgId}/analytics/burn-rate [get]
func (h *Handler) GetBurnRate(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "orgId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	params := analytics.DateRangeParams{}
	if startStr := r.URL.Query().Get("startDate"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			params.StartDate = &t
		}
	}
	if endStr := r.URL.Query().Get("endDate"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			params.EndDate = &t
		}
	}

	data, err := h.service.GetBurnRateData(r.Context(), orgID, params)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if data == nil {
		data = []analytics.BurnRateDataPoint{}
	}

	_ = httputil.WriteJSON(w, http.StatusOK, data)
}

// GetByCategory godoc
// @Summary Get spending breakdown by category
// @Description Returns budgeted vs actuals grouped by expense category
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {array} analytics.CategoryBreakdown
// @Failure 400 {string} string "invalid org id"
// @Router /organisations/{orgId}/analytics/by-category [get]
func (h *Handler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "orgId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	data, err := h.service.GetCategoryBreakdown(r.Context(), orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if data == nil {
		data = []analytics.CategoryBreakdown{}
	}

	_ = httputil.WriteJSON(w, http.StatusOK, data)
}

// GetByProject godoc
// @Summary Get project health summary
// @Description Returns health indicators for each project in the organization
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {array} analytics.ProjectHealthSummary
// @Failure 400 {string} string "invalid org id"
// @Router /organisations/{orgId}/analytics/by-project [get]
func (h *Handler) GetByProject(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "orgId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	data, err := h.service.GetProjectHealth(r.Context(), orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if data == nil {
		data = []analytics.ProjectHealthSummary{}
	}

	_ = httputil.WriteJSON(w, http.StatusOK, data)
}

// GetByFunding godoc
// @Summary Get spending breakdown by funding source
// @Description Returns budgeted vs actuals grouped by funding source
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {array} analytics.FundingBreakdown
// @Failure 400 {string} string "invalid org id"
// @Router /organisations/{orgId}/analytics/by-funding [get]
func (h *Handler) GetByFunding(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "orgId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	data, err := h.service.GetFundingBreakdown(r.Context(), orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if data == nil {
		data = []analytics.FundingBreakdown{}
	}

	_ = httputil.WriteJSON(w, http.StatusOK, data)
}

// Helper to extract user ID from context
func getUserIDFromContext(r *http.Request) uuid.UUID {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		return uuid.Nil
	}
	if id, ok := userIDVal.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
