package projecthandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/SURF-Innovatie/MORIS/internal/project"
)

type Handler struct {
	svc project.Service
}

func NewHandler(svc project.Service) *Handler {
	return &Handler{svc: svc}
}

// GET /projects/{id}
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}
	proj, err := h.svc.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

type personRequest struct {
	Name string `json:"name"`
}

// POST /projects/{id}/person/add
func (h *Handler) AddPerson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var req personRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.AddPerson(r.Context(), id, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

// POST /projects/{id}/person/remove
func (h *Handler) RemovePerson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var req personRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.RemovePerson(r.Context(), id, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}
