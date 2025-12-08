package project

import (
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/changelogdto"
	_ "github.com/SURF-Innovatie/MORIS/internal/api/changelogdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/eventdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/api/productdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/projectdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	_ "github.com/SURF-Innovatie/MORIS/internal/domain/entities"

	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/project"
	"github.com/google/uuid"
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}
	proj, err := h.svc.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, "invalid startDate", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		http.Error(w, "invalid endDate", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	req := projectdto.Request{}
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		http.Error(w, "invalid startDate", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		http.Error(w, "invalid endDate", http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resps := make([]projectdto.Response, 0, len(projs))
	for _, p := range projs {
		resps = append(resps, toProjectResponse(p))
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	personID, err := httputil.ParseUUIDParam(r, "personId")
	if err != nil {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.AddPerson(r.Context(), id, personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	personId, err := httputil.ParseUUIDParam(r, "personId")
	if err != nil {
		http.Error(w, "invalid personId", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.RemovePerson(r.Context(), id, personId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	productID, err := httputil.ParseUUIDParam(r, "productID")
	if err != nil {
		http.Error(w, "invalid productID", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.AddProduct(r.Context(), id, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	productID, err := httputil.ParseUUIDParam(r, "productID")
	if err != nil {
		http.Error(w, "invalid productID", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.RemoveProduct(r.Context(), id, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toProjectResponse(proj))
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	changeLog, err := h.svc.GetChangeLog(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	changeDto := changelogdto.Changelog{
		Entries: make([]changelogdto.ChangelogEntry, 0, len(changeLog.Entries)),
	}
	for _, entry := range changeLog.Entries {
		changeDto.Entries = append(changeDto.Entries, changelogdto.ChangelogEntry{
			Event: entry.Event,
			At:    entry.At,
		})
	}

	_ = httputil.WriteJSON(w, http.StatusOK, changeDto)
}

func toOrganisationDTO(org entities.Organisation) organisationdto.Response {
	return organisationdto.Response{
		ID:   org.Id,
		Name: org.Name,
	}
}

func toPersonDTO(p entities.Person) persondto.Response {
	return persondto.Response{
		ID:          p.Id,
		UserID:      p.UserID,
		Name:        p.Name,
		GivenName:   p.GivenName,
		FamilyName:  p.FamilyName,
		Email:       p.Email,
		AvatarUrl:   p.AvatarUrl,
		ORCiD:       p.ORCiD,
		Description: p.Description,
	}
}

func toProductDTO(p entities.Product) productdto.Response {
	return productdto.Response{
		ID:       p.Id,
		Name:     p.Name,
		Language: p.Language,
		Type:     p.Type,
		DOI:      p.DOI,
	}
}

func toProjectResponse(d *entities.ProjectDetails) projectdto.Response {
	peopleDTOs := make([]persondto.Response, 0, len(d.People))
	for _, p := range d.People {
		peopleDTOs = append(peopleDTOs, toPersonDTO(p))
	}

	productDTOs := make([]productdto.Response, 0, len(d.Products))
	for _, p := range d.Products {
		productDTOs = append(productDTOs, toProductDTO(p))
	}

	return projectdto.Response{
		Id:           d.Project.Id,
		ProjectAdmin: d.Project.ProjectAdmin,
		Version:      d.Project.Version,
		Title:        d.Project.Title,
		Description:  d.Project.Description,
		StartDate:    d.Project.StartDate,
		EndDate:      d.Project.EndDate,
		Organization: toOrganisationDTO(d.Organisation),
		People:       peopleDTOs,
		Products:     productDTOs,
	}
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
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	events, err := h.svc.GetPendingEvents(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := eventdto.Response{
		Events: make([]eventdto.Event, 0, len(events)),
	}
	for _, e := range events {
		resp.Events = append(resp.Events, eventdto.Event{
			ID:        e.GetID(),
			ProjectID: e.AggregateID(),
			Type:      e.Type(),
			Status:    "pending",                                              // We filtered for pending
			CreatedBy: e.(interface{ CreatedByID() uuid.UUID }).CreatedByID(), // Safe cast as Base has it
			At:        e.OccurredAt(),
			Details:   e.String(),
		})
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}
