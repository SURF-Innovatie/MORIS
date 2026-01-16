package zenodo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client provides low-level HTTP access to the Zenodo API
type Client struct {
	httpClient *http.Client
	opts       Options
}

// NewClient creates a new Zenodo API client
func NewClient(httpClient *http.Client, opts Options) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		httpClient: httpClient,
		opts:       opts,
	}
}

func (c *Client) AuthURL(state string) (string, error) {
	return c.opts.authURL(state)
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	if err := c.opts.ValidateOAuth(); err != nil {
		return nil, err
	}
	if code == "" {
		return nil, fmt.Errorf("missing code")
	}

	data := url.Values{}
	data.Set("client_id", c.opts.ClientID)
	data.Set("client_secret", c.opts.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.opts.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.opts.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = resp.StatusCode
			return nil, &apiErr
		}
		return nil, fmt.Errorf("token request failed: status %d", resp.StatusCode)
	}

	var out TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if out.AccessToken == "" {
		return nil, fmt.Errorf("no access token returned")
	}
	return &out, nil
}

// RefreshToken refreshes an expired access token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if err := c.opts.ValidateOAuth(); err != nil {
		return nil, err
	}
	if refreshToken == "" {
		return nil, fmt.Errorf("missing refresh token")
	}

	data := url.Values{}
	data.Set("client_id", c.opts.ClientID)
	data.Set("client_secret", c.opts.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.opts.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = resp.StatusCode
			return nil, &apiErr
		}
		return nil, fmt.Errorf("refresh request failed: status %d", resp.StatusCode)
	}

	var out TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode refresh response: %w", err)
	}
	if out.AccessToken == "" {
		return nil, fmt.Errorf("no access token returned")
	}
	return &out, nil
}

// doRequest performs an authenticated request to the Zenodo API
func (c *Client) doRequest(ctx context.Context, method, path, accessToken string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	u := strings.TrimRight(c.opts.APIURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// parseResponse parses the API response into the provided target
func parseResponse[T any](resp *http.Response, target *T) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Message != "" {
			apiErr.Status = resp.StatusCode
			return &apiErr
		}
		return fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}
	return nil
}

// CreateDeposition creates a new empty deposition
func (c *Client) CreateDeposition(ctx context.Context, accessToken string) (*Deposition, error) {
	resp, err := c.doRequest(ctx, "POST", "/deposit/depositions", accessToken, struct{}{})
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("create deposition failed: %w", err)
	}
	return &deposition, nil
}

// GetDeposition retrieves a deposition by ID
func (c *Client) GetDeposition(ctx context.Context, accessToken string, depositionID int) (*Deposition, error) {
	path := fmt.Sprintf("/deposit/depositions/%d", depositionID)
	resp, err := c.doRequest(ctx, "GET", path, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("get deposition failed: %w", err)
	}
	return &deposition, nil
}

// UpdateDeposition updates the metadata of a deposition
func (c *Client) UpdateDeposition(ctx context.Context, accessToken string, depositionID int, metadata *DepositionMetadata) (*Deposition, error) {
	path := fmt.Sprintf("/deposit/depositions/%d", depositionID)
	body := map[string]interface{}{
		"metadata": metadata,
	}
	resp, err := c.doRequest(ctx, "PUT", path, accessToken, body)
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("update deposition failed: %w", err)
	}
	return &deposition, nil
}

// DeleteDeposition deletes a deposition
func (c *Client) DeleteDeposition(ctx context.Context, accessToken string, depositionID int) error {
	path := fmt.Sprintf("/deposit/depositions/%d", depositionID)
	resp, err := c.doRequest(ctx, "DELETE", path, accessToken, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Message != "" {
			return &apiErr
		}
		return fmt.Errorf("delete deposition failed: status %d", resp.StatusCode)
	}
	return nil
}

// ListDepositions lists all depositions for the authenticated user
func (c *Client) ListDepositions(ctx context.Context, accessToken string) ([]Deposition, error) {
	resp, err := c.doRequest(ctx, "GET", "/deposit/depositions", accessToken, nil)
	if err != nil {
		return nil, err
	}

	var depositions []Deposition
	if err := parseResponse(resp, &depositions); err != nil {
		return nil, fmt.Errorf("list depositions failed: %w", err)
	}
	return depositions, nil
}

// UploadFile uploads a file to a deposition's bucket
// The new Zenodo Files API supports files up to 50GB
func (c *Client) UploadFile(ctx context.Context, accessToken string, bucketURL, filename string, data io.Reader) (*DepositionFile, error) {
	url := fmt.Sprintf("%s/%s", bucketURL, filename)

	// Read all data into a buffer to get content length
	// This is required because Zenodo's API needs Content-Length header
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = size

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request failed: %w", err)
	}

	var file DepositionFile
	if err := parseResponse(resp, &file); err != nil {
		return nil, fmt.Errorf("upload file failed: %w", err)
	}
	return &file, nil
}

// ListFiles lists all files in a deposition
func (c *Client) ListFiles(ctx context.Context, accessToken string, depositionID int) ([]DepositionFile, error) {
	path := fmt.Sprintf("/deposit/depositions/%d/files", depositionID)
	resp, err := c.doRequest(ctx, "GET", path, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var files []DepositionFile
	if err := parseResponse(resp, &files); err != nil {
		return nil, fmt.Errorf("list files failed: %w", err)
	}
	return files, nil
}

// DeleteFile deletes a file from a deposition
func (c *Client) DeleteFile(ctx context.Context, accessToken string, depositionID int, fileID string) error {
	path := fmt.Sprintf("/deposit/depositions/%d/files/%s", depositionID, fileID)
	resp, err := c.doRequest(ctx, "DELETE", path, accessToken, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("delete file failed: status %d", resp.StatusCode)
	}
	return nil
}

// Publish publishes a deposition, making it publicly available and minting a DOI
func (c *Client) Publish(ctx context.Context, accessToken string, depositionID int) (*Deposition, error) {
	path := fmt.Sprintf("/deposit/depositions/%d/actions/publish", depositionID)
	resp, err := c.doRequest(ctx, "POST", path, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("publish failed: %w", err)
	}
	return &deposition, nil
}

// Edit unlocks a published deposition for editing (creates a new version)
func (c *Client) Edit(ctx context.Context, accessToken string, depositionID int) (*Deposition, error) {
	path := fmt.Sprintf("/deposit/depositions/%d/actions/edit", depositionID)
	resp, err := c.doRequest(ctx, "POST", path, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("edit failed: %w", err)
	}
	return &deposition, nil
}

// Discard discards changes in the current editing session
func (c *Client) Discard(ctx context.Context, accessToken string, depositionID int) (*Deposition, error) {
	path := fmt.Sprintf("/deposit/depositions/%d/actions/discard", depositionID)
	resp, err := c.doRequest(ctx, "POST", path, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("discard failed: %w", err)
	}
	return &deposition, nil
}

// NewVersion creates a new version of a published deposition
func (c *Client) NewVersion(ctx context.Context, accessToken string, depositionID int) (*Deposition, error) {
	path := fmt.Sprintf("/deposit/depositions/%d/actions/newversion", depositionID)
	resp, err := c.doRequest(ctx, "POST", path, accessToken, nil)
	if err != nil {
		return nil, err
	}

	var deposition Deposition
	if err := parseResponse(resp, &deposition); err != nil {
		return nil, fmt.Errorf("new version failed: %w", err)
	}
	return &deposition, nil
}
