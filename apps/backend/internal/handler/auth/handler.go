package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	"github.com/SURF-Innovatie/MORIS/internal/app/orcid"
	surfconextapp "github.com/SURF-Innovatie/MORIS/internal/app/surfconext"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

type Handler struct {
	userService  user.Service
	authService  coreauth.Service
	orcidService orcid.Service
	surfService  surfconextapp.Service
}

func NewHandler(userService user.Service, authService coreauth.Service, orcidService orcid.Service, surfService surfconextapp.Service) *Handler {
	return &Handler{
		userService:  userService,
		authService:  authService,
		orcidService: orcidService,
		surfService:  surfService,
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} httputil.BackendError "Invalid request body"
// @Failure 401 {object} httputil.BackendError "Invalid credentials"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	token, authUser, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		_ = httputil.WriteError(w, r, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	resp := dto.FromEntity(token, authUser)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// GetSurfconextAuthURL godoc
// @Summary Get SURFconext authorization URL
// @Description Returns the URL to redirect the user to for SURFconext (OIDC) login
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.SurfconextAuthURLResponse
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/surfconext/url [get]
func (h *Handler) GetSurfconextAuthURL(w http.ResponseWriter, r *http.Request) {
	url, err := h.surfService.AuthURL(r.Context())
	if err != nil {
		_ = httputil.WriteError(w, r, http.StatusInternalServerError, "Failed to create authorization URL", map[string]any{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.SurfconextAuthURLResponse{URL: url})
}

// LoginWithSurfconext godoc
// @Summary Login with SURFconext
// @Description Exchanges an OIDC authorization code for a MORIS JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SurfconextLoginRequest true "SURFconext login request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not allowed / not found"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/surfconext/login [post]
func (h *Handler) LoginWithSurfconext(w http.ResponseWriter, r *http.Request) {
	var req dto.SurfconextLoginRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	token, authUser, err := h.surfService.LoginWithCode(r.Context(), req.Code)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch {
		case errors.Is(err, surfconextapp.ErrMissingCode):
			statusCode = http.StatusBadRequest
		case errors.Is(err, surfconextapp.ErrNoEmail):
			statusCode = http.StatusUnauthorized
		default:
			// If the wrapped error contains "invalid credentials" from LoginByEmail
			statusCode = http.StatusUnauthorized
		}

		_ = httputil.WriteError(w, r, statusCode, err.Error(), nil)
		return
	}

	resp := dto.FromEntity(token, authUser)
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// Profile godoc
// @Summary Get user profile
// @Description Returns the authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Router /profile [get]
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := httputil.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "User not found in context", nil)
		return
	}

	// Fetch fresh user data from database
	freshUser, err := h.userService.GetAccount(r.Context(), userCtx.User.ID)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "Failed to fetch user profile", nil)
		return
	}

	dtoResp := transform.ToDTOItem[dto.UserResponse](freshUser)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dtoResp)
}

// GetORCIDAuthURL godoc
// @Summary Get ORCID authorization URL
// @Description Returns the URL to redirect the user to for ORCID authentication
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ORCIDAuthURLResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/orcid/url [get]
func (h *Handler) GetORCIDAuthURL(w http.ResponseWriter, r *http.Request) {
	u, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	url, err := h.orcidService.GetAuthURL(r.Context(), u.User.ID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, orcid.ErrUnauthenticated) {
			statusCode = http.StatusUnauthorized
		}

		httputil.WriteError(w, r, statusCode, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := dto.ORCIDAuthURLResponse{URL: url}
	_ = json.NewEncoder(w).Encode(resp)
}

// LinkORCID godoc
// @Summary Link ORCID ID
// @Description Links an ORCID ID to the authenticated user's account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LinkORCIDRequest true "ORCID link request"
// @Success 200 {object} httputil.StatusResponse
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 409 {object} httputil.BackendError "ORCID ID already linked"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/orcid/link [post]
func (h *Handler) LinkORCID(w http.ResponseWriter, r *http.Request) {
	u, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.LinkORCIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "Invalid request body", nil)
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

		httputil.WriteError(w, r, statusCode, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = httputil.WriteStatus(w)
}

// UnlinkORCID godoc
// @Summary Unlink ORCID ID
// @Description Unlinks the ORCID ID from the authenticated user's account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httputil.StatusResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/orcid/unlink [post]
func (h *Handler) UnlinkORCID(w http.ResponseWriter, r *http.Request) {
	u, ok := httputil.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, r, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := h.orcidService.Unlink(r.Context(), u.User.ID); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = httputil.WriteStatus(w)
}
