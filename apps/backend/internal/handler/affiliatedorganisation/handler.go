package affiliatedorganisation

import (
	"net/http"
	"strings"

	"github.com/SURF-Innovatie/MORIS/external/vies"
	"github.com/SURF-Innovatie/MORIS/internal/api/dto"
	affiliatedorganisationsvc "github.com/SURF-Innovatie/MORIS/internal/app/affiliatedorganisation"
	"github.com/SURF-Innovatie/MORIS/internal/common/transform"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
)

// Handler handles HTTP requests for affiliated organisations.
type Handler struct {
	svc        affiliatedorganisationsvc.Service
	viesClient *vies.Client
}

// NewHandler creates a new affiliated organisation handler.
func NewHandler(svc affiliatedorganisationsvc.Service, viesClient *vies.Client) *Handler {
	return &Handler{
		svc:        svc,
		viesClient: viesClient,
	}
}

// Create godoc
// @Summary Create an affiliated organisation
// @Description Creates a new affiliated organisation
// @Tags affiliated-organisations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organisation body dto.AffiliatedOrganisationRequest true "Organisation data"
// @Success 200 {object} dto.AffiliatedOrganisationResponse
// @Failure 400 {string} string "invalid body"
// @Failure 500 {string} string "internal server error"
// @Router /affiliated-organisations [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.AffiliatedOrganisationRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	org, err := h.svc.Create(r.Context(), entities.AffiliatedOrganisation{
		Name:      req.Name,
		KvkNumber: req.KvkNumber,
		RorID:     req.RorID,
		VatNumber: req.VatNumber,
		City:      req.City,
		Country:   req.Country,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.AffiliatedOrganisationResponse](*org))
}

// GetAll godoc
// @Summary List affiliated organisations
// @Description Returns all affiliated organisations
// @Tags affiliated-organisations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.AffiliatedOrganisationResponse
// @Failure 500 {string} string "internal server error"
// @Router /affiliated-organisations [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.svc.GetAll(r.Context())
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	dtos := make([]dto.AffiliatedOrganisationResponse, len(orgs))
	for i, org := range orgs {
		dtos[i] = transform.ToDTOItem[dto.AffiliatedOrganisationResponse](*org)
	}

	_ = httputil.WriteJSON(w, http.StatusOK, dtos)
}

// Get godoc
// @Summary Get an affiliated organisation
// @Description Get a single affiliated organisation by ID
// @Tags affiliated-organisations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID (UUID)"
// @Success 200 {object} dto.AffiliatedOrganisationResponse
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /affiliated-organisations/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	org, err := h.svc.Get(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.AffiliatedOrganisationResponse](*org))
}

// Update godoc
// @Summary Update an affiliated organisation
// @Description Update an existing affiliated organisation
// @Tags affiliated-organisations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID (UUID)"
// @Param organisation body dto.AffiliatedOrganisationRequest true "Organisation data"
// @Success 200 {object} dto.AffiliatedOrganisationResponse
// @Failure 400 {string} string "invalid id or body"
// @Failure 500 {string} string "internal server error"
// @Router /affiliated-organisations/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req dto.AffiliatedOrganisationRequest
	if !httputil.ReadJSON(w, r, &req) {
		return
	}

	if req.Name == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "name is required", nil)
		return
	}

	org, err := h.svc.Update(r.Context(), id, entities.AffiliatedOrganisation{
		Name:      req.Name,
		KvkNumber: req.KvkNumber,
		RorID:     req.RorID,
		VatNumber: req.VatNumber,
		City:      req.City,
		Country:   req.Country,
	})
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	_ = httputil.WriteJSON(w, http.StatusOK, transform.ToDTOItem[dto.AffiliatedOrganisationResponse](*org))
}

// Delete godoc
// @Summary Delete an affiliated organisation
// @Description Delete an affiliated organisation by ID
// @Tags affiliated-organisations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Organisation ID (UUID)"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal server error"
// @Router /affiliated-organisations/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.ParseUUIDParam(r, "id")
	if err != nil {
		httputil.WriteError(w, r, http.StatusBadRequest, "invalid id", nil)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to delete organisation", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// VatLookupResponse is the response for VAT lookup.
type VatLookupResponse struct {
	Valid       bool   `json:"valid"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	CountryCode string `json:"country_code"`
	VatNumber   string `json:"vat_number"`
	City        string `json:"city"`
}

// LookupVAT godoc
// @Summary Lookup VAT number via VIES
// @Description Validates a VAT number using the EU VIES API and returns company info
// @Tags affiliated-organisations
// @Produce json
// @Security BearerAuth
// @Param vat_number query string true "VAT number with country code prefix (e.g., NL822655287B01)"
// @Success 200 {object} VatLookupResponse
// @Failure 400 {string} string "invalid vat number"
// @Failure 500 {string} string "internal server error"
// @Router /affiliated-organisations/vat/lookup [get]
func (h *Handler) LookupVAT(w http.ResponseWriter, r *http.Request) {
	vatNumber := r.URL.Query().Get("vat_number")
	if vatNumber == "" {
		httputil.WriteError(w, r, http.StatusBadRequest, "vat_number query parameter is required", nil)
		return
	}

	// Remove spaces and normalize
	vatNumber = strings.ReplaceAll(vatNumber, " ", "")
	vatNumber = strings.ToUpper(vatNumber)

	result, err := h.viesClient.CheckVatNumber(r.Context(), vatNumber)
	if err != nil {
		httputil.WriteError(w, r, http.StatusInternalServerError, "failed to lookup VAT number", err.Error())
		return
	}

	city, _ := result.ParsedAddress()

	resp := VatLookupResponse{
		Valid:       result.Valid,
		Name:        strings.TrimSpace(result.Name),
		Address:     strings.TrimSpace(result.Address),
		CountryCode: result.CountryCode,
		VatNumber:   result.CountryCode + result.VatNumber,
		City:        city,
	}

	_ = httputil.WriteJSON(w, http.StatusOK, resp)
}
