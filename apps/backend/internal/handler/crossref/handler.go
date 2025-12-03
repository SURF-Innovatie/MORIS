package crossref

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/infra/external/crossref"
)

type Handler struct {
	svc crossref.Service
}

func NewHandler(svc crossref.Service) *Handler {
	return &Handler{svc: svc}
}

// GetWork godoc
// @Summary Get a work by DOI
// @Description Retrieves a single work from Crossref by its DOI
// @Tags crossref
// @Accept json
// @Produce json
// @Param doi query string true "DOI"
// @Success 200 {object} crossref.Work
// @Failure 400 {string} string "invalid doi"
// @Failure 404 {string} string "work not found"
// @Failure 500 {string} string "internal server error"
// @Router /crossref/works [get]
func (h *Handler) GetWork(w http.ResponseWriter, r *http.Request) {
	doi := r.URL.Query().Get("doi")
	if doi == "" {
		http.Error(w, "doi is required", http.StatusBadRequest)
		return
	}

	// No need to unescape path, query params are decoded by net/http
	work, err := h.svc.GetWork(r.Context(), doi)
	if err != nil {
		// TODO: Handle 404 specifically if the service returns a specific error for it
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(work); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
