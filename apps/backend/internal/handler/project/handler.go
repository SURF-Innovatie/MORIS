package project

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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

func (h *Handler) StartProject(w http.ResponseWriter, r *http.Request) {
	req := StartRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		http.Error(w, "invalid startDate", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		http.Error(w, "invalid endDate", http.StatusBadRequest)
		return
	}

	params := project.StartProjectParams{
		Title:        req.Title,
		Description:  req.Description,
		Organisation: req.Organisation,
		StartDate:    start,
		EndDate:      end,
	}

	proj, err := h.svc.StartProject(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

func (h *Handler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	proj, err := h.svc.GetAllProjects(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

type personRequest struct {
	Name string `json:"name"`
}

// AddPerson POST /projects/{id}/person/add
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

	person := entities.NewPerson(req.Name)

	proj, err := h.svc.AddPerson(r.Context(), id, person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}

// RemovePerson DELETE /projects/{id}/person/{personId}
func (h *Handler) RemovePerson(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	personIdStr := chi.URLParam(r, "personId")
	personId, err := uuid.Parse(personIdStr)
	if err != nil {
		http.Error(w, "invalid personId", http.StatusBadRequest)
		return
	}

	proj, err := h.svc.RemovePerson(r.Context(), id, personId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proj)
}
