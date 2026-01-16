package env

import (
	"os"

	"github.com/SURF-Innovatie/MORIS/external/zenodo"
)

func ZenodoOptionsFromEnv() zenodo.Options {
	clientID := os.Getenv("ZENODO_CLIENT_ID")
	clientSecret := os.Getenv("ZENODO_CLIENT_SECRET")
	redirectURL := os.Getenv("ZENODO_REDIRECT_URL")
	sandbox := os.Getenv("ZENODO_SANDBOX") == "true"

	opts := zenodo.DefaultOptions(sandbox)
	opts.ClientID = clientID
	opts.ClientSecret = clientSecret
	opts.RedirectURL = redirectURL
	return opts
}
