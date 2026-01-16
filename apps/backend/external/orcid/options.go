package orcid

// Options holds configuration for the ORCID client.
type Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string

	// BaseURL is "https://orcid.org" or "https://sandbox.orcid.org"
	BaseURL string

	// PublicBaseURL is "https://pub.orcid.org/v3.0" or "https://pub.sandbox.orcid.org/v3.0"
	PublicBaseURL string
}

func DefaultOptions(sandbox bool) Options {
	if sandbox {
		return Options{
			BaseURL:       "https://sandbox.orcid.org",
			PublicBaseURL: "https://pub.sandbox.orcid.org/v3.0",
		}
	}
	return Options{
		BaseURL:       "https://orcid.org",
		PublicBaseURL: "https://pub.orcid.org/v3.0",
	}
}
