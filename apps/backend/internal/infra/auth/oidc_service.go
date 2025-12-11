package auth

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/env"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCProvider interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error)
	Verify(ctx context.Context, tokenString string) (*oidc.IDToken, error)
}

type oidcProvider struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauth2   *oauth2.Config
}

func NewOIDCProvider(ctx context.Context) (OIDCProvider, error) {
	provider, err := oidc.NewProvider(ctx, env.Global.OIDCIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     env.Global.OIDCClientID,
		ClientSecret: env.Global.OIDCClientSecret,
		RedirectURL:  env.Global.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: env.Global.OIDCClientID})

	return &oidcProvider{
		provider: provider,
		verifier: verifier,
		oauth2:   oauth2Config,
	}, nil
}

func (p *oidcProvider) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return p.oauth2.AuthCodeURL(state, opts...)
}

func (p *oidcProvider) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return p.oauth2.Exchange(ctx, code, opts...)
}

func (p *oidcProvider) UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error) {
	return p.provider.UserInfo(ctx, tokenSource)
}

func (p *oidcProvider) Verify(ctx context.Context, tokenString string) (*oidc.IDToken, error) {
	return p.verifier.Verify(ctx, tokenString)
}
