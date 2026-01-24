package di

import (
	"net/http"

	excrossref "github.com/SURF-Innovatie/MORIS/external/crossref"
	exkvk "github.com/SURF-Innovatie/MORIS/external/kvk"
	exnwo "github.com/SURF-Innovatie/MORIS/external/nwo"
	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/external/raid"
	exror "github.com/SURF-Innovatie/MORIS/external/ror"
	exsurfconext "github.com/SURF-Innovatie/MORIS/external/surfconext"
	exvies "github.com/SURF-Innovatie/MORIS/external/vies"
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

func provideKVKClient(i do.Injector) (exkvk.Client, error) {
	cfg := &exkvk.Config{
		BaseURL: env.Global.KVK.BaseURL,
		APIKey:  env.Global.KVK.APIKey,
	}
	return exkvk.NewClient(cfg), nil
}

func provideRORClient(i do.Injector) (*exror.Client, error) {
	return exror.NewClient(http.DefaultClient), nil
}

func provideVIESClient(i do.Injector) (*exvies.Client, error) {
	return exvies.NewClient(http.DefaultClient), nil
}
