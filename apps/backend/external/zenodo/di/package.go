package di

import (
	"net/http"

	exzenodo "github.com/SURF-Innovatie/MORIS/external/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideZenodoClient),
)

func provideZenodoClient(i do.Injector) (*exzenodo.Client, error) {
	opts := env.ZenodoOptionsFromEnv()
	return exzenodo.NewClient(http.DefaultClient, opts), nil
}
