package event

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/event"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc event.Service
}

func NewHandler(svc event.Service) *Handler {
	return &Handler{svc: svc}
}

// ApproveEvent godoc
// @Summary Approve an event
// @Description Approves a pending event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid event id"
// @Failure 500 {string} string "internal server error"
// @Router /events/{id}/approve [post]
func (h *Handler) ApproveEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.svc.ApproveEvent(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// RejectEvent godoc
// @Summary Reject an event
// @Description Rejects a pending event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid event id"
// @Failure 500 {string} string "internal server error"
// @Router /events/{id}/reject [post]
func (h *Handler) RejectEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.svc.RejectEvent(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
