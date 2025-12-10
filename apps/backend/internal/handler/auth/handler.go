package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/api/authdto"
	"github.com/SURF-Innovatie/MORIS/internal/api/userdto"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/auth"
	"github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/SURF-Innovatie/MORIS/internal/user"
)

type Handler struct {
	userService  user.Service
	authService  coreauth.Service
	orcidService orcid.Service
}

func NewHandler(userService user.Service, authService coreauth.Service, orcidService orcid.Service) *Handler {
	return &Handler{
		userService:  userService,
		authService:  authService,
		orcidService: orcidService,
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body authdto.LoginRequest true "User login credentials"
// @Success 200 {object} authdto.LoginResponse
// @Failure 400 {object} httputil.BackendError "Invalid request body"
// @Failure 401 {object} httputil.BackendError "Invalid credentials"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authdto.LoginRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	token, authUser, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		_ = httputil.WriteError(w, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	resp := authdto.LoginResponse{
		Token: token,
		User:  authUser,
	}
	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}

// Profile godoc
// @Summary Get user profile
// @Description Returns the authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userdto.Response
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Router /profile [get]
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx == nil {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated or found in context", nil)
		return
	}

	// Fetch fresh user data from database
	freshUser, err := h.userService.GetAccount(r.Context(), userCtx.User.ID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "Failed to fetch user profile", nil)
		return
	}

	dto := userdto.Response{
		ID:          freshUser.User.ID,
		PersonID:    freshUser.User.PersonID,
		ORCiD:       freshUser.Person.ORCiD,
		Name:        freshUser.Person.Name,
		GivenName:   freshUser.Person.GivenName,
		FamilyName:  freshUser.Person.FamilyName,
		Email:       freshUser.Person.Email,
		AvatarURL:   freshUser.Person.AvatarUrl,
		Description: freshUser.Person.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto)
}

// GetORCIDAuthURL godoc
// @Summary Get ORCID authorization URL
// @Description Returns the URL to redirect the user to for ORCID authentication
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authdto.ORCIDAuthURLResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/orcid/url [get]
func (h *Handler) GetORCIDAuthURL(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	url, err := h.orcidService.GetAuthURL(r.Context(), u.User.ID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, orcid.ErrUnauthenticated) {
			statusCode = http.StatusUnauthorized
		}

		httputil.WriteError(w, statusCode, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := authdto.ORCIDAuthURLResponse{URL: url}
	_ = json.NewEncoder(w).Encode(resp)
}

// LinkORCID godoc
// @Summary Link ORCID ID
// @Description Links an ORCID ID to the authenticated user's account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body authdto.LinkORCIDRequest true "ORCID link request"
// @Success 200 {object} httputil.StatusResponse
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 409 {object} httputil.BackendError "ORCID ID already linked"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /auth/orcid/link [post]
func (h *Handler) LinkORCID(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req authdto.LinkORCIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "Invalid request body", nil)
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

		httputil.WriteError(w, statusCode, err.Error(), nil)
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
	u, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := h.orcidService.Unlink(r.Context(), u.User.ID); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = httputil.WriteStatus(w)
}
