package adapter

import (
	"fmt"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/domain/identity"
	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
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
// @Description Triggers a project export to the specified sink. For file-based sinks, returns a file download.
// @Tags projects
// @Param id path string true "Project ID"
// @Param sink path string true "Sink Name"
// @Produce json
// @Produce octet-stream
// @Security BearerAuth
// @Success 200 {object} map[string]string "Export successful or file download"
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

	// Get project details
	projDetails, err := h.projectSvc.GetProject(r.Context(), projectID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, "project not found", nil)
		return
	}

	// Load event stream
	evts, err := h.projectSvc.GetEvents(r.Context(), projectID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to load project events", nil)
		return
	}

	// Resolve members
	members := lo.Map(projDetails.Members, func(m project.MemberDetail, _ int) identity.Person {
		return m.Person
	})

	// Create context
	pc := adapter.ProjectContext{
		ProjectID: projectID,
		Events:    evts,
		Project:   &projDetails.Project,
		Members:   members,
		OrgNode:   &projDetails.OwningOrgNode,
	}

	// Connect
	if err := sink.Connect(r.Context()); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to connect to sink", nil)
		return
	}

	// Check if file-based or API-based
	outputInfo := sink.OutputInfo()
	if outputInfo.Type == adapter.TransferTypeFile {
		// File-based export - stream file to client
		result, err := sink.ExportProjectData(r.Context(), pc)
		if err != nil {
			httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		w.Header().Set("Content-Type", result.MimeType)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", result.Filename))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result.Data)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(result.Data)
		return
	}

	// API-based export
	if err := sink.PushProject(r.Context(), pc); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "Export successful"})
}
