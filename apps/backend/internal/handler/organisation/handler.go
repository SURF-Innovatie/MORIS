package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	organisationsvc "github.com/SURF-Innovatie/MORIS/internal/organisation"
)

type Handler struct {
	svc organisationsvc.Service
}

func NewHandler(s organisationsvc.Service) *Handler {
	return &Handler{svc: s}
}

// CreateRoot godoc
// @Summary Create a root organisation node
// @Description Creates a new root node in the organisation tree (no parent)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body organisationdto.CreateRootRequest true "Create root node request"
// @Success 200 {object} organisationdto.Response
// @Failure 400 {string} string "name is required / invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes [post]
func (h *Handler) CreateRoot(w http.ResponseWriter, r *http.Request) {
	var req organisationdto.CreateRootRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	node, err := h.svc.CreateRoot(r.Context(), req.Name)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.FromEntity(*node))
}

// CreateChild godoc
// @Summary Create a child organisation node
// @Description Creates a new child node under the given parent node
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Parent node ID"
// @Param body body organisationdto.CreateChildRequest true "Create child node request"
// @Success 200 {object} organisationdto.Response
// @Failure 400 {string} string "invalid id / name is required / invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/children [post]
func (h *Handler) CreateChild(w http.ResponseWriter, r *http.Request) {
	parentID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req organisationdto.CreateChildRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	node, err := h.svc.CreateChild(r.Context(), parentID, req.Name)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.FromEntity(*node))
}

// GetOrganisationNode godoc
// @Summary GetOrganisationNode an organisation node
// @Description Retrieves a single organisation node by ID
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation node ID"
// @Success 200 {object} organisationdto.Response
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id} [get]
func (h *Handler) GetOrganisationNode(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	node, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.FromEntity(*node))
}

// UpdateOrganisationNode godoc
// @Summary UpdateOrganisationNode an organisation node
// @Description Updates an organisation node's name and/or parent (re-parenting). If parentID is null, the node becomes a root node.
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation node ID"
// @Param body body organisationdto.UpdateRequest true "UpdateOrganisationNode organisation node request"
// @Success 200 {object} organisationdto.Response
// @Failure 400 {string} string "invalid id / name is required / invalid request body"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id} [patch]
func (h *Handler) UpdateOrganisationNode(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req organisationdto.UpdateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	node, err := h.svc.Update(r.Context(), id, req.Name, req.ParentID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.FromEntity(*node))
}

// ListRoots godoc
// @Summary List root organisation nodes
// @Description Returns all organisation nodes that have no parent
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} organisationdto.Response
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/roots [get]
func (h *Handler) ListRoots(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.svc.ListRoots(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := make([]organisationdto.Response, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, organisationdto.FromEntity(n))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}

// ListChildren godoc
// @Summary List children of an organisation node
// @Description Returns the direct children of a given organisation node
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Parent node ID"
// @Success 200 {array} organisationdto.Response
// @Failure 401 {string} string "unauthorized"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/children [get]
func (h *Handler) ListChildren(w http.ResponseWriter, r *http.Request) {
	parentID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	nodes, err := h.svc.ListChildren(r.Context(), parentID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := make([]organisationdto.Response, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, organisationdto.FromEntity(n))
	}

	_ = httputil.WriteJSON(w, http.StatusOK, out)
}
