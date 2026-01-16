package crossref

import (
	"errors"
	"net/http"

	_ "github.com/SURF-Innovatie/MORIS/internal/api/dto"
	app "github.com/SURF-Innovatie/MORIS/internal/app/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	svc app.Service
}

func NewHandler(svc app.Service) *Handler {
	return &Handler{svc: svc}
}

// GetWork godoc
// @Summary Get a work by DOI
// @Description Retrieves a single work from Crossref by its DOI
// @Tags crossref
// @Accept json
// @Produce json
// @Param doi query string true "DOI"
// @Success 200 {object} dto.Work
// @Failure 400 {object} httputil.BackendError "doi is required"
// @Failure 404 {object} httputil.BackendError "work not found"
// @Failure 500 {object} httputil.BackendError "internal server error"
// @Router /crossref/works [get]
func (h *Handler) GetWork(w http.ResponseWriter, r *http.Request) {
	doi := r.URL.Query().Get("doi")
	if doi == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "doi is required", nil)
		return
	}

	work, err := h.svc.GetWork(r.Context(), doi)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httputil.WriteError(w, r, http.StatusNotFound, "work not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, *work)
}
