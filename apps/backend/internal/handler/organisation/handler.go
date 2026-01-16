package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	organisationsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	rbacsvc "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	orgrole "github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/app/projectrole"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	svc            organisationsvc.Service
	rbac           rbacsvc.Service
	roleSvc        projectrole.Service
	customFieldSvc customfield.Service
}

func NewHandler(s organisationsvc.Service, r rbacsvc.Service, rs projectrole.Service, cfs customfield.Service) *Handler {
	return &Handler{svc: s, rbac: r, roleSvc: rs, customFieldSvc: cfs}
}

// CreateRoot godoc
// @Summary Create a root organisation node
// @Description Creates a new root node in the organisation tree (no parent)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.OrganisationCreateRootRequest true "Create root node request"
// @Success 200 {object} dto.OrganisationResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "name is required / invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes [post]
func (h *Handler) CreateRoot(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if !user.User.IsSysAdmin {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	var req dto.OrganisationCreateRootRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	node, err := h.svc.CreateRoot(r.Context(), req.Name, req.RorID, req.Description, req.AvatarURL)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationResponse](*node))
}

// CreateChild godoc
// @Summary Create a child organisation node
// @Description Creates a new child node under the given parent node
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Parent node ID"
// @Param body body dto.OrganisationCreateChildRequest true "Create child node request"
// @Success 200 {object} dto.OrganisationResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id / name is required / invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/children [post]
func (h *Handler) CreateChild(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	parentID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	// Check permissions: SysAdmin OR Admin access on parent
	if !user.User.IsSysAdmin {
		hasAccess, err := h.rbac.HasAdminAccess(r.Context(), user.Person.ID, parentID)
		if err != nil {
			httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		if !hasAccess {
			httputil.WriteError(w, r, http.StatusForbidden, "forbidden: requires admin rights on parent", nil)
			return
		}
	}

	var req dto.OrganisationCreateChildRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	node, err := h.svc.CreateChild(r.Context(), parentID, req.Name, req.RorID, req.Description, req.AvatarURL)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationResponse](*node))
}

// GetOrganisationNode godoc
// @Summary Get an organisation node
// @Description Retrieves a single organisation node by ID
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "OwningOrgNode node ID"
// @Success 200 {object} dto.OrganisationResponse
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

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationResponse](*node))
}

// UpdateOrganisationNode godoc
// @Summary Update an organisation node
// @Description Updates an organisation node's name and/or parent (re-parenting). If parentID is null, the node becomes a root node.
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "OwningOrgNode node ID"
// @Param body body dto.OrganisationUpdateRequest true "UpdateOrganisationNode organisation node request"
// @Success 200 {object} dto.OrganisationResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id / name is required / invalid request body"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id} [patch]
func (h *Handler) UpdateOrganisationNode(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.OrganisationUpdateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	// Check permissions
	if !user.User.IsSysAdmin {
		// Must have access to the node itself (to edit/move it)
		hasAccess, err := h.rbac.HasAdminAccess(r.Context(), user.Person.ID, id)
		if err != nil {
			httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		if !hasAccess {
			httputil.WriteError(w, r, http.StatusForbidden, "forbidden: requires admin rights on node", nil)
			return
		}

		// If re-parenting, must have access to new parent too
		if req.ParentID != nil {
			accessToNewParent, err := h.rbac.HasAdminAccess(r.Context(), user.Person.ID, *req.ParentID)
			if err != nil {
				httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
				return
			}
			if !accessToNewParent {
				httputil.WriteError(w, r, http.StatusForbidden, "forbidden: requires admin rights on new parent", nil)
				return
			}
		} else {
			// Moving to root requires SysAdmin (or we can decide if top-level admins can do this, but safe to restrict to SysAdmin)
			httputil.WriteError(w, r, http.StatusForbidden, "forbidden: only sysadmin can verify root nodes", nil)
			return
		}
	}

	node, err := h.svc.Update(r.Context(), id, req.Name, req.ParentID, req.RorID, req.Description, req.AvatarURL)
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.OrganisationResponse](*node))
}

