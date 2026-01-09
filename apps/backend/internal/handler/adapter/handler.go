package adapter

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

type Handler struct {
	registry   *adapter.Registry
	projectSvc queries.Service
}

func NewHandler(registry *adapter.Registry, projectSvc queries.Service) *Handler {
	return &Handler{
		registry:   registry,
		projectSvc: projectSvc,
	}
}

// ListAdapters godoc
// @Summary List all registered adapters
// @Description Returns a list of all available source and sink adapters with their metadata
// @Tags adapters
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AdapterListResponse
// @Router /adapters [get]
func (h *Handler) ListAdapters(w http.ResponseWriter, r *http.Request) {
	sources := lo.Map(h.registry.ListSources(), func(name string, _ int) dto.SourceInfoResponse {
		s, _ := h.registry.GetSource(name)
		return dto.SourceInfoResponse{
			AdapterInfoResponse: dto.AdapterInfoResponse{
				Name:           s.Name(),
				DisplayName:    s.DisplayName(),
				SupportedTypes: s.SupportedTypes(),
			},
			Input: s.InputInfo(),
		}
	})

	sinks := lo.Map(h.registry.ListSinks(), func(name string, _ int) dto.SinkInfoResponse {
		s, _ := h.registry.GetSink(name)
		return dto.SinkInfoResponse{
			AdapterInfoResponse: dto.AdapterInfoResponse{
				Name:           s.Name(),
				DisplayName:    s.DisplayName(),
				SupportedTypes: s.SupportedTypes(),
			},
			Output: s.OutputInfo(),
		}
	})

	_ = httputil.WriteJSON(w, http.StatusOK, dto.AdapterListResponse{
		Sources: sources,
		Sinks:   sinks,
	})
}

// ExportProject godoc
// @Summary Export a project using a sink adapter
// @Description Triggers a project export to the specified sink
// @Tags projects
// @Param id path string true "Project ID"
// @Param sink path string true "Sink Name"
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "Export successful"
// @Failure 400 {string} string "Invalid project id"
// @Failure 404 {string} string "Project or sink not found"
// @Failure 500 {string} string "Export failed"
// @Router /projects/{id}/export/{sink} [post]
func (h *Handler) ExportProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid project id", nil)
		return
	}

	sinkName := chi.URLParam(r, "sink")
	sink, ok := h.registry.GetSink(sinkName)
	if !ok {
		httputil.WriteError(w, r, http.StatusNotFound, "sink not found", nil)
		return
	}

	// 1. Get project data
	// Note: We need the full context (events, entities, org node)
	// This might require adding a specialized method to projectSvc or gathering it here
	projDetails, err := h.projectSvc.GetProject(r.Context(), projectID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "project not found", nil)
		return
	}

	// 2. Load event stream
	// This should be available via the project service or a dedicated event store service
	// For this example, let's assume we can get it from the service
	evts, err := h.projectSvc.GetEvents(r.Context(), projectID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to load project events", nil)
		return
	}

	// 3. Resolve members
	members := lo.Map(projDetails.Members, func(m entities.ProjectMemberDetail, _ int) entities.Person {
		return m.Person
	})

	// 4. Create context
	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events:    evts,
		Project:   &projDetails.Project,
		Members:   members,
		OrgNode:   &projDetails.OwningOrgNode,
	}

	// 5. Connect and push
	if err := sink.Connect(r.Context()); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to connect to sink", nil)
		return
	}

	if err := sink.PushProject(r.Context(), pc); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, "Export successful")
}
