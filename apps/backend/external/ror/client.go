package ror

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is the client for the ROR API.
type Client struct {
	httpClient *http.Client
	options    ClientOptions
}

// NewClient creates a new ROR API client.
func NewClient(httpClient *http.Client, opts ...ClientOption) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &Client{
		httpClient: httpClient,
		options:    options,
	}
}

// PerformQuery performs a query against the ROR API.
func (c *Client) PerformQuery(ctx context.Context, query string) (*OrganizationsResult, error) {
	reqUrl := fmt.Sprintf("%s?%s", c.options.BaseUrl, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result OrganizationsResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetOrganization retrieves a single organization by ROR ID.
func (c *Client) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	// Handle potential URL encoding if ID contains special chars, though ROR IDs usually don't.
	reqUrl := fmt.Sprintf("%s/%s", c.options.BaseUrl, url.PathEscape(id))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var org Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &org, nil
}

// Query returns a new query builder.
func (c *Client) Query() *OrganizationQueryBuilder {
	return NewOrganizationQueryBuilder(c)
}
