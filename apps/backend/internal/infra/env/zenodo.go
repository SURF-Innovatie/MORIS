package env

import (
	"github.com/SURF-Innovatie/MORIS/external/zenodo"
)

func ZenodoOptionsFromEnv() zenodo.Options {
	opts := zenodo.DefaultOptions(Global.Zenodo.Sandbox)
	if Global.Zenodo.ClientID != "" {
		opts.ClientID = Global.Zenodo.ClientID
	}
	if Global.Zenodo.ClientSecret != "" {
		opts.ClientSecret = Global.Zenodo.ClientSecret
	}
	if Global.Zenodo.RedirectURL != "" {
		opts.RedirectURL = Global.Zenodo.RedirectURL
	}
	return opts
}
