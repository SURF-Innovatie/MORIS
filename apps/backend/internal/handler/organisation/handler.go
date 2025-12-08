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
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	p, err := h.svc.Create(r.Context(), entities.Organisation{
		Name: req.Name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.Response{
		ID:   p.Id,
		Name: p.Name,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = httputil.WriteJSON(w, http.StatusOK, organisationdto.Response{
		ID:   p.Id,
		Name: p.Name,
	})
}

// List / Update implementations omitted for brevity; same pattern.
