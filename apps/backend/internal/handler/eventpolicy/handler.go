package eventpolicy

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/domain/policy"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Handler handles HTTP requests for event policies
type Handler struct {
	svc eventpolicy.Service
}

// NewHandler creates a new event policy handler
func NewHandler(svc eventpolicy.Service) *Handler {
	return &Handler{svc: svc}
}

// ListForOrgNode godoc
// @Summary List event policies for an organisation node
// @Description Retrieves all event policies for an org node, including inherited policies
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation Node ID (UUID)"
// @Param inherited query bool false "Include inherited policies from parent orgs"
// @Success 200 {array} dto.EventPolicyResponse
// @Failure 400 {string} string "invalid org node id"
// @Failure 500 {string} string "internal server error"
// @Router /organisations/{id}/policies [get]
func (h *Handler) ListForOrgNode(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org node id", nil)
		return
	}

	includeInherited := r.URL.Query().Get("inherited") != "false" // default true

	policies, err := h.svc.ListForOrgNode(r.Context(), id, includeInherited)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := lo.Map(policies, func(p policy.EventPolicy, _ int) dto.EventPolicyResponse {
		var r dto.EventPolicyResponse
		r.FromEntity(&p)
		return r
	})

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// CreateForOrgNode godoc
// @Summary Create an event policy for an organisation node
// @Description Creates a new event policy scoped to an organisation node
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation Node ID (UUID)"
// @Param body body dto.EventPolicyRequest true "Policy details"
// @Success 201 {object} dto.EventPolicyResponse
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /organisations/{id}/policies [post]
func (h *Handler) CreateForOrgNode(w http.ResponseWriter, r *http.Request) {
	orgNodeID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org node id", nil)
		return
	}

	var req dto.EventPolicyRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	policy, err := h.requestToEntity(req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	policy.OrgNodeID = &orgNodeID

	created, err := h.svc.Create(r.Context(), policy)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var resp dto.EventPolicyResponse
	resp.FromEntity(created)
	_ = httputil.WriteJSON(w, http.StatusCreated, resp)
}

// ListForProject godoc
// @Summary List event policies for a project
// @Description Retrieves all event policies for a project, including inherited org policies
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Param org_node_id query string true "Owning Organisation Node ID (UUID)"
// @Param inherited query bool false "Include inherited policies from org hierarchy"
// @Success 200 {array} dto.EventPolicyResponse
// @Failure 400 {string} string "invalid project id"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/policies [get]
func (h *Handler) ListForProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	orgNodeIDStr := r.URL.Query().Get("org_node_id")
	orgNodeID, err := uuid.Parse(orgNodeIDStr)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org_node_id query param", nil)
		return
	}

	includeInherited := r.URL.Query().Get("inherited") != "false"

	policies, err := h.svc.ListForProject(r.Context(), projectID, orgNodeID, includeInherited)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := lo.Map(policies, func(p policy.EventPolicy, _ int) dto.EventPolicyResponse {
		var r dto.EventPolicyResponse
		r.FromEntity(&p)
		return r
	})

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// CreateForProject godoc
// @Summary Create an event policy for a project
// @Description Creates a new event policy scoped to a project
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Param body body dto.EventPolicyRequest true "Policy details"
// @Success 201 {object} dto.EventPolicyResponse
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/policies [post]
func (h *Handler) CreateForProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	var req dto.EventPolicyRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	policy, err := h.requestToEntity(req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	policy.ProjectID = &projectID

	created, err := h.svc.Create(r.Context(), policy)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var resp dto.EventPolicyResponse
	resp.FromEntity(created)
	_ = httputil.WriteJSON(w, http.StatusCreated, resp)
}

// GetPolicy godoc
// @Summary Get an event policy by ID
// @Description Retrieves a single event policy
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Policy ID (UUID)"
// @Success 200 {object} dto.EventPolicyResponse
// @Failure 400 {string} string "invalid policy id"
// @Failure 404 {string} string "policy not found"
// @Router /policies/{id} [get]
func (h *Handler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid policy id", nil)
		return
	}

	policy, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "policy not found", nil)
		return
	}

	var resp dto.EventPolicyResponse
	resp.FromEntity(policy)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// UpdatePolicy godoc
// @Summary Update an event policy
// @Description Updates an existing event policy
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Policy ID (UUID)"
// @Param body body dto.EventPolicyRequest true "Updated policy details"
// @Success 200 {object} dto.EventPolicyResponse
// @Failure 400 {string} string "invalid request"
// @Failure 404 {string} string "policy not found"
// @Router /policies/{id} [put]
func (h *Handler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid policy id", nil)
		return
	}

	var req dto.EventPolicyRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	policy, err := h.requestToEntity(req)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	updated, err := h.svc.Update(r.Context(), id, policy)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var resp dto.EventPolicyResponse
	resp.FromEntity(updated)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// DeletePolicy godoc
// @Summary Delete an event policy
// @Description Deletes an event policy
// @Tags event-policies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Policy ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {string} string "invalid policy id"
// @Failure 500 {string} string "internal server error"
// @Router /policies/{id} [delete]
func (h *Handler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid policy id", nil)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// requestToEntity converts request DTO to domain entity
func (h *Handler) requestToEntity(req dto.EventPolicyRequest) (policy.EventPolicy, error) {
	eventPolicy := policy.EventPolicy{
		Name:             req.Name,
		Description:      req.Description,
		EventTypes:       req.EventTypes,
		ActionType:       policy.ActionType(req.ActionType),
		MessageTemplate:  req.MessageTemplate,
		RecipientDynamic: req.RecipientDynamic,
		Enabled:          req.Enabled,
	}

	// Convert conditions
	for _, c := range req.Conditions {
		eventPolicy.Conditions = append(eventPolicy.Conditions, policy.Condition{
			Field:    c.Field,
			Operator: c.Operator,
			Value:    c.Value,
		})
	}

	// Parse user IDs
	for _, uidStr := range req.RecipientUserIDs {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			return eventPolicy, err
		}
		eventPolicy.RecipientUserIDs = append(eventPolicy.RecipientUserIDs, uid)
	}

	// Parse project role IDs
	for _, ridStr := range req.RecipientProjectRoleIDs {
		rid, err := uuid.Parse(ridStr)
		if err != nil {
			return eventPolicy, err
		}
		eventPolicy.RecipientProjectRoleIDs = append(eventPolicy.RecipientProjectRoleIDs, rid)
	}

	// Parse org role IDs
	for _, ridStr := range req.RecipientOrgRoleIDs {
		rid, err := uuid.Parse(ridStr)
		if err != nil {
			return eventPolicy, err
		}
		eventPolicy.RecipientOrgRoleIDs = append(eventPolicy.RecipientOrgRoleIDs, rid)
	}

	return eventPolicy, nil
}
