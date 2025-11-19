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
	authService auth.Service
}

func NewHandler(userService user.Service, authService auth.Service) *Handler {
	return &Handler{
		userService: userService,
		authService: authService,
	}
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

// Health (GET /health)
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Register (POST /register)
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Name, email, and password are required",
		})
		return
	}

	usr, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := RegisterResponse{
		ID:    usr.ID,
		Email: usr.Email,
		Name:  usr.Name,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Login (POST /login)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
		})
		return
	}

	token, authUser, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "Invalid credentials",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := LoginResponse{
		Token: token,
	}
	resp.User.ID = authUser.ID
	resp.User.Email = authUser.Email
	resp.User.Roles = authUser.Roles
	_ = json.NewEncoder(w).Encode(resp)
}

// Profile (GET /profile)
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userObj, ok := auth.GetUserFromContext(r.Context())
	if !ok || userObj == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "User not authenticated or found in context",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(userObj)
}

// TotalUserCount (GET /users/count)
func (h *Handler) TotalUserCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.userService.TotalUserCount(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "failed to get user count: " + err.Error(),
		})
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
