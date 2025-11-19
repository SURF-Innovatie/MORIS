package custom

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/user"
)

type Handler struct {
	userService user.Service
}

func NewHandler(userService user.Service) *Handler {
	return &Handler{userService: userService}
}

// Status (GET /status)
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Login (POST /login) – later you can inject an AuthService as well.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := LoginResponse{Token: "supersecrettoken-dummy-jwt"}
	_ = json.NewEncoder(w).Encode(resp)
}

// Profile (GET /profile)
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userObj, ok := auth.GetUserFromContext(r.Context())
	if !ok || userObj == nil {
		http.Error(w, "User not authenticated or found in context", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userObj)
}

// TotalUserCount (GET /users/count)
func (h *Handler) TotalUserCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.userService.TotalUserCount(r.Context())
	if err != nil {
		http.Error(w, "failed to get user count: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := TotalUsersResponse{TotalUsers: count}
	_ = json.NewEncoder(w).Encode(resp)
}

// AdminUserList (GET /admin/users/list)
func (h *Handler) AdminUserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Admin-only user list!", "users": [{"id":1,"name":"Admin User"}]}`))
}
