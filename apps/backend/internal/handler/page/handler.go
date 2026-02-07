package page

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/page"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Handler struct {
	service page.Service
}

func NewHandler(service page.Service) *Handler {
	return &Handler{service: service}
}

// Get Page by slug
// @Summary Get a page by slug
// @Description get page by slug
// @Tags page
// @Accept  json
// @Produce  json
// @Param slug path string true "Page Slug"
// @Success 200 {object} dto.PageResponse
// @Failure 404 {object} string
// @Router /pages/{slug} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug required", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetPage(r.Context(), slug)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, mapEntityToResponse(p))
}

// Create Page
// @Summary Create a new page
// @Description create page with content
// @Tags page
// @Accept  json
// @Produce  json
// @Param page body dto.CreatePageRequest true "Create Page Request"
// @Success 201 {object} dto.PageResponse
// @Failure 400 {object} string
// @Router /pages [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.service.CreatePage(r.Context(), req.Title, req.Slug, req.Type, req.Content, req.ProjectID, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, mapEntityToResponse(p))
}

// Update Page
// @Summary Update a page
// @Description update page content
// @Tags page
// @Accept  json
// @Produce  json
// @Param id path string true "Page ID"
// @Param page body dto.UpdatePageRequest true "Update Page Request"
// @Success 200 {object} dto.PageResponse
// @Failure 404 {object} string
// @Router /pages/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.service.UpdatePage(r.Context(), id, req.Title, req.Content, req.IsPublished)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, mapEntityToResponse(p))
}

// ListByProject
// @Summary List pages for a project
// @Description list pages associated with a project
// @Tags page
// @Accept  json
// @Produce  json
// @Param projectId path string true "Project ID"
// @Success 200 {array} dto.PageResponse
// @Router /projects/{projectId}/pages [get]
func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "projectId")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	pages, err := h.service.ListProjectPages(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res []dto.PageResponse
	for _, p := range pages {
		res = append(res, mapEntityToResponse(p))
	}
	render.JSON(w, r, res)
}

// ListByUser
// @Summary List pages for a user
// @Description list pages associated with a user
// @Tags page
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Success 200 {array} dto.PageResponse
// @Router /users/{userId}/pages [get]
func (h *Handler) ListByUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	pages, err := h.service.ListUserPages(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res []dto.PageResponse
	for _, p := range pages {
		res = append(res, mapEntityToResponse(p))
	}
	render.JSON(w, r, res)
}

func mapEntityToResponse(p *ent.Page) dto.PageResponse {
	return dto.PageResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Type:        string(p.Type),
		Content:     p.Content,
		IsPublished: p.IsPublished,
		ProjectID:   p.ProjectID,
		UserID:      p.UserID, // Needs edge loading or just use ID? Wait, ProjectID is ID, UserID is ID. Correct.
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
