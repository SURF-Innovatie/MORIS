package project

import (
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/changelogdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/eventdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/projectdto"

	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/project"
)

type Handler struct {
	svc project.Service
}

func NewHandler(svc project.Service) *Handler {
	return &Handler{svc: svc}
}

// GetProject godoc
// @Summary Get a project by ID
// @Description Retrieves a single project by its unique identifier
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid project id"
// @Failure 404 {string} string "project not found"
// @Router /projects/{id} [get]
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}
	proj, err := h.svc.GetProject(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// StartProject godoc
// @Summary Start a new project
// @Description Creates and starts a new project with the provided details
// @Tags projects
// @Accept json
// @Produce json
// @Param project body projectdto.Request true "Project details"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid body or date format"
// @Failure 500 {string} string "internal server error"
// @Router /projects [post]
func (h *Handler) StartProject(w http.ResponseWriter, r *http.Request) {
	req := projectdto.Request{}
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid startDate", nil)
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid endDate", nil)
		return
	}

	params := project.StartProjectParams{
		ProjectAdmin:   req.ProjectAdmin,
		Title:          req.Title,
		Description:    req.Description,
		OrganisationID: req.OrganisationID,
		StartDate:      start,
		EndDate:        end,
	}

	proj, err := h.svc.StartProject(r.Context(), params)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// UpdateProject godoc
// @Summary Update a project
// @Description Updates an existing project with the provided details
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param project body projectdto.Request true "Project details"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid body, id or date format"
// @Failure 404 {string} string "project not found"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id} [put]
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	req := projectdto.Request{}
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid startDate", nil)
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid endDate", nil)
		return
	}

	params := project.UpdateProjectParams{
		Title:          req.Title,
		Description:    req.Description,
		OrganisationID: req.OrganisationID,
		StartDate:      start,
		EndDate:        end,
	}

	proj, err := h.svc.UpdateProject(r.Context(), id, params)
	if err != nil {
		if err == project.ErrNotFound {
			httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// GetAllProjects godoc
// @Summary Get all projects
// @Description Retrieves a list of all projects
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} projectdto.Response
// @Failure 500 {string} string "internal server error"
// @Router /projects [get]
func (h *Handler) GetAllProjects(w http.ResponseWriter, r *http.Request) {

	projs, err := h.svc.GetAllProjects(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	resps := make([]projectdto.Response, 0, len(projs))
	for _, p := range projs {
		resps = append(resps, projectdto.FromEntity(*p))
	}
	_ = httputil.WriteJSON(w, http.StatusOK, resps)
}

// AddPerson godoc
// @Summary Add a person to a project
// @Description Adds a new person to the specified project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param personId path string true "Person ID (UUID)"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid project id or body"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/people/{personId} [post]
func (h *Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	personID, err := httputil.ParseUUIDParam(r, "personId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid person id", nil)
		return
	}

	proj, err := h.svc.AddPerson(r.Context(), id, personID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// RemovePerson godoc
// @Summary Remove a person from a project
// @Description Removes a person from the specified project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param personId path string true "Person ID (UUID)"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid project id or person id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/people/{personId} [delete]
func (h *Handler) RemovePerson(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	personId, err := httputil.ParseUUIDParam(r, "personId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid personId", nil)
		return
	}

	proj, err := h.svc.RemovePerson(r.Context(), id, personId)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// AddProduct godoc
// @Summary Add a product to a project
// @Description Adds an existing product to the specified project using its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param productID path string true "Product ID (UUID)"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid project id or product id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/products/{productID} [post]
func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	productID, err := httputil.ParseUUIDParam(r, "productID")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid productID", nil)
		return
	}

	proj, err := h.svc.AddProduct(r.Context(), id, productID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// RemoveProduct godoc
// @Summary Remove a product from a project
// @Description Removes a product association from the specified project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param productID path string true "Product ID (UUID)"
// @Success 200 {object} projectdto.Response
// @Failure 400 {string} string "invalid project id or product id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/products/{productID} [delete]
func (h *Handler) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	productID, err := httputil.ParseUUIDParam(r, "productID")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid productID", nil)
		return
	}

	proj, err := h.svc.RemoveProduct(r.Context(), id, productID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, projectdto.FromEntity(*proj))
}

// GetChangelog godoc
// @Summary Get change log for a project
// @Description Retrieves the change log for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} changelogdto.Changelog
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/changelog [get]
func (h *Handler) GetChangelog(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	changeLog, err := h.svc.GetChangeLog(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, changelogdto.FromEntity(*changeLog))
}

// GetPendingEvents godoc
// @Summary Get pending events for a project
// @Description Retrieves a list of pending events for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Success 200 {object} eventdto.Response
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/pending-events [get]
func (h *Handler) GetPendingEvents(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	pendingEvents, err := h.svc.GetPendingEvents(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := eventdto.Response{
		Events: make([]eventdto.Event, 0, len(pendingEvents)),
	}
	for _, e := range pendingEvents {
		resp.Events = append(resp.Events, eventdto.FromEntity(e))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}
