package user

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/eventdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/project"
	"github.com/SURF-Innovatie/MORIS/internal/user"
	"github.com/google/uuid"
)

type Handler struct {
	svc     user.Service
	projSvc project.Service
}

func NewHandler(svc user.Service, projSvc project.Service) *Handler {
	return &Handler{svc: svc, projSvc: projSvc}
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
		httputil.WriteError(w, r, http.StatusBadRequest, "person ID is required", nil)
		return
	}

	u, err := h.svc.Create(r.Context(), entities.User{
		PersonID: req.PersonID,
		Password: req.Password,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), u.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, userdto.FromEntity(acc))
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
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, userdto.FromEntity(acc))
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
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
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
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, userdto.FromEntity(acc))
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
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
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
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	events, err := h.svc.GetApprovedEvents(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Map to DTO
	dtos := make([]eventdto.Event, 0, len(events))
	for _, e := range events {

		proj, _ := h.projSvc.GetProject(r.Context(), e.AggregateID())
		projectTitle := ""
		if proj != nil {
			projectTitle = proj.Project.Title
		}

		dtos = append(dtos, eventdto.FromEntityWithTitle(e, projectTitle))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, eventdto.Response{Events: dtos})
}

// ListUsers godoc
// @Summary Get all users (Admin only)
// @Description Returns a list of all users - requires admin role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "Admin user list"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 403 {object} httputil.BackendError "Insufficient permissions"
// @Router /admin/users/list [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Admin-only user list!", "users": [{"id":1,"name":"Admin User"}]}`))
}
