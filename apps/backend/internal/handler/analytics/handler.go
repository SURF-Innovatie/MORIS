package analytics

import (
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/analytics"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
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
// @Success 200 {object} dto.OrgAnalyticsSummaryResponse
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

	resp := transform.ToDTOItem[dto.OrgAnalyticsSummaryResponse](summary)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
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
// @Success 200 {array} dto.BurnRateDataPointResponse
// @Failure 400 {string} string "invalid parameters"
// @Router /organisations/{orgId}/analytics/burn-rate [get]
func (h *Handler) GetBurnRate(w http.ResponseWriter, r *http.Request) {
	orgID, err := httputil.ParseUUIDParam(r, "orgId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	params := entities.DateRangeParams{}
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

	resp := transform.ToDTOs[dto.BurnRateDataPointResponse](data)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetByCategory godoc
// @Summary Get spending breakdown by category
// @Description Returns budgeted vs actuals grouped by expense category
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {array} dto.CategoryBreakdownResponse
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

	resp := transform.ToDTOs[dto.CategoryBreakdownResponse](data)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetByProject godoc
// @Summary Get project health summary
// @Description Returns health indicators for each project in the organization
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {array} dto.ProjectHealthSummaryResponse
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

	resp := transform.ToDTOs[dto.ProjectHealthSummaryResponse](data)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetByFunding godoc
// @Summary Get spending breakdown by funding source
// @Description Returns budgeted vs actuals grouped by funding source
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orgId path string true "Organisation ID (UUID)"
// @Success 200 {array} dto.FundingBreakdownResponse
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

	resp := transform.ToDTOs[dto.FundingBreakdownResponse](data)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}
