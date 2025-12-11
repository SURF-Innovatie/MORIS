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
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} userdto.Response
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /users/{id} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
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
// @Param id path string true "User ID (UUID)"
// @Param user body userdto.Request true "User update payload"
// @Success 200 {object} userdto.Response
// @Failure 400 {string} string "invalid id or request body"
// @Failure 500 {string} string "internal server error"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	// Permission check: User can update themselves, SysAdmin can update anyone
	authUser, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if authUser.User.ID != id && !authUser.User.IsSysAdmin {
		httputil.WriteError(w, r, http.StatusForbidden, "Forbidden: Can only update own profile", nil)
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
// @Param id path string true "User ID (UUID)"
// @Success 204 {string} string "no content"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	authUser, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	if !authUser.User.IsSysAdmin {
		httputil.WriteError(w, r, http.StatusForbidden, "insufficient permissions", nil)
		return
	}

	// Prevent deleting yourself
	if authUser.User.ID == id {
		httputil.WriteError(w, r, http.StatusBadRequest, "cannot delete yourself", nil)
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
// @Description Returns a paginated list of all users - requires admin role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10)"
// @Success 200 {object} userdto.PaginatedResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 403 {object} httputil.BackendError "Insufficient permissions"
// @Router /admin/users/list [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := httputil.ParseIntQuery(r, "page", 1)
	pageSize := httputil.ParseIntQuery(r, "page_size", 10)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, total, err := h.svc.ListAll(r.Context(), pageSize, offset)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	dtos := make([]userdto.Response, 0, len(users))
	for _, acc := range users {
		dtos = append(dtos, userdto.FromEntity(acc))
	}

	totalPages := (total + pageSize - 1) / pageSize

	resp := userdto.PaginatedResponse{
		Data:       dtos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// ToggleActive godoc
// @Summary Toggle user active status (Admin only)
// @Description Toggle user active status - requires admin role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body userdto.ToggleActiveRequest true "Toggle active status payload"
// @Success 200 {object} httputil.StatusResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 403 {object} httputil.BackendError "Insufficient permissions"
// @Router /admin/users/{id}/toggle-active [post]
func (h *Handler) ToggleActive(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	authUser, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	if !authUser.User.IsSysAdmin {
		httputil.WriteError(w, r, http.StatusForbidden, "insufficient permissions", nil)
		return
	}

	// Prevent deactivating yourself
	if authUser.User.ID == id {
		httputil.WriteError(w, r, http.StatusBadRequest, "cannot deactivate yourself", nil)
		return
	}

	// Get current status to toggle
	_, err = h.svc.GetAccount(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "user not found", nil)
		return
	}

	var req userdto.ToggleActiveRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	err = h.svc.ToggleActive(r.Context(), id, req.IsActive)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteStatus(w)
}
