package nwo_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SURF-Innovatie/MORIS/external/nwo"
)

func TestClient_GetProjects_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/NWOpen-API/api/Projects") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"meta": {
				"api_type": "NWO Projects API",
				"version": "1.0.2",
				"count": 1,
				"per_page": 100,
				"pages": 1,
				"page": 1
			},
			"projects": [
				{
					"project_id": "NWO-123",
					"title": "Test Project",
					"funding_scheme": "NWO Open Competition",
					"award_amount": 500000,
					"project_members": [
						{
							"role": "Main Applicant",
							"last_name": "Janssen",
							"first_name": "Jan",
							"orcid": "0000-0001-2345-6789"
						}
					]
				}
			]
		}`))
	}))
	defer ts.Close()

	cfg := &nwo.Config{BaseURL: ts.URL}
	c := nwo.NewClientWithHTTP(cfg, &http.Client{Timeout: 2 * time.Second})

	resp, err := c.GetProjects(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetProjects error: %v", err)
	}
	if resp.Meta.Count != 1 {
		t.Fatalf("expected count 1, got %d", resp.Meta.Count)
	}
	if len(resp.Projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(resp.Projects))
	}
	if resp.Projects[0].ProjectID != "NWO-123" {
		t.Fatalf("expected project_id NWO-123, got %q", resp.Projects[0].ProjectID)
	}
	if resp.Projects[0].Title != "Test Project" {
		t.Fatalf("unexpected title: %q", resp.Projects[0].Title)
	}
	if len(resp.Projects[0].ProjectMembers) != 1 {
		t.Fatalf("expected 1 project member, got %d", len(resp.Projects[0].ProjectMembers))
	}
	if resp.Projects[0].ProjectMembers[0].Role != "Main Applicant" {
		t.Fatalf("unexpected role: %q", resp.Projects[0].ProjectMembers[0].Role)
	}
}

func TestClient_GetProjects_WithQueryOptions(t *testing.T) {
	var capturedQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"meta": {"count": 0}, "projects": []}`))
	}))
	defer ts.Close()

	cfg := &nwo.Config{BaseURL: ts.URL}
	c := nwo.NewClient(cfg)

	opts := &nwo.QueryOptions{
		Title:         "climate",
		Organisation:  "University of Amsterdam",
		ReportingYear: 2024,
		PerPage:       50,
		Page:          2,
	}

	_, err := c.GetProjects(context.Background(), opts)
	if err != nil {
		t.Fatalf("GetProjects error: %v", err)
	}

	// Check that query parameters were set
	if !strings.Contains(capturedQuery, "title=climate") {
		t.Fatalf("missing title parameter in query: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "reporting_year=2024") {
		t.Fatalf("missing reporting_year parameter in query: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "per_page=50") {
		t.Fatalf("missing per_page parameter in query: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "page=2") {
		t.Fatalf("missing page parameter in query: %s", capturedQuery)
	}
}

func TestClient_GetProject_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "project_id=NWO-456") {
			t.Fatalf("expected project_id in query, got: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"meta": {"count": 1},
			"projects": [{"project_id": "NWO-456", "title": "Specific Project"}]
		}`))
	}))
	defer ts.Close()

	cfg := &nwo.Config{BaseURL: ts.URL}
	c := nwo.NewClient(cfg)

	project, err := c.GetProject(context.Background(), "NWO-456")
	if err != nil {
		t.Fatalf("GetProject error: %v", err)
	}
	if project.ProjectID != "NWO-456" {
		t.Fatalf("expected project_id NWO-456, got %q", project.ProjectID)
	}
}

func TestClient_GetProject_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"meta": {"count": 0}, "projects": []}`))
	}))
	defer ts.Close()

	cfg := &nwo.Config{BaseURL: ts.URL}
	c := nwo.NewClient(cfg)

	_, err := c.GetProject(context.Background(), "nonexistent")
	if err == nil || !errors.Is(err, nwo.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestClient_GetProjects_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{
			"exception": {
				"status": 400,
				"error": "Bad Request",
				"message": "Invalid parameter value"
			}
		}`))
	}))
	defer ts.Close()

	cfg := &nwo.Config{BaseURL: ts.URL}
	c := nwo.NewClient(cfg)

	_, err := c.GetProjects(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "Invalid parameter value") {
		t.Fatalf("expected error message to contain 'Invalid parameter value', got: %v", err)
	}
}

func TestClient_GetProjects_404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cfg := &nwo.Config{BaseURL: ts.URL}
	c := nwo.NewClient(cfg)

	_, err := c.GetProjects(context.Background(), nil)
	if err == nil || !errors.Is(err, nwo.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
