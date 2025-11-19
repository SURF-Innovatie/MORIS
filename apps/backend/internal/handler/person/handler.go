package person

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/person"
)

type Handler struct {
	svc personsvc.Service
}

func NewHandler(s personsvc.Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	p, err := h.svc.Create(r.Context(), entities.Person{
		Name:       req.Name,
		GivenName:  req.GivenName,
		FamilyName: req.FamilyName,
		Email:      req.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, PersonResponse{
		ID:         p.Id.String(),
		Name:       p.Name,
		GivenName:  p.GivenName,
		FamilyName: p.FamilyName,
		Email:      p.Email,
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
	writeJSON(w, PersonResponse{
		ID:         p.Id.String(),
		Name:       p.Name,
		GivenName:  p.GivenName,
		FamilyName: p.FamilyName,
		Email:      p.Email,
	})
}

// List / Update implementations omitted for brevity; same pattern.

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
