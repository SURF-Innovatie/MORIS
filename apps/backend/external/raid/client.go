package raid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Client is the client for the RAiD API.
type Client struct {
	httpClient *http.Client
	options    Options
	
	// Auth state
	token      string
	tokenMutex sync.RWMutex
}

// NewClient creates a new RAiD API client.
func NewClient(httpClient *http.Client, opts Options) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: httpClient,
		options:    opts,
	}
}

// authenticate retrieves a new access token using the configured credentials.
func (c *Client) authenticate(ctx context.Context) error {
	if c.options.Username == "" || c.options.Password == "" {
		return fmt.Errorf("username and password must be set")
	}

	data := url.Values{}
	data.Set("client_id", "raid-api")
	data.Set("username", c.options.Username)
	data.Set("password", c.options.Password)
	data.Set("grant_type", "password")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.options.AuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp RAiDAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.tokenMutex.Lock()
	c.token = authResp.AccessToken
	c.tokenMutex.Unlock()

	return nil
}

// doRequest performs an HTTP request with automatic token injection and retry on 401.
func (c *Client) doRequest(ctx context.Context, method, path string, payload interface{}, retry bool) (*http.Response, error) {
	reqUrl := fmt.Sprintf("%s%s", strings.TrimRight(c.options.BaseURL, "/"), path)

	var body io.Reader
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqUrl, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.tokenMutex.RLock()
	token := c.token
	c.tokenMutex.RUnlock()

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized && !retry {
		resp.Body.Close() // Close the failed response body
		
		// Re-authenticate
		if err := c.authenticate(ctx); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}

		// Retry once
		return c.doRequest(ctx, method, path, payload, true)
	}

	return resp, nil
}

// decodeResponse decodes the JSON response into the target struct.
func decodeResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// MintRaid creates a new RAiD.
func (c *Client) MintRaid(ctx context.Context, req *RAiDCreateRequest) (*RAiDDto, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/raid/", req, false)
	if err != nil {
		return nil, err
	}
	return decodeResponse[RAiDDto](resp)
}

// UpdateRaid updates an existing RAiD.
func (c *Client) UpdateRaid(ctx context.Context, prefix, suffix string, req *RAiDUpdateRequest) (*RAiDDto, error) {
	path := fmt.Sprintf("/raid/%s/%s", prefix, suffix)
	resp, err := c.doRequest(ctx, http.MethodPut, path, req, false)
	if err != nil {
		return nil, err
	}
	return decodeResponse[RAiDDto](resp)
}

// FindRaid retrieves a RAiD by handle.
func (c *Client) FindRaid(ctx context.Context, prefix, suffix string) (*RAiDDto, error) {
	path := fmt.Sprintf("/raid/%s/%s", prefix, suffix)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, false)
	if err != nil {
		return nil, err
	}
	return decodeResponse[RAiDDto](resp)
}

// FindAllRaids retrieves all RAiDs.
func (c *Client) FindAllRaids(ctx context.Context) ([]RAiDDto, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/raid/", nil, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result []RAiDDto
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

