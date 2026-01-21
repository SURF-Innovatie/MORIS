package portfolio

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	appauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/portfolio"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/google/uuid"
)

type Handler struct {
	svc         portfolio.Service
	currentUser appauth.CurrentUserProvider
}

func NewHandler(svc portfolio.Service, cu appauth.CurrentUserProvider) *Handler {
	return &Handler{svc: svc, currentUser: cu}
}

// GetMe godoc
// @Summary Get portfolio settings for current user
// @Description Returns portfolio customizations such as pinned projects and deliverables
// @Tags portfolio
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.PortfolioResponse
// @Failure 500 {string} string "internal server error"
// @Router /portfolio/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	portfolioSettings, err := h.svc.GetForPerson(r.Context(), u.PersonID())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if portfolioSettings == nil {
		defaults := defaultPortfolio(u.PersonID())
		_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.PortfolioResponse](defaults))
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.PortfolioResponse](*portfolioSettings))
}

// UpdateMe godoc
// @Summary Update portfolio settings for current user
// @Description Updates portfolio customizations such as headline and pinned items
// @Tags portfolio
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param portfolio body dto.PortfolioRequest true "Portfolio settings"
// @Success 200 {object} dto.PortfolioResponse
// @Failure 400 {string} string "invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /portfolio/me [put]
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	var req dto.PortfolioRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	u, err := h.currentUser.Current(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	existing, err := h.svc.GetForPerson(r.Context(), u.PersonID())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	merged := defaultPortfolio(u.PersonID())
	if existing != nil {
		merged = *existing
	}

	if req.Headline != nil {
		merged.Headline = req.Headline
	}
	if req.Summary != nil {
		merged.Summary = req.Summary
	}
	if req.Website != nil {
		merged.Website = req.Website
	}
	if req.ShowEmail != nil {
		merged.ShowEmail = *req.ShowEmail
	}
	if req.ShowOrcid != nil {
		merged.ShowOrcid = *req.ShowOrcid
	}
	if req.PinnedProjectIDs != nil {
		merged.PinnedProjectIDs = req.PinnedProjectIDs
	}
	if req.PinnedProductIDs != nil {
		merged.PinnedProductIDs = req.PinnedProductIDs
	}

	if merged.PinnedProjectIDs == nil {
		merged.PinnedProjectIDs = []uuid.UUID{}
	}
	if merged.PinnedProductIDs == nil {
		merged.PinnedProductIDs = []uuid.UUID{}
	}

	updated, err := h.svc.UpdateForPerson(r.Context(), u.PersonID(), merged)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.PortfolioResponse](*updated))
}

func defaultPortfolio(personID uuid.UUID) entities.Portfolio {
	return entities.Portfolio{
		PersonID:         personID,
		ShowEmail:        true,
		ShowOrcid:        true,
		PinnedProjectIDs: []uuid.UUID{},
		PinnedProductIDs: []uuid.UUID{},
	}
}
