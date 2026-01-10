package zenodo

import (
	"io"
	"net/http"
	"strconv"

	"github.com/SURF-Innovatie/MORIS/external/zenodo"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
)

// Handler provides HTTP handlers for Zenodo integration
type Handler struct {
	zenodoService zenodo.Service
	currentUser   appauth.CurrentUserProvider
}

// NewHandler creates a new Zenodo handler
func NewHandler(zenodoService zenodo.Service, currentUser appauth.CurrentUserProvider) *Handler {
	return &Handler{
		zenodoService: zenodoService,
		currentUser:   currentUser,
	}
}

// MountRoutes mounts the Zenodo routes on the router
func MountRoutes(r chi.Router, h *Handler) {
	r.Get("/zenodo/auth-url", h.GetAuthURL)
	r.Post("/zenodo/link", h.Link)
	r.Delete("/zenodo/unlink", h.Unlink)
	r.Get("/zenodo/status", h.GetStatus)
	r.Get("/zenodo/depositions", h.ListDepositions)
	r.Post("/zenodo/depositions", h.CreateDeposition)
	r.Get("/zenodo/depositions/{id}", h.GetDeposition)
	r.Put("/zenodo/depositions/{id}", h.UpdateDeposition)
	r.Delete("/zenodo/depositions/{id}", h.DeleteDeposition)
	r.Post("/zenodo/depositions/{id}/files", h.UploadFile)
	r.Post("/zenodo/depositions/{id}/publish", h.Publish)
}

// GetAuthURL godoc
// @Summary Get Zenodo OAuth authorization URL
// @Description Returns a URL to redirect the user to for Zenodo OAuth authorization
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "auth_url"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/auth-url [get]
func (h *Handler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	authURL, err := h.zenodoService.GetAuthURL(r.Context(), u.UserID())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"auth_url": authURL})
}

// Link godoc
// @Summary Link Zenodo account
// @Description Exchanges an authorization code for tokens and links the Zenodo account
// @Tags zenodo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body LinkRequest true "OAuth code"
// @Success 200 {object} map[string]string "status: linked"
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/link [post]
func (h *Handler) Link(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req LinkRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Code == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "code is required", nil)
		return
	}

	if err := h.zenodoService.Link(r.Context(), u.UserID(), req.Code); err != nil {
		if err == zenodo.ErrAlreadyLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account already linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "linked"})
}

// LinkRequest is the request body for linking a Zenodo account
type LinkRequest struct {
	Code string `json:"code"`
}

// Unlink godoc
// @Summary Unlink Zenodo account
// @Description Removes the Zenodo account link from the current user
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "status: unlinked"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/unlink [delete]
func (h *Handler) Unlink(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if err := h.zenodoService.Unlink(r.Context(), u.UserID()); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "unlinked"})
}

// GetStatus godoc
// @Summary Get Zenodo link status
// @Description Returns whether the current user has a linked Zenodo account
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatusResponse
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/status [get]
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	linked, err := h.zenodoService.IsLinked(r.Context(), u.UserID())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, StatusResponse{Linked: linked})
}

// StatusResponse is the response for the status endpoint
type StatusResponse struct {
	Linked bool `json:"linked"`
}

// ListDepositions godoc
// @Summary List Zenodo depositions
// @Description Returns all depositions for the current user
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Success 200 {array} zenodo.Deposition
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions [get]
func (h *Handler) ListDepositions(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	depositions, err := h.zenodoService.ListDepositions(r.Context(), u.UserID())
	if err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, depositions)
}

// CreateDeposition godoc
// @Summary Create a new Zenodo deposition
// @Description Creates a new empty deposition
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Success 201 {object} zenodo.Deposition
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions [post]
func (h *Handler) CreateDeposition(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	deposition, err := h.zenodoService.CreateDeposition(r.Context(), u.UserID())
	if err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, deposition)
}

// GetDeposition godoc
// @Summary Get a Zenodo deposition
// @Description Returns a single deposition by ID
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deposition ID"
// @Success 200 {object} zenodo.Deposition
// @Failure 400 {object} httputil.BackendError "Invalid ID"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions/{id} [get]
func (h *Handler) GetDeposition(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid deposition ID", nil)
		return
	}

	deposition, err := h.zenodoService.GetDeposition(r.Context(), u.UserID(), id)
	if err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, deposition)
}

// UpdateDeposition godoc
// @Summary Update a Zenodo deposition
// @Description Updates the metadata of a deposition
// @Tags zenodo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deposition ID"
// @Param body body zenodo.DepositionMetadata true "Deposition metadata"
// @Success 200 {object} zenodo.Deposition
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions/{id} [put]
func (h *Handler) UpdateDeposition(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid deposition ID", nil)
		return
	}

	var metadata zenodo.DepositionMetadata
	if !httputil.ReadJSON(w, r, &metadata) {
		return
	}

	deposition, err := h.zenodoService.UpdateDeposition(r.Context(), u.UserID(), id, &metadata)
	if err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, deposition)
}

// DeleteDeposition godoc
// @Summary Delete a Zenodo deposition
// @Description Deletes an unpublished deposition
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deposition ID"
// @Success 204 "No Content"
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions/{id} [delete]
func (h *Handler) DeleteDeposition(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid deposition ID", nil)
		return
	}

	if err := h.zenodoService.DeleteDeposition(r.Context(), u.UserID(), id); err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UploadFile godoc
// @Summary Upload a file to a deposition
// @Description Uploads a file to an existing deposition
// @Tags zenodo
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deposition ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} zenodo.DepositionFile
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions/{id}/files [post]
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid deposition ID", nil)
		return
	}

	// Parse multipart form (max 50GB as per Zenodo limits)
	if err := r.ParseMultipartForm(50 << 30); err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "failed to parse form", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "file is required", nil)
		return
	}
	defer file.Close()

	// Create a reader that properly handles the file
	reader := io.Reader(file)

	depositionFile, err := h.zenodoService.UploadFile(r.Context(), u.UserID(), id, header.Filename, reader)
	if err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, depositionFile)
}

// Publish godoc
// @Summary Publish a deposition
// @Description Publishes a deposition, minting a DOI
// @Tags zenodo
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deposition ID"
// @Success 202 {object} zenodo.Deposition
// @Failure 400 {object} httputil.BackendError "Invalid request"
// @Failure 401 {object} httputil.BackendError "User not authenticated"
// @Failure 500 {object} httputil.BackendError "Internal server error"
// @Router /zenodo/depositions/{id}/publish [post]
func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid deposition ID", nil)
		return
	}

	deposition, err := h.zenodoService.Publish(r.Context(), u.UserID(), id)
	if err != nil {
		if err == zenodo.ErrNotLinked {
			httputil.WriteError(w, r, http.StatusBadRequest, "Zenodo account not linked", nil)
			return
		}
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httputil.WriteJSON(w, http.StatusAccepted, deposition)
}
