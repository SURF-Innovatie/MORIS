package organisation

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	organisationsvc "github.com/SURF-Innovatie/MORIS/internal/organisation"
)

type Handler struct {
	svc organisationsvc.Service
}

func NewHandler(s organisationsvc.Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req organisationdto.CreateRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	p, err := h.svc.Create(r.Context(), entities.Organisation{
		Name: req.Name,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.FromEntity(*p))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.FromEntity(*p))
}

// List / Update implementations omitted for brevity; same pattern.
