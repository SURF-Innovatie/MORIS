package surfconext

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
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

type Claims struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

func (c *Client) oauthConfig(provider *oidc.Provider) (*oauth2.Config, error) {
	if c.opts.IssuerURL == "" || c.opts.ClientID == "" || c.opts.ClientSecret == "" || c.opts.RedirectURL == "" {
		return nil, fmt.Errorf("surfconext options missing: IssuerURL, ClientID, ClientSecret, RedirectURL")
	}

	scopes := c.opts.Scopes
	if len(scopes) == 0 {
		scopes = DefaultOptions().Scopes
	}

	return &oauth2.Config{
		ClientID:     c.opts.ClientID,
		ClientSecret: c.opts.ClientSecret,
		RedirectURL:  c.opts.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}, nil
}

// AuthURL returns the URL to redirect the user to for SURFconext login.
func (c *Client) AuthURL(ctx context.Context) (string, error) {
	provider, err := oidc.NewProvider(ctx, c.opts.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("discover oidc provider: %w", err)
	}

	cfg, err := c.oauthConfig(provider)
	if err != nil {
		return "", err
	}

	// We don't currently persist/validate state (consistent with existing ORCID/Zenodo flows).
	state := uuid.NewString()
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

// ExchangeCode exchanges an authorization code for user identity claims.
func (c *Client) ExchangeCode(ctx context.Context, code string) (*Claims, error) {
	if code == "" {
		return nil, fmt.Errorf("missing authorization code")
	}

	provider, err := oidc.NewProvider(ctx, c.opts.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("discover oidc provider: %w", err)
	}

	cfg, err := c.oauthConfig(provider)
	if err != nil {
		return nil, err
	}

	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	rawIDToken, ok := tok.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, fmt.Errorf("no id_token in token response")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: c.opts.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("verify id_token: %w", err)
	}

	claims := &Claims{}
	if err := idToken.Claims(claims); err != nil {
		return nil, fmt.Errorf("decode id_token claims: %w", err)
	}

	// SURFconext can provide additional attributes/claims via the userinfo endpoint.
	if claims.Email == "" {
		ui, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(tok))
		if err == nil {
			_ = ui.Claims(claims)
		}
	}

	return claims, nil
}
