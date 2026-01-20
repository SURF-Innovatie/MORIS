package env

import (
	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
)

func ORCIDOptionsFromEnv() exorcid.Options {
	opts := exorcid.DefaultOptions(Global.ORCID.Sandbox)
	if Global.ORCID.ClientID != "" {
		opts.ClientID = Global.ORCID.ClientID
	}
	if Global.ORCID.ClientSecret != "" {
		opts.ClientSecret = Global.ORCID.ClientSecret
	}
	if Global.ORCID.RedirectURL != "" {
		opts.RedirectURL = Global.ORCID.RedirectURL
	}
	return opts
}
