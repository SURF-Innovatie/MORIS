package custom

import (
	"encoding/json"
	"net/http"
	"strings"
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
	userCtx, ok := auth.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "User not authenticated or found in context",
		})
		return
	}

	// Fetch fresh user data from database
	freshUser, err := h.authService.GetUserByID(r.Context(), userCtx.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to fetch user profile",
		})
		return
	}

	authUser := &auth.AuthenticatedUser{
		ID:      freshUser.ID,
		Email:   freshUser.Email,
		OrcidID: freshUser.OrcidID,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authUser)
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

// GetORCIDAuthURL godoc
// @Summary Get ORCID authorization URL
// @Description Returns the URL to redirect the user to for ORCID authentication
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ORCIDAuthURLResponse
// @Failure 401 {object} auth.BackendError "User not authenticated"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /auth/orcid/url [get]
func (h *Handler) GetORCIDAuthURL(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	url, err := h.authService.GenerateORCIDAuthURL(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusInternalServerError,
			Status:  "Internal Server Error",
			Message: "Failed to generate ORCID auth URL: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := ORCIDAuthURLResponse{URL: url}
	_ = json.NewEncoder(w).Encode(resp)
}

// LinkORCID godoc
// @Summary Link ORCID ID
// @Description Links an ORCID ID to the authenticated user's account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LinkORCIDRequest true "ORCID link request"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} auth.BackendError "Invalid request"
// @Failure 401 {object} auth.BackendError "User not authenticated"
// @Failure 409 {object} auth.BackendError "ORCID ID already linked"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /auth/orcid/link [post]
func (h *Handler) LinkORCID(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req LinkORCIDRequest
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

	if req.Code == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Authorization code is required",
		})
		return
	}

	err := h.authService.LinkORCID(r.Context(), user.ID, req.Code)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already linked") {
			statusCode = http.StatusConflict
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    statusCode,
			Status:  http.StatusText(statusCode),
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// UnlinkORCID godoc
// @Summary Unlink ORCID ID
// @Description Unlinks the ORCID ID from the authenticated user's account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} auth.BackendError "User not authenticated"
// @Failure 500 {object} auth.BackendError "Internal server error"
// @Router /auth/orcid/unlink [post]
func (h *Handler) UnlinkORCID(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(auth.BackendError{
			Code:    http.StatusUnauthorized,
			Status:  "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	err := h.authService.UnlinkORCID(r.Context(), user.ID)
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
	resp := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}
