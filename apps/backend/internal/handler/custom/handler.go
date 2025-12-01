package custom

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	"github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/user"
)

type Handler struct {
	userService  user.Service
	authService  auth.Service
	orcidService orcid.Service
}

func NewHandler(userService user.Service, authService auth.Service, orcidService orcid.Service) *Handler {
	return &Handler{
		userService:  userService,
		authService:  authService,
		orcidService: orcidService,
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
// @Success 201 {object} userdto.Response
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

	var usrReq userdto.Request
	usr, err := h.authService.Register(r.Context(), usrReq)
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

	_ = json.NewEncoder(w).Encode(usr)
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
		User:  authUser,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Profile godoc
// @Summary Get user profile
// @Description Returns the authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userdto.Response
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
	freshUser, err := h.userService.GetAccount(r.Context(), userCtx.User.ID)
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

	dto := userdto.Response{
		ID:         freshUser.User.ID,
		PersonID:   freshUser.User.PersonID,
		ORCiD:      freshUser.Person.ORCiD,
		Name:       freshUser.Person.Name,
		GivenName:  freshUser.Person.GivenName,
		FamilyName: freshUser.Person.FamilyName,
		Email:      freshUser.Person.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto)
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
	u, ok := auth.GetUserFromContext(r.Context())
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

	url, err := h.orcidService.GetAuthURL(r.Context(), u.User.ID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, orcid.ErrUnauthenticated) {
			statusCode = http.StatusUnauthorized
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
	u, ok := auth.GetUserFromContext(r.Context())
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

	err := h.orcidService.Link(r.Context(), u.User.ID, req.Code)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch {
		case errors.Is(err, orcid.ErrMissingCode):
			statusCode = http.StatusBadRequest
		case errors.Is(err, orcid.ErrAlreadyLinked):
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
	u, ok := auth.GetUserFromContext(r.Context())
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

	if err := h.orcidService.Unlink(r.Context(), u.User.ID); err != nil {
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
