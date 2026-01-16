package orcid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	http *http.Client
	opts Options
}

func NewClient(httpClient *http.Client, opts Options) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{http: httpClient, opts: opts}
}

func (c *Client) AuthURL() (string, error) {
	if c.opts.ClientID == "" || c.opts.RedirectURL == "" || c.opts.BaseURL == "" {
		return "", fmt.Errorf("orcid options missing: ClientID, RedirectURL, BaseURL")
	}

	authURL := strings.TrimRight(c.opts.BaseURL, "/") + "/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", c.opts.ClientID)
	params.Add("response_type", "code")
	params.Add("scope", "/authenticate")
	params.Add("redirect_uri", c.opts.RedirectURL)

	return fmt.Sprintf("%s?%s", authURL, params.Encode()), nil
}

// ExchangeCode exchanges an authorization code for the ORCID iD.
func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	if c.opts.ClientID == "" || c.opts.ClientSecret == "" || c.opts.RedirectURL == "" || c.opts.BaseURL == "" {
		return "", fmt.Errorf("orcid options missing")
	}

	tokenURL := strings.TrimRight(c.opts.BaseURL, "/") + "/oauth/token"

	data := url.Values{}
	data.Set("client_id", c.opts.ClientID)
	data.Set("client_secret", c.opts.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.opts.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed: status %d", resp.StatusCode)
	}

	var result struct {
		ORCID string `json:"orcid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if result.ORCID == "" {
		return "", fmt.Errorf("no orcid returned")
	}
	return result.ORCID, nil
}

func (c *Client) clientCredentialsToken(ctx context.Context) (string, error) {
	if c.opts.ClientID == "" || c.opts.ClientSecret == "" || c.opts.BaseURL == "" {
		return "", fmt.Errorf("orcid options missing")
	}

	tokenURL := strings.TrimRight(c.opts.BaseURL, "/") + "/oauth/token"

	data := url.Values{}
	data.Set("client_id", c.opts.ClientID)
	data.Set("client_secret", c.opts.ClientSecret)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "/read-public")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed: status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("no access token returned")
	}
	return result.AccessToken, nil
}

// SearchExpanded searches the ORCID public registry via expanded-search endpoint.
func (c *Client) SearchExpanded(ctx context.Context, query string) ([]OrcidPerson, error) {
	if c.opts.PublicBaseURL == "" {
		return nil, fmt.Errorf("orcid options missing: PublicBaseURL")
	}

	token, err := c.clientCredentialsToken(ctx)
	if err != nil {
		return nil, err
	}

	searchURL := fmt.Sprintf("%s/expanded-search?q=%s", strings.TrimRight(c.opts.PublicBaseURL, "/"), url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create search request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.orcid+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: status %d", resp.StatusCode)
	}

	var response expandedSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	out := make([]OrcidPerson, 0, len(response.ExpandedResult))
	for _, r := range response.ExpandedResult {
		out = append(out, r.ToPerson())
	}
	return out, nil
}
