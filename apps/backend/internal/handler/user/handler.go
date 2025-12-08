package user

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/eventdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	"github.com/google/uuid"
)

type Handler struct {
	svc user.Service
}

func NewHandler(svc user.Service) *Handler {
	return &Handler{svc: svc}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account for an existing person using the provided person ID and password.
// @Tags users
// @Accept json
// @Produce json
// @Param user body userdto.Request true "User creation payload"
// @Success 200 {object} userdto.Response
// @Failure 400 {string} string "invalid request body or missing person ID"
// @Failure 500 {string} string "internal server error"
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userdto.Request
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.PersonID == uuid.Nil {
		http.Error(w, "person ID is required", http.StatusBadRequest)
		return
	}

	u, err := h.svc.Create(r.Context(), entities.User{
		PersonID: req.PersonID,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toUserResponse(acc))
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves a single user by its ID, provided as the `id` query parameter.
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User ID (UUID)"
// @Success 200 {object} userdto.Response
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /users [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDQuery(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toUserResponse(acc))
}

// UpdateUser godoc
// @Summary Update a user
// @Description Updates an existing user's person reference and/or password based on the given ID and request body.
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User ID (UUID)"
// @Param user body userdto.Request true "User update payload"
// @Success 200 {object} userdto.Response
// @Failure 400 {string} string "invalid id or request body"
// @Failure 500 {string} string "internal server error"
// @Router /users [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDQuery(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req userdto.Request
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	_, err = h.svc.Update(r.Context(), id, entities.User{
		PersonID: req.PersonID,
		Password: req.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, toUserResponse(acc))
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user by its ID, provided as the `id` query parameter.
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User ID (UUID)"
// @Success 204 {string} string "no content"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /users [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDQuery(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetApprovedEvents godoc
// @Summary Get approved events for a user
// @Description Retrieves all approved events created by the user with the given ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} eventdto.Response
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /users/{id}/events/approved [get]
func (h *Handler) GetApprovedEvents(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	events, err := h.svc.GetApprovedEvents(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map to DTO
	dtos := make([]eventdto.Event, 0, len(events))
	for _, e := range events {
		dtos = append(dtos, eventdto.Event{
			ID:        e.GetID(),
			ProjectID: e.AggregateID(),
			Type:      e.Type(),
			Status:    e.GetStatus(),
			CreatedBy: e.CreatedByID(),
			At:        e.OccurredAt(),
			Details:   e.String(),
		})
	}

	_ = httputil.WriteJSON(w, http.StatusOK, eventdto.Response{Events: dtos})
}

func toUserResponse(acc *entities.UserAccount) userdto.Response {
	p := acc.Person
	u := acc.User

	return userdto.Response{
		ID:         u.ID,
		PersonID:   u.PersonID,
		Email:      p.Email,
		Name:       p.Name,
		ORCiD:      p.ORCiD,
		GivenName:  p.GivenName,
		FamilyName: p.FamilyName,
	}
}
