package env

import (
	"os"
	"strings"

	"github.com/SURF-Innovatie/MORIS/external/surfconext"
)

func SurfconextOptionsFromEnv() surfconext.Options {
	opts := surfconext.DefaultOptions()
	opts.IssuerURL = os.Getenv("SURFCONEXT_ISSUER_URL")
	opts.ClientID = os.Getenv("SURFCONEXT_CLIENT_ID")
	opts.ClientSecret = os.Getenv("SURFCONEXT_CLIENT_SECRET")
	opts.RedirectURL = os.Getenv("SURFCONEXT_REDIRECT_URL")

	if scopes := strings.TrimSpace(os.Getenv("SURFCONEXT_SCOPES")); scopes != "" {
		// Allow both comma-separated and space-separated scopes.
		scopes = strings.ReplaceAll(scopes, ",", " ")
		opts.Scopes = strings.Fields(scopes)
	}

	return opts
}
