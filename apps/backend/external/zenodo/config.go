package zenodo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Config holds the configuration for Zenodo OAuth and API access
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	APIURL       string
}

// GetConfig returns the Zenodo configuration from environment variables
func GetConfig() (*Config, error) {
	clientID := os.Getenv("ZENODO_CLIENT_ID")
	clientSecret := os.Getenv("ZENODO_CLIENT_SECRET")
	redirectURL := os.Getenv("ZENODO_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, fmt.Errorf("ZENODO_CLIENT_ID, ZENODO_CLIENT_SECRET, and ZENODO_REDIRECT_URL must be set")
	}

	// Use sandbox if configured, otherwise production
	isSandbox := os.Getenv("ZENODO_SANDBOX") == "true"
	baseURL := "https://zenodo.org"
	if isSandbox {
		baseURL = "https://sandbox.zenodo.org"
	}

	return &Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		AuthURL:      fmt.Sprintf("%s/oauth/authorize", baseURL),
		TokenURL:     fmt.Sprintf("%s/oauth/token", baseURL),
		APIURL:       fmt.Sprintf("%s/api", baseURL),
	}, nil
}

// IsSandbox returns true if using the sandbox environment
func (c *Config) IsSandbox() bool {
	return strings.Contains(c.AuthURL, "sandbox")
}

// GenerateAuthURL generates the Zenodo OAuth authorization URL
func (c *Config) GenerateAuthURL(state string) string {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("response_type", "code")
	params.Add("scope", "deposit:write deposit:actions")
	params.Add("redirect_uri", c.RedirectURL)
	if state != "" {
		params.Add("state", state)
	}

	return fmt.Sprintf("%s?%s", c.AuthURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for access and refresh tokens
func (c *Config) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Message != "" {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("failed to exchange code, status: %d", resp.StatusCode)
	}

	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.AccessToken == "" {
		return nil, fmt.Errorf("no access token returned")
	}

	return &result, nil
}

// RefreshToken refreshes an expired access token
func (c *Config) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Message != "" {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("failed to refresh token, status: %d", resp.StatusCode)
	}

	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