// ListRoots godoc
// @Summary List root organisation nodes
// @Description Returns all organisation nodes that have no parent
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.OrganisationResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/roots [get]
func (h *Handler) ListRoots(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.svc.ListRoots(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.OrganisationResponse](nodes))
}

// ListChildren godoc
// @Summary List children of an organisation node
// @Description Returns the direct children of a given organisation node
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Parent node ID"
// @Success 200 {array} dto.OrganisationResponse
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

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.OrganisationResponse](nodes))
}

// Search godoc
// @Summary Search organisation nodes
// @Description Search for organisation nodes by name
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param permission query string false "Permission filter"
// @Success 200 {array} dto.OrganisationResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/search [get]
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		_ = httputil.WriteJSON(w, http.StatusOK, []dto.OrganisationResponse{})
		return
	}

	permission := r.URL.Query().Get("permission")

	var nodes []entities.OrganisationNode
	var err error

	if permission == "create_project" {
		user, ok := httputil.GetUserFromContext(r.Context())
		if !ok {
			httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		nodes, err = h.svc.SearchForProjectCreation(r.Context(), query, user.Person.ID)
	} else {
		nodes, err = h.svc.Search(r.Context(), query)
	}

	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.OrganisationResponse](nodes))
}

// SearchROR godoc
// @Summary Search ROR
// @Description Search for organizations in ROR
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {array} RORItem
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/ror/search [get]
func (h *Handler) SearchROR(w http.ResponseWriter, r *http.Request) {
	_, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		_ = httputil.WriteJSON(w, http.StatusOK, []RORItem{})
		return
	}

	items, err := SearchROR(query)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, items)
}

// CreateProjectRole godoc
// @Summary Create a project role for an organisation
// @Description Creates a new custom project role defined at this organisation node
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param body body dto.ProjectRoleCreateRequest true "Create role request"
// @Success 200 {object} dto.ProjectRoleResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id / invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/roles [post]
func (h *Handler) CreateProjectRole(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		logrus.Debugf("CreateProjectRole: invalid id param: %v", err)
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	hasAccess, err := h.rbac.HasPermission(r.Context(), user.Person.ID, id, orgrole.PermissionManageProjectRoles)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !hasAccess {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	var req dto.ProjectRoleCreateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Key == "" || req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "key and name are required", nil)
		return
	}

	role, err := h.roleSvc.Create(r.Context(), req.Key, req.Name, id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dto.ProjectRoleResponse{
		ID:   role.ID,
		Key:  role.Key,
		Name: role.Name,
	})
}

// DeleteProjectRole godoc
// @Summary Delete a project role
// @Description Deletes a custom project role defined at this organisation node
// @Tags organisation
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param roleId path string true "Role ID"
// @Success 204 "no content"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/roles/{roleId} [delete]
func (h *Handler) DeleteProjectRole(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	roleID, err := httputil.ParseUUIDParam(r, "roleId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid role id", nil)
		return
	}

	hasAccess, err := h.rbac.HasPermission(r.Context(), user.Person.ID, orgID, orgrole.PermissionManageProjectRoles)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !hasAccess {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	err = h.roleSvc.Delete(r.Context(), roleID, orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListProjectRoles godoc
// @Summary List project roles defined at an organisation node
// @Description Returns roles that are available to projects in this organisation (including inherited ones)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Success 200 {array} dto.ProjectRoleResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/roles [get]
func (h *Handler) ListProjectRoles(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	roles, err := h.roleSvc.ListAvailableForNode(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.ProjectRoleResponse](roles))
}

// UpdateProjectRole godoc
// @Summary Update a project role's allowed event types
// @Description Updates which event types a project role can use (EBAC)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param roleId path string true "Role ID"
// @Param body body dto.ProjectRoleUpdateRequest true "Update request"
// @Success 200 {object} dto.ProjectRoleResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/roles/{roleId} [patch]
func (h *Handler) UpdateProjectRole(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	roleID, err := httputil.ParseUUIDParam(r, "roleId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid role id", nil)
		return
	}

	hasAccess, err := h.rbac.HasPermission(r.Context(), user.Person.ID, orgID, orgrole.PermissionManageProjectRoles)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !hasAccess {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	var req dto.ProjectRoleUpdateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	role, err := h.roleSvc.UpdateAllowedEventTypes(r.Context(), roleID, req.AllowedEventTypes)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.ProjectRoleResponse](*role))
}

