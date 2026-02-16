package affiliatedorganisation

import (
	"encoding/json"
	"net/http"

	"github.com/SURF-Innovatie/MORIS/internal/app/affiliatedorganisation"
	domain "github.com/SURF-Innovatie/MORIS/internal/domain/affiliatedorganisation"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Handler struct {
	service affiliatedorganisation.Service
}

func NewHandler(service affiliatedorganisation.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/vat-lookup", h.LookupVAT)
	r.Get("/", h.GetAll)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// Create godoc
// @Summary      Create affiliated organisation
// @Description  Create a new affiliated organisation
// @Tags         affiliated-organisations
// @Accept       json
// @Produce      json
// @Param        request body domain.AffiliatedOrganisation true "Organisation data"
// @Success      200  {object}  domain.AffiliatedOrganisation
// @Router       /affiliated-organisations [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body domain.AffiliatedOrganisation
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	created, err := h.service.Create(r.Context(), body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, created)
}

// LookupVAT godoc
// @Summary      Lookup VAT number
// @Description  Validate and look up details for a VAT number
// @Tags         affiliated-organisations
// @Param        vat_number query string true "VAT Number (e.g. NL822655287B01)"
// @Success      200  {object}  map[string]interface{}
// @Router       /affiliated-organisations/vat-lookup [get]
func (h *Handler) LookupVAT(w http.ResponseWriter, r *http.Request) {
	vatNumber := r.URL.Query().Get("vat_number")
	if vatNumber == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "vat_number is required"})
		return
	}

	// Mock implementation
	if len(vatNumber) > 5 {
		render.JSON(w, r, map[string]interface{}{
			"valid":        true,
			"vat_number":   vatNumber,
			"name":         "Mock Company",
			"address":      "Mock Address 123",
			"country_code": "NL",
		})
	} else {
		render.JSON(w, r, map[string]interface{}{
			"valid": false,
		})
	}
}

// GetAll godoc
// @Summary      Get all affiliated organisations
// @Description  Retrieve a list of all affiliated organisations
// @Tags         affiliated-organisations
// @Produce      json
// @Success      200  {array}  domain.AffiliatedOrganisation
// @Router       /affiliated-organisations [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.service.GetAll(r.Context())
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, r, orgs)
}

// Get godoc
// @Summary      Get affiliated organisation by ID
// @Description  Retrieve a single affiliated organisation by its ID
// @Tags         affiliated-organisations
// @Produce      json
// @Param        id   path      string  true  "Organisation ID"
// @Success      200  {object}  domain.AffiliatedOrganisation
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /affiliated-organisations/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	org, err := h.service.Get(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, r, org)
}

// Update godoc
// @Summary      Update affiliated organisation
// @Description  Update an existing affiliated organisation by ID
// @Tags         affiliated-organisations
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "Organisation ID"
// @Param        request body      domain.AffiliatedOrganisation true "Organisation data"
// @Success      200     {object}  domain.AffiliatedOrganisation
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /affiliated-organisations/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	var body domain.AffiliatedOrganisation
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	updated, err := h.service.Update(r.Context(), id, body)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, r, updated)
}

// Delete godoc
// @Summary      Delete affiliated organisation
// @Description  Delete an affiliated organisation by ID
// @Tags         affiliated-organisations
// @Param        id   path      string  true  "Organisation ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /affiliated-organisations/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.NoContent(w, r)
}
