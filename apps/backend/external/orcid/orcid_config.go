package orcid

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// SetupORCIDOAuth2 initializes the OIDC provider and OAuth2 configuration for ORCID
func SetupORCIDOAuth2(ctx context.Context) (*oauth2.Config, *oidc.Provider, error) {
	clientID := os.Getenv("ORCID_CLIENT_ID")
	clientSecret := os.Getenv("ORCID_CLIENT_SECRET")
	redirectURL := os.Getenv("ORCID_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, nil, fmt.Errorf("ORCID_CLIENT_ID, ORCID_CLIENT_SECRET, and ORCID_REDIRECT_URL must be set")
	}

	// Use sandbox if configured, otherwise production
	isSandbox := os.Getenv("ORCID_SANDBOX") == "true"
	issuerURL := "https://orcid.org"
	if isSandbox {
		issuerURL = "https://sandbox.orcid.org"
	}

	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get ORCID provider: %w", err)
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "/authenticate"}, // ORCID requires /authenticate scope for authentication
	}

	return conf, provider, nil
}
