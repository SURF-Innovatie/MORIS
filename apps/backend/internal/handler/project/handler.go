package project

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

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
// @Success 200 {object} entities.Project
// @Failure 400 {string} string "invalid project id"
// @Failure 404 {string} string "project not found"
// @Router /projects/{id} [get]
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}
	proj, err := h.svc.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

// StartProject godoc
// @Summary Start a new project
// @Description Creates and starts a new project with the provided details
// @Tags projects
// @Accept json
// @Produce json
// @Param project body StartRequest true "Project details"
// @Success 200 {object} entities.Project
// @Failure 400 {string} string "invalid body or date format"
// @Failure 500 {string} string "internal server error"
// @Router /projects [post]
func (h *Handler) StartProject(w http.ResponseWriter, r *http.Request) {
	req := StartRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
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
		Title:        req.Title,
		Description:  req.Description,
		Organisation: req.Organisation,
		StartDate:    start,
		EndDate:      end,
	}

	proj, err := h.svc.StartProject(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

// GetAllProjects godoc
// @Summary Get all projects
// @Description Retrieves a list of all projects
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} entities.Project
// @Failure 500 {string} string "internal server error"
// @Router /projects [get]
func (h *Handler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	proj, err := h.svc.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

type personRequest struct {
	Name string `json:"name"`
}

// AddPerson godoc
// @Summary Add a person to a project
// @Description Adds a new person to the specified project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param person body personRequest true "Person details"
// @Success 200 {object} entities.Project
// @Failure 400 {string} string "invalid project id or body"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/people [post]
func (h *Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	personIDstr := chi.URLParam(r, "personId")
	personID, err := uuid.Parse(personIDstr)

	proj, err := h.svc.AddPerson(r.Context(), id, personID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

// RemovePerson godoc
// @Summary Remove a person from a project
// @Description Removes a person from the specified project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID (UUID)"
// @Param personId path string true "Person ID (UUID)"
// @Success 200 {object} entities.Project
// @Failure 400 {string} string "invalid project id or person id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/people/{personId} [delete]
func (h *Handler) RemovePerson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	personIdStr := chi.URLParam(r, "personId")
	personId, err := uuid.Parse(personIdStr)
	if err != nil {
		http.Error(w, "invalid personId", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.RemovePerson(r.Context(), id, personId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}