// CreateCustomField godoc
// @Summary Create a custom field definition for an organisation
// @Description Creates a new custom field definition at this organisation node
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param body body dto.CustomFieldDefinitionCreateRequest true "Create custom field definition request"
// @Success 200 {object} dto.CustomFieldDefinitionResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id / invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/custom-fields [post]
func (h *Handler) CreateCustomField(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	hasAccess, err := h.rbac.HasPermission(r.Context(), user.Person.ID, id, orgrole.PermissionManageCustomFields)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !hasAccess {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	var req dto.CustomFieldDefinitionCreateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" || req.Type == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name and type are required", nil)
		return
	}

	fd, err := h.customFieldSvc.Create(r.Context(), id, req.Name, entities.CustomFieldType(req.Type), entities.CustomFieldCategory(req.Category), req.Description, req.ValidationRegex, req.ExampleValue, req.Required)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.CustomFieldDefinitionResponse](*fd))
}

// DeleteCustomField godoc
// @Summary Delete a custom field definition
// @Description Deletes a custom field definition at this organisation node
// @Tags organisation
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param fieldId path string true "Field Definition ID"
// @Success 204 "no content"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/custom-fields/{fieldId} [delete]
func (h *Handler) DeleteCustomField(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	fieldID, err := httputil.ParseUUIDParam(r, "fieldId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid field id", nil)
		return
	}

	hasAccess, err := h.rbac.HasPermission(r.Context(), user.Person.ID, orgID, orgrole.PermissionManageCustomFields)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !hasAccess {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	err = h.customFieldSvc.Delete(r.Context(), fieldID, orgID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCustomFields godoc
// @Summary List custom field definitions at an organisation node
// @Description Returns definitions that are available to projects in this organisation (including inherited ones)
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param category query string false "Filter by category (PROJECT, PERSON)"
// @Success 200 {array} dto.CustomFieldDefinitionResponse
// @Failure 401 {string} string "unauthorized"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/custom-fields [get]
func (h *Handler) ListCustomFields(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var category *entities.CustomFieldCategory
	if c := r.URL.Query().Get("category"); c != "" {
		val := entities.CustomFieldCategory(c)
		category = &val
	}

	defs, err := h.customFieldSvc.ListAvailableForNode(r.Context(), id, category)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOs[dto.CustomFieldDefinitionResponse](defs))
}

// UpdateMemberCustomFields godoc
// @Summary Update custom fields for a member in an organisation context
// @Description Updates the custom field values for a specific person within the context of an organisation
// @Tags organisation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID"
// @Param personId path string true "Person ID"
// @Param body body dto.MemberCustomFieldUpdateValues true "Custom field values"
// @Success 204 "no content"
// @Failure 401 {string} string "unauthorized"
// @Failure 403 {string} string "forbidden"
// @Failure 400 {string} string "invalid id / invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /organisation-nodes/{id}/members/{personId}/custom-fields [put]
func (h *Handler) UpdateMemberCustomFields(w http.ResponseWriter, r *http.Request) {
	user, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	orgID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid org id", nil)
		return
	}

	personID, err := httputil.ParseUUIDParam(r, "personId")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid person id", nil)
		return
	}

	// Check Admin Access (Manage Members)
	hasAccess, err := h.rbac.HasPermission(r.Context(), user.Person.ID, orgID, orgrole.PermissionManageMembers)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if !hasAccess {
		httputil.WriteError(w, r, http.StatusForbidden, "forbidden", nil)
		return
	}

	var req dto.MemberCustomFieldUpdateValues
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	err = h.svc.UpdateMemberCustomFields(r.Context(), orgID, personID, req.Values)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
