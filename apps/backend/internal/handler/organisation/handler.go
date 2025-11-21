package organisation

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/organisationdto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
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

	writeJSON(w, organisationdto.Response{
		ID:   p.Id,
		Name: p.Name,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, organisationdto.Response{
		ID:   p.Id,
		Name: p.Name,
	})
}

// List / Update implementations omitted for brevity; same pattern.

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
