package orcid

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	orcidService orcid.Service
}

func NewHandler(orcidService orcid.Service) *Handler {
	return &Handler{
		orcidService: orcidService,
	}
}

func MountRoutes(r chi.Router, h *Handler) {
	r.Get("/orcid/search", h.Search)
}

// Search godoc
// @Summary Search for people in ORCID
// @Description Searches for people in the ORCID public registry
// @Tags orcid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {array} dto.OrcidPerson
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /orcid/search [get]
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "query parameter 'q' is required", nil)
		return
	}

	results, err := h.orcidService.Search(r.Context(), query)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	out := make([]dto.OrcidPerson, 0, len(results))
	for _, p := range results {
		out = append(out, dto.OrcidPerson{
			FirstName:  p.FirstName,
			LastName:   p.LastName,
			CreditName: p.CreditName,
			Biography:  p.Biography,
			ORCID:      p.ORCID,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, results)
}
