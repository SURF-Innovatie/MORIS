package doi

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/SURF-Innovatie/MORIS/internal/api/dto"
	app "github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	svc app.Service
}

func NewHandler(svc app.Service) *Handler {
	return &Handler{svc: svc}
}

// Resolve godoc
// @Summary Resolve a DOI
// @Description Resolves a DOI to product data using content negotiation with the registration agency
// @Tags doi
// @Accept json
// @Produce json
// @Param doi query string true "DOI"
// @Success 200 {object} dto.Work
// @Failure 400 {object} httputil.BackendError "doi is required"
// @Failure 404 {object} httputil.BackendError "doi not found"
// @Failure 500 {object} httputil.BackendError "internal server error"
// @Router /doi/resolve [get]
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	doi := r.URL.Query().Get("doi")
	if doi == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "doi is required", nil)
		return
	}

	work, err := h.svc.Resolve(r.Context(), doi)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			httputil.WriteError(w, r, http.StatusNotFound, "doi not found", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, *work)
}

func MountRoutes(r chi.Router, h *Handler) {
	r.Route("/doi", func(r chi.Router) {
		r.Get("/resolve", h.Resolve)
	})
}
