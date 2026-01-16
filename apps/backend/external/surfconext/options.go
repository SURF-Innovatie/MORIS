package surfconext

// Options configures the SURFconext OpenID Connect (OIDC) client.
//
// Note: SURFconext is the federation; the OIDC "issuer" is environment-specific.
// Use the appropriate IssuerURL for test vs production.
type Options struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

func DefaultOptions() Options {
	return Options{
		Scopes: []string{"openid", "email", "profile"},
	}
}
