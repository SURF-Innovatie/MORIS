package kvk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("kvk_not_found")
)

// Client defines the interface for KVK API operations
type Client interface {
	// Search queries the KVK register
	Search(ctx context.Context, query string) (*SearchResponse, error)
	// GetBasicProfile retrieves the basic profile for a given KVK number
	GetBasicProfile(ctx context.Context, kvkNumber string) (*BasicProfile, error)
}

type client struct {
	httpClient *http.Client
	config     *Config
}

// NewClient creates a new KVK API client
func NewClient(config *Config) Client {
	return &client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		config:     config,
	}
}

// NewClientWithHTTP creates a new KVK API client with a custom HTTP client
func NewClientWithHTTP(config *Config, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &client{
		httpClient: httpClient,
		config:     config,
	}
}

func (c *client) Search(ctx context.Context, query string) (*SearchResponse, error) {
	u, err := c.buildURL("/v2/zoeken")
	if err != nil {
		return nil, fmt.Errorf("build url: %w", err)
	}

	q := u.Query()
	q.Set("naam", query) // Use 'naam' as the primary query param; 'handelsnaam' is invalid in v2 test env
	// Note: The KVK API strictly separates parameters. Simple search usually involves 'handelsnaam' or 'kvkNummer'.
	// Assuming 'query' is meant for name search here.
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return c.doSearchRequest(req)
}

func (c *client) GetBasicProfile(ctx context.Context, kvkNumber string) (*BasicProfile, error) {
	path := fmt.Sprintf("/v1/basisprofielen/%s", kvkNumber)
	u, err := c.buildURL(path)
	if err != nil {
		return nil, fmt.Errorf("build url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	var out BasicProfile
	if err := c.doRequest(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) buildURL(path string) (*url.URL, error) {
	baseURL := c.config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.kvk.nl/test/api"
	}

	// Ensure no double slashes if base has trailing or path has leading
	base := strings.TrimRight(baseURL, "/")
	cleanPath := "/" + strings.TrimLeft(path, "/")

	fullURL := base + cleanPath

	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	return u, nil
}

func (c *client) doSearchRequest(req *http.Request) (*SearchResponse, error) {
	var out SearchResponse
	if err := c.doRequest(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) doRequest(req *http.Request, out interface{}) error {
	if c.config.APIKey != "" {
		req.Header.Set("apikey", c.config.APIKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		// Read body to see error detail
		var bodySample string
		if b, err := io.ReadAll(resp.Body); err == nil {
			bodySample = string(b)
		}
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, bodySample)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
