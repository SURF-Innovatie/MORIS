package nwo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrNotFound = errors.New("nwo_not_found")
)

// Client defines the interface for NWO Open API operations
type Client interface {
	// GetProjects queries projects with optional filters
	GetProjects(ctx context.Context, opts *QueryOptions) (*ProjectsResponse, error)
	// GetProject retrieves a single project by project_id
	GetProject(ctx context.Context, projectID string) (*Project, error)
}

// QueryOptions contains all available query parameters for the Projects endpoint
type QueryOptions struct {
	ProjectID      string
	GrantID        string
	RORID          string
	Organisation   string
	Title          string
	ReportingYear  int
	RSStartDate    *time.Time // Range start for start_date
	REStartDate    *time.Time // Range end for start_date
	RSEndDate      *time.Time // Range start for end_date
	REEndDate      *time.Time // Range end for end_date
	Summary        string
	MemberLastName string
	Role           ProjectRole
	ORCID          string
	PerPage        int
	Page           int
}

type client struct {
	httpClient *http.Client
	config     *Config
}

// NewClient creates a new NWO API client with default HTTP settings
func NewClient(config *Config) Client {
	return &client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		config:     config,
	}
}

// NewClientWithHTTP creates a new NWO API client with a custom HTTP client
func NewClientWithHTTP(config *Config, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &client{
		httpClient: httpClient,
		config:     config,
	}
}

func (c *client) GetProjects(ctx context.Context, opts *QueryOptions) (*ProjectsResponse, error) {
	u, err := c.buildURL(opts)
	if err != nil {
		return nil, fmt.Errorf("build url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var out ProjectsResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		return &out, nil
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusBadRequest:
		var errResp ExceptionResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("bad request: status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("bad request: %s", errResp.Exception.Message)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (c *client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	opts := &QueryOptions{ProjectID: projectID}
	resp, err := c.GetProjects(ctx, opts)
	if err != nil {
		return nil, err
	}
	if len(resp.Projects) == 0 {
		return nil, ErrNotFound
	}
	return &resp.Projects[0], nil
}

func (c *client) buildURL(opts *QueryOptions) (string, error) {
	baseURL := c.config.BaseURL
	if baseURL == "" {
		baseURL = "https://nwopen-api.nwo.nl"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}
	u.Path = "/NWOpen-API/api/Projects"

	if opts == nil {
		return u.String(), nil
	}

	q := u.Query()

	if opts.ProjectID != "" {
		q.Set("project_id", opts.ProjectID)
	}
	if opts.GrantID != "" {
		q.Set("grant_id", opts.GrantID)
	}
	if opts.RORID != "" {
		q.Set("ror_id", opts.RORID)
	}
	if opts.Organisation != "" {
		q.Set("organisation", opts.Organisation)
	}
	if opts.Title != "" {
		q.Set("title", opts.Title)
	}
	if opts.ReportingYear > 0 {
		q.Set("reporting_year", strconv.Itoa(opts.ReportingYear))
	}
	if opts.RSStartDate != nil {
		q.Set("rs_start_date", opts.RSStartDate.Format("2006-01-02"))
	}
	if opts.REStartDate != nil {
		q.Set("re_start_date", opts.REStartDate.Format("2006-01-02"))
	}
	if opts.RSEndDate != nil {
		q.Set("rs_end_date", opts.RSEndDate.Format("2006-01-02"))
	}
	if opts.REEndDate != nil {
		q.Set("re_end_date", opts.REEndDate.Format("2006-01-02"))
	}
	if opts.Summary != "" {
		q.Set("summary", opts.Summary)
	}
	if opts.MemberLastName != "" {
		q.Set("member_last_name", opts.MemberLastName)
	}
	if opts.Role != "" {
		q.Set("role", string(opts.Role))
	}
	if opts.ORCID != "" {
		q.Set("orcid", opts.ORCID)
	}
	if opts.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if opts.Page > 0 {
		q.Set("page", strconv.Itoa(opts.Page))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
