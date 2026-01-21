package nwo

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	ex "github.com/SURF-Innovatie/MORIS/external/nwo"
	app "github.com/SURF-Innovatie/MORIS/internal/app/nwo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

// Handler handles HTTP requests for NWO Open API operations
type Handler struct {
	svc app.Service
}

// NewHandler creates a new Handler with the given service
func NewHandler(svc app.Service) *Handler {
	return &Handler{svc: svc}
}

// GetProjects godoc
// @Summary Query NWO projects
// @Description Retrieves NWO-subsidized projects with optional filters
// @Tags nwo
// @Accept json
// @Produce json
// @Param project_id query string false "Project ID (dossiernummer)"
// @Param grant_id query string false "Grant ID (DOI)"
// @Param ror_id query string false "ROR ID of organisation"
// @Param organisation query string false "Organisation name"
// @Param title query string false "Project title (partial match)"
// @Param reporting_year query int false "Reporting year"
// @Param member_last_name query string false "Project member last name"
// @Param orcid query string false "ORCID of project member"
// @Param per_page query int false "Results per page"
// @Param page query int false "Page number"
// @Success 200 {object} ex.ProjectsResponse
// @Failure 400 {object} httputil.BackendError "bad request"
// @Failure 500 {object} httputil.BackendError "internal server error"
// @Router /nwo/projects [get]
func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	opts := parseQueryOptions(r)

	resp, err := h.svc.GetProjects(r.Context(), opts)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetProject godoc
// @Summary Get a single NWO project
// @Description Retrieves a single NWO project by its project_id
// @Tags nwo
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID (dossiernummer)"
// @Success 200 {object} ex.Project
// @Failure 400 {object} httputil.BackendError "project_id is required"
// @Failure 404 {object} httputil.BackendError "project not found"
// @Failure 500 {object} httputil.BackendError "internal server error"
// @Router /nwo/projects/{project_id} [get]
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	if projectID == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "project_id is required", nil)
		return
	}

	project, err := h.svc.GetProject(r.Context(), projectID)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httputil.WriteError(w, r, http.StatusNotFound, "project not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, project)
}

func parseQueryOptions(r *http.Request) *ex.QueryOptions {
	opts := &ex.QueryOptions{}

	opts.ProjectID = r.URL.Query().Get("project_id")
	opts.GrantID = r.URL.Query().Get("grant_id")
	opts.RORID = r.URL.Query().Get("ror_id")
	opts.Organisation = r.URL.Query().Get("organisation")
	opts.Title = r.URL.Query().Get("title")
	opts.Summary = r.URL.Query().Get("summary")
	opts.MemberLastName = r.URL.Query().Get("member_last_name")
	opts.ORCID = r.URL.Query().Get("orcid")

	if roleStr := r.URL.Query().Get("role"); roleStr != "" {
		opts.Role = ex.ProjectRole(roleStr)
	}

	if yearStr := r.URL.Query().Get("reporting_year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			opts.ReportingYear = year
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil {
			opts.PerPage = perPage
		}
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			opts.Page = page
		}
	}

	return opts
}
