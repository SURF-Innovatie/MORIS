package di

import (
	"net/http"

	excrossref "github.com/SURF-Innovatie/MORIS/external/crossref"
	exnwo "github.com/SURF-Innovatie/MORIS/external/nwo"
	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/external/raid"
	exsurfconext "github.com/SURF-Innovatie/MORIS/external/surfconext"
	exzenodo "github.com/SURF-Innovatie/MORIS/external/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

func provideORCIDClient(i do.Injector) (*exorcid.Client, error) {
	opts := env.ORCIDOptionsFromEnv()
	return exorcid.NewClient(http.DefaultClient, opts), nil
}

func provideSurfconextClient(i do.Injector) (*exsurfconext.Client, error) {
	opts := env.SurfconextOptionsFromEnv()
	return exsurfconext.NewClient(http.DefaultClient, opts), nil
}

func provideZenodoClient(i do.Injector) (*exzenodo.Client, error) {
	opts := env.ZenodoOptionsFromEnv()
	return exzenodo.NewClient(http.DefaultClient, opts), nil
}

func provideCrossrefClient(i do.Injector) (excrossref.Client, error) {
	cfg := &excrossref.Config{
		BaseURL:   env.Global.Crossref.BaseURL,
		UserAgent: env.Global.Crossref.UserAgent,
		Mailto:    env.Global.Crossref.Mailto,
	}
	return excrossref.NewClient(cfg), nil
}

func provideRAiDClient(i do.Injector) (*raid.Client, error) {
	opts := raid.DefaultOptions()
	return raid.NewClient(http.DefaultClient, opts), nil
}

func provideNWOClient(i do.Injector) (exnwo.Client, error) {
	cfg := &exnwo.Config{
		BaseURL: env.Global.NWO.BaseURL,
	}
	return exnwo.NewClient(cfg), nil
}
