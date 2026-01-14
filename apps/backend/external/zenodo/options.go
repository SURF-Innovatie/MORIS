package zenodo

import (
	"fmt"
	"net/url"
	"strings"
)

type Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string

	AuthURL  string
	TokenURL string
	APIURL   string
}

func DefaultOptions(sandbox bool) Options {
	base := "https://zenodo.org"
	if sandbox {
		base = "https://sandbox.zenodo.org"
	}
	return Options{
		AuthURL:  base + "/oauth/authorize",
		TokenURL: base + "/oauth/token",
		APIURL:   base + "/api",
	}
}

func (o Options) ValidateOAuth() error {
	if o.ClientID == "" || o.ClientSecret == "" || o.RedirectURL == "" {
		return fmt.Errorf("zenodo options missing: ClientID, ClientSecret, RedirectURL")
	}
	if o.AuthURL == "" || o.TokenURL == "" || o.APIURL == "" {
		return fmt.Errorf("zenodo options missing: AuthURL, TokenURL, APIURL")
	}
	return nil
}

func (o Options) authURL(state string) (string, error) {
	if o.ClientID == "" || o.RedirectURL == "" || o.AuthURL == "" {
		return "", fmt.Errorf("zenodo options missing: ClientID, RedirectURL, AuthURL")
	}

	params := url.Values{}
	params.Add("client_id", o.ClientID)
	params.Add("response_type", "code")
	params.Add("scope", "deposit:write deposit:actions")
	params.Add("redirect_uri", o.RedirectURL)
	if state != "" {
		params.Add("state", state)
	}

	return strings.TrimRight(o.AuthURL, "/") + "?" + params.Encode(), nil
}
