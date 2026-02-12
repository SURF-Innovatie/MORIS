package bulkimport

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	bulkimport2 "github.com/SURF-Innovatie/MORIS/internal/app/project/bulkimport"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	svc         bulkimport2.Service
	currentUser appauth.CurrentUserProvider
}

func NewHandler(svc bulkimport2.Service, cu appauth.CurrentUserProvider) *Handler {
	return &Handler{svc: svc, currentUser: cu}
}

// BulkImportIntoProject godoc
// @Summary Bulk import products by DOI into a project
// @Description Resolves DOI metadata, creates products and adds them to the specified project.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID (UUID)"
// @Param request body dto.BulkImportRequest true "List of DOIs"
// @Success 200 {object} dto.BulkImportResponse
// @Failure 400 {string} string "invalid id or body"
// @Failure 401 {string} string "unauthorized"
// @Failure 500 {string} string "internal server error"
// @Router /projects/{id}/bulk-import [post]
func (h *Handler) BulkImportIntoProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.BulkImportRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if len(req.Dois) == 0 {
		httputil.WriteError(w, r, http.StatusBadRequest, "dois is required", nil)
		return
	}

	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	entries := make([]bulkimport2.Entry, 0, len(req.Dois))
	for _, d := range req.Dois {
		entries = append(entries, bulkimport2.Entry{DOI: d})
	}

	result, err := h.svc.BulkImport(r.Context(), u.UserID, u.PersonID, projectID, entries)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	resp := dto.BulkImportResponse{
		ProjectID:       result.ProjectID,
		CreatedProducts: result.CreatedProducts,
		Items:           make([]dto.BulkImportItemResult, 0, len(result.Items)),
	}

	for _, it := range result.Items {
		item := dto.BulkImportItemResult{
			DOI: it.DOI,
		}

		if it.Error != "" {
			e := it.Error
			item.Error = &e
		}

		if it.Product != nil {
			pdto := transform.ToDTOItem[dto.ProductResponse](*it.Product)
			item.Product = &pdto
		}

		resp.Items = append(resp.Items, item)
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}
