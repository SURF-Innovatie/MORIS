package env

import (
	"strings"

	"github.com/SURF-Innovatie/MORIS/external/surfconext"
)

func SurfconextOptionsFromEnv() surfconext.Options {
	opts := surfconext.DefaultOptions()
	if Global.Surfconext.IssuerURL != "" {
		opts.IssuerURL = Global.Surfconext.IssuerURL
	}
	if Global.Surfconext.ClientID != "" {
		opts.ClientID = Global.Surfconext.ClientID
	}
	if Global.Surfconext.ClientSecret != "" {
		opts.ClientSecret = Global.Surfconext.ClientSecret
	}
	if Global.Surfconext.RedirectURL != "" {
		opts.RedirectURL = Global.Surfconext.RedirectURL
	}

	if scopes := strings.TrimSpace(Global.Surfconext.Scopes); scopes != "" {
		// Allow both comma-separated and space-separated scopes.
		scopes = strings.ReplaceAll(scopes, ",", " ")
		opts.Scopes = strings.Fields(scopes)
	}

	return opts
}
