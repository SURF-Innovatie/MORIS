package odata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	appOdata "github.com/SURF-Innovatie/MORIS/internal/app/odata"
	"github.com/SURF-Innovatie/MORIS/internal/infra/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// OData constants
const (
	ODataVersion     = "4.0"
	ODataContentType = "application/json;odata.metadata=minimal;odata.streaming=true;charset=utf-8"
	ODataNamespace   = "MORIS.OData"
)

// Handler handles OData HTTP requests for Power BI integration
type Handler struct {
	service *appOdata.Service
}

// NewHandler creates a new OData handler
func NewHandler(service *appOdata.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers OData routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/odata", func(r chi.Router) {
		// OData discovery endpoints (required by Power BI)
		r.Get("/", h.GetServiceDocument)
		r.Get("/$metadata", h.GetMetadata)

		// Entity set endpoints
		r.Get("/Budgets", h.GetBudgets)
		r.Get("/BudgetLineItems", h.GetBudgetLineItems)
		r.Get("/BudgetActuals", h.GetBudgetActuals)
		r.Get("/BudgetAnalytics", h.GetBudgetAnalytics)

		// Keep old routes for backward compatibility
		r.Get("/budgets", h.GetBudgets)
		r.Get("/budget-line-items", h.GetBudgetLineItems)
		r.Get("/budget-actuals", h.GetBudgetActuals)
		r.Get("/budget-analytics", h.GetBudgetAnalytics)
	})
}

// GetServiceDocument returns the OData service document listing all entity sets
// This is required by Power BI to discover available data
// @Summary Get OData service document
// @Description Returns the OData service document listing all available entity sets
// @Tags odata
// @Produce json
// @Success 200 {object} map[string]any "OData service document"
// @Router /odata [get]
func (h *Handler) GetServiceDocument(w http.ResponseWriter, r *http.Request) {
	baseURL := getBaseURL(r)

	serviceDoc := map[string]any{
		"@odata.context": baseURL + "/$metadata",
		"value": []map[string]string{
			{"name": "Budgets", "kind": "EntitySet", "url": "Budgets"},
			{"name": "BudgetLineItems", "kind": "EntitySet", "url": "BudgetLineItems"},
			{"name": "BudgetActuals", "kind": "EntitySet", "url": "BudgetActuals"},
			{"name": "BudgetAnalytics", "kind": "EntitySet", "url": "BudgetAnalytics"},
		},
	}

	w.Header().Set("Content-Type", ODataContentType)
	w.Header().Set("OData-Version", ODataVersion)
	_ = httputil.WriteJSON(w, http.StatusOK, serviceDoc)
}

// GetMetadata returns the OData $metadata document describing entity types
// This is required by Power BI to understand the data model
// @Summary Get OData metadata document
// @Description Returns the CSDL/EDMX metadata document describing entity types
// @Tags odata
// @Produce application/xml
// @Success 200 {string} string "OData CSDL/EDMX metadata"
// @Router /odata/$metadata [get]
func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	metadata := generateMetadataXML()
	w.Header().Set("Content-Type", "application/xml;charset=utf-8")
	w.Header().Set("OData-Version", ODataVersion)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(metadata))
}

// generateMetadataXML generates the CSDL/EDMX metadata document
func generateMetadataXML() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<edmx:Edmx Version="4.0" xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx">
  <edmx:DataServices>
    <Schema Namespace="` + ODataNamespace + `" xmlns="http://docs.oasis-open.org/odata/ns/edm">
      <EntityType Name="Budget">
        <Key><PropertyRef Name="id"/></Key>
        <Property Name="id" Type="Edm.Guid" Nullable="false"/>
        <Property Name="projectId" Type="Edm.Guid" Nullable="false"/>
        <Property Name="title" Type="Edm.String"/>
        <Property Name="status" Type="Edm.String"/>
        <Property Name="totalAmount" Type="Edm.Double"/>
        <Property Name="totalBudgeted" Type="Edm.Double"/>
        <Property Name="totalActuals" Type="Edm.Double"/>
        <Property Name="burnRate" Type="Edm.Double"/>
        <Property Name="currency" Type="Edm.String"/>
      </EntityType>
      <EntityType Name="BudgetLineItem">
        <Key><PropertyRef Name="id"/></Key>
        <Property Name="id" Type="Edm.Guid" Nullable="false"/>
        <Property Name="budgetId" Type="Edm.Guid" Nullable="false"/>
        <Property Name="category" Type="Edm.String"/>
        <Property Name="description" Type="Edm.String"/>
        <Property Name="budgetedAmount" Type="Edm.Double"/>
        <Property Name="year" Type="Edm.Int32"/>
        <Property Name="fundingSource" Type="Edm.String"/>
        <Property Name="totalActuals" Type="Edm.Double"/>
      </EntityType>
      <EntityType Name="BudgetActual">
        <Key><PropertyRef Name="id"/></Key>
        <Property Name="id" Type="Edm.Guid" Nullable="false"/>
        <Property Name="lineItemId" Type="Edm.Guid" Nullable="false"/>
        <Property Name="amount" Type="Edm.Double"/>
        <Property Name="description" Type="Edm.String"/>
        <Property Name="recordedDate" Type="Edm.String"/>
        <Property Name="source" Type="Edm.String"/>
      </EntityType>
      <EntityType Name="BudgetAnalytics">
        <Key><PropertyRef Name="budgetId"/></Key>
        <Property Name="projectId" Type="Edm.Guid" Nullable="false"/>
        <Property Name="projectTitle" Type="Edm.String"/>
        <Property Name="budgetId" Type="Edm.Guid" Nullable="false"/>
        <Property Name="totalBudgeted" Type="Edm.Double"/>
        <Property Name="totalActuals" Type="Edm.Double"/>
        <Property Name="remaining" Type="Edm.Double"/>
        <Property Name="burnRate" Type="Edm.Double"/>
        <Property Name="status" Type="Edm.String"/>
      </EntityType>
      <EntityContainer Name="Container">
        <EntitySet Name="Budgets" EntityType="` + ODataNamespace + `.Budget"/>
        <EntitySet Name="BudgetLineItems" EntityType="` + ODataNamespace + `.BudgetLineItem"/>
        <EntitySet Name="BudgetActuals" EntityType="` + ODataNamespace + `.BudgetActual"/>
        <EntitySet Name="BudgetAnalytics" EntityType="` + ODataNamespace + `.BudgetAnalytics"/>
      </EntityContainer>
    </Schema>
  </edmx:DataServices>
