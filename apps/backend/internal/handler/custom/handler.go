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

// Status godoc
// @Summary Check API status
// @Description Returns the current status of the API
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /status [get]
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Health godoc
// @Summary Health check
// @Description Returns the health status of the API
// @Tags status
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with the provided credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} auth.BackendError "Invalid request body or missing fields"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /register [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} auth.BackendError "Invalid request body"
// @Failure 401 {object} auth.BackendError "Invalid credentials"
// @Router /login [post]
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

// Profile godoc
// @Summary Get user profile
// @Description Returns the authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.AuthenticatedUser
// @Failure 401 {object} auth.BackendError "User not authenticated"
// @Router /profile [get]
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

// TotalUserCount godoc
// @Summary Get total user count
// @Description Returns the total number of registered users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TotalUsersResponse
// @Failure 401 {object} auth.BackendError "User not authenticated"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /users/count [get]
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

// AdminUserList godoc
// @Summary Get all users (Admin only)
// @Description Returns a list of all users - requires admin role
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "Admin user list"
// @Failure 401 {object} auth.BackendError "User not authenticated"
// @Failure 403 {object} auth.BackendError "Insufficient permissions"
// @Router /admin/users/list [get]
func (h *Handler) AdminUserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Admin-only user list!", "users": [{"id":1,"name":"Admin User"}]}`))
}
