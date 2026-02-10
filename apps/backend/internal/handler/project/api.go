package project

import (
	"github.com/SURF-Innovatie/MORIS/internal/handler/bulkimport"
	"github.com/go-chi/chi/v5"
)

func MountProjectRoutes(r chi.Router, h *Handler, bulk *bulkimport.Handler) {
	r.Get("/", h.GetAllProjects)
	r.Get("/{id}", h.GetProject)
	r.Get("/{id}/changelog", h.GetChangelog)
	r.Get("/{id}/pending-events", h.GetPendingEvents)
	r.Get("/{id}/allowed-events", h.GetAllowedEvents)
	r.Get("/{id}/custom-fields", h.ListAvailableCustomFields)
	r.Post("/{id}/bulk-import", bulk.BulkImportIntoProject)
}