</edmx:Edmx>`
}

// getBaseURL extracts the base URL for the OData service
func getBaseURL(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		if fwdProto := r.Header.Get("X-Forwarded-Proto"); fwdProto != "" {
			scheme = fwdProto
		} else {
			scheme = "http"
		}
	}

	host := r.Host
	if fwdHost := r.Header.Get("X-Forwarded-Host"); fwdHost != "" {
		host = fwdHost
	}

	// Get the path up to /odata
	path := r.URL.Path
	if idx := strings.Index(path, "/odata"); idx >= 0 {
		path = path[:idx+6] // include "/odata"
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

// GetBudgets godoc
// @Summary Get budgets via OData
// @Description OData endpoint for Power BI - returns budgets with $select, $filter, $orderby, $top, $skip support
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param $select query string false "Fields to select (comma-separated)"
// @Param $filter query string false "OData filter expression"
// @Param $orderby query string false "Order by field (e.g., 'title desc')"
// @Param $top query integer false "Maximum number of results"
// @Param $skip query integer false "Number of results to skip"
// @Param $count query boolean false "Include total count in response"
// @Success 200 {object} map[string]any "OData response with @odata.count and value array"
// @Failure 400 {object} map[string]any "OData error response"
// @Failure 401 {string} string "unauthorized"
// @Router /odata/budgets [get]
func (h *Handler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetBudgets(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponseWithContext(w, r, "Budgets", result)
}

// GetBudgetLineItems godoc
// @Summary Get budget line items via OData
// @Description OData endpoint for Power BI - returns budget line items
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param $select query string false "Fields to select"
// @Param $filter query string false "OData filter expression"
// @Param $orderby query string false "Order by field"
// @Param $top query integer false "Maximum number of results"
// @Param $skip query integer false "Number of results to skip"
// @Success 200 {object} map[string]any "OData response"
// @Failure 400 {object} map[string]any "OData error response"
// @Router /odata/budget-line-items [get]
func (h *Handler) GetBudgetLineItems(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetLineItems(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponseWithContext(w, r, "BudgetLineItems", result)
}

// GetBudgetActuals godoc
// @Summary Get budget actuals via OData
// @Description OData endpoint for Power BI - returns recorded actuals
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param $select query string false "Fields to select"
// @Param $filter query string false "OData filter expression"
// @Param $top query integer false "Maximum number of results"
// @Param $skip query integer false "Number of results to skip"
// @Success 200 {object} map[string]any "OData response"
// @Failure 400 {object} map[string]any "OData error response"
// @Router /odata/budget-actuals [get]
func (h *Handler) GetBudgetActuals(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetActuals(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponseWithContext(w, r, "BudgetActuals", result)
}

// GetBudgetAnalytics godoc
// @Summary Get budget analytics via OData
// @Description OData endpoint for Power BI - returns pre-computed analytics summary per project
// @Tags odata
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any "OData response with analytics data"
// @Failure 401 {string} string "unauthorized"
// @Router /odata/budget-analytics [get]
func (h *Handler) GetBudgetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == uuid.Nil {
		writeODataError(w, http.StatusUnauthorized, "unauthorized", "Valid authentication required")
		return
	}

	result, err := h.service.GetAnalytics(r.Context(), userID, r.URL.Query())
	if err != nil {
		writeODataError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	writeODataResponseWithContext(w, r, "BudgetAnalytics", result)
}

// Helper functions

func getUserIDFromContext(r *http.Request) uuid.UUID {
	userIDPtr := httputil.GetUserIDFromContext(r.Context())
	if userIDPtr == nil {
		return uuid.Nil
	}
	return *userIDPtr
}

// writeODataResponseWithContext wraps the result with @odata.context
func writeODataResponseWithContext(w http.ResponseWriter, r *http.Request, entitySet string, result any) {
	w.Header().Set("Content-Type", ODataContentType)
	w.Header().Set("OData-Version", ODataVersion)

	baseURL := getBaseURL(r)
	contextURL := fmt.Sprintf("%s/$metadata#%s", baseURL, entitySet)

	// Convert result to map to flatten the structure
	// This avoids nesting the ODataResult struct inside a "value" field
	var response map[string]any

	data, err := json.Marshal(result)
	if err != nil {
		writeODataError(w, http.StatusInternalServerError, "serialization_error", "Failed to serialize response")
		return
	}

	if err := json.Unmarshal(data, &response); err != nil {
		writeODataError(w, http.StatusInternalServerError, "serialization_error", "Failed to prepare response")
		return
	}

	// Add context annotation
	response["@odata.context"] = contextURL

	_ = httputil.WriteJSON(w, http.StatusOK, response)
}

func writeODataError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", ODataContentType)
	w.Header().Set("OData-Version", ODataVersion)
	_ = httputil.WriteJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
