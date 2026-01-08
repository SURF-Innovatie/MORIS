package person

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
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
// @Security BearerAuth
// @Param request body dto.PersonRequest true "Person details"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {string} string "Invalid body or missing required fields"
// @Failure 500 {string} string "Internal server error"
// @Router /people [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.PersonRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	p, err := h.svc.Create(r.Context(), entities.Person{
		Name:        req.Name,
		GivenName:   req.GivenName,
		FamilyName:  req.FamilyName,
		Email:       req.Email,
		ORCiD:       req.ORCiD,
		AvatarUrl:   req.AvatarURL,
		Description: req.Description,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.PersonResponse](*p))
}

// Get retrieves a person by ID
// @Summary Get a person
// @Description Get a person by their ID
// @Tags people
// @Produce json
// @Security BearerAuth
// @Param id path string true "Person ID"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Person not found"
// @Router /people/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.PersonResponse](*p))
}

// Update updates a person
// @Summary Update a person
// @Description Update a person's details
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Person ID"
// @Param request body dto.PersonRequest true "Person details"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {string} string "Invalid body or ID"
// @Failure 404 {string} string "Person not found"
// @Failure 500 {string} string "Internal server error"
// @Router /people/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.PersonRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	p, err := h.svc.Update(r.Context(), id, entities.Person{
		Name:        req.Name,
		GivenName:   req.GivenName,
		FamilyName:  req.FamilyName,
		Email:       req.Email,
		ORCiD:       req.ORCiD,
		AvatarUrl:   req.AvatarURL,
		Description: req.Description,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.PersonResponse](*p))
}
