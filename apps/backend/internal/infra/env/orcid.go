package env

import (
	"os"

	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
)

func ORCIDOptionsFromEnv() exorcid.Options {
	clientID := os.Getenv("ORCID_CLIENT_ID")
	clientSecret := os.Getenv("ORCID_CLIENT_SECRET")
	redirectURL := os.Getenv("ORCID_REDIRECT_URL")

	sandbox := os.Getenv("ORCID_SANDBOX") == "true"
	opts := exorcid.DefaultOptions(sandbox)
	opts.ClientID = clientID
	opts.ClientSecret = clientSecret
	opts.RedirectURL = redirectURL
	return opts
}
