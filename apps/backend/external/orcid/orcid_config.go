package orcid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ORCIDConfig holds the configuration for ORCID OIDC
type ORCIDConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
}

// GetORCIDConfig returns the ORCID configuration from environment variables
func GetORCIDConfig() (*ORCIDConfig, error) {
	clientID := os.Getenv("ORCID_CLIENT_ID")
	clientSecret := os.Getenv("ORCID_CLIENT_SECRET")
	redirectURL := os.Getenv("ORCID_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, fmt.Errorf("ORCID_CLIENT_ID, ORCID_CLIENT_SECRET, and ORCID_REDIRECT_URL must be set")
	}

	// Use sandbox if configured, otherwise production
	isSandbox := os.Getenv("ORCID_SANDBOX") == "true"
	baseURL := "https://orcid.org"
	if isSandbox {
		baseURL = "https://sandbox.orcid.org"
	}

	return &ORCIDConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		AuthURL:      fmt.Sprintf("%s/oauth/authorize", baseURL),
		TokenURL:     fmt.Sprintf("%s/oauth/token", baseURL),
	}, nil
}

// GenerateAuthURL generates the ORCID authorization URL
func (c *ORCIDConfig) GenerateAuthURL() string {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("response_type", "code")
	params.Add("scope", "/authenticate")
	params.Add("redirect_uri", c.RedirectURL)

	return fmt.Sprintf("%s?%s", c.AuthURL, params.Encode())
}

// ExchangeCode exchanges the authorization code for an access token and ORCID ID
func (c *ORCIDConfig) ExchangeCode(ctx context.Context, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to exchange code, status: %d", resp.StatusCode)
	}

	var result struct {
		ORCID string `json:"orcid"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.ORCID == "" {
		return "", fmt.Errorf("no ORCID ID returned")
	}

	return result.ORCID, nil
}

// GetPublicAPIURL returns the base URL for the public API
func (c *ORCIDConfig) GetPublicAPIURL() string {
	// ORCID Search API requires the /v3.0/ endpoint
	// sandbox: https://pub.sandbox.orcid.org/v3.0/
	// production: https://pub.orcid.org/v3.0/

	if strings.Contains(c.AuthURL, "sandbox") {
		return "https://pub.sandbox.orcid.org/v3.0"
	}
	return "https://pub.orcid.org/v3.0"
}

// GetClientCredentialsToken retrieves an access token using client_credentials flow for public API access
func (c *ORCIDConfig) GetClientCredentialsToken(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "/read-public")

	req, err := http.NewRequestWithContext(ctx, "POST", c.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to obtain token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to obtain token, status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("no access token returned")
	}

	return result.AccessToken, nil
}
