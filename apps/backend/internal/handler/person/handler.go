package person

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/persondto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/person"
)

type Handler struct {
	svc personsvc.Service
}

func NewHandler(s personsvc.Service) *Handler {
	return &Handler{svc: s}
}

// Create creates a new person
// @Summary Create a new person
// @Description Create a new person with the provided details
// @Tags people
// @Accept json
// @Produce json
// @Param request body persondto.Request true "Person details"
// @Success 200 {object} persondto.Response
// @Failure 400 {string} string "Invalid body or missing required fields"
// @Failure 500 {string} string "Internal server error"
// @Router /people [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req persondto.Request
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	p, err := h.svc.Create(r.Context(), entities.Person{
		Name:       req.Name,
		GivenName:  req.GivenName,
		FamilyName: req.FamilyName,
		Email:      req.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, persondto.Response{
		ID:         p.Id,
		UserID:     p.UserID,
		Name:       p.Name,
		GivenName:  p.GivenName,
		FamilyName: p.FamilyName,
		Email:      p.Email,
	})
}

// Get retrieves a person by ID
// @Summary Get a person
// @Description Get a person by their ID
// @Tags people
// @Produce json
// @Param id path string true "Person ID"
// @Success 200 {object} persondto.Response
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Person not found"
// @Router /people/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, persondto.Response{
		ID:         p.Id,
		UserID:     p.UserID,
		Name:       p.Name,
		GivenName:  p.GivenName,
		FamilyName: p.FamilyName,
		Email:      p.Email,
	})
}

// Update updates a person
// @Summary Update a person
// @Description Update a person's details
// @Tags people
// @Accept json
// @Produce json
// @Param id path string true "Person ID"
// @Param request body persondto.Request true "Person details"
// @Success 200 {object} persondto.Response
// @Failure 400 {string} string "Invalid body or ID"
// @Failure 404 {string} string "Person not found"
// @Failure 500 {string} string "Internal server error"
// @Router /people/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req persondto.Request
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	p, err := h.svc.Update(r.Context(), id, entities.Person{
		Name:        req.Name,
		GivenName:   req.GivenName,
		FamilyName:  req.FamilyName,
		Email:       req.Email,
		AvatarUrl:   req.AvatarUrl,
		Description: req.Description,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, persondto.Response{
		ID:          p.Id,
		UserID:      p.UserID,
		Name:        p.Name,
		GivenName:   p.GivenName,
		FamilyName:  p.FamilyName,
		Email:       p.Email,
		AvatarUrl:   p.AvatarUrl,
		Description: p.Description,
	})
}

// List implementation omitted for brevity.
