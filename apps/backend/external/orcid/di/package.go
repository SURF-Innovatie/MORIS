package di

import (
	"net/http"

	exorcid "github.com/SURF-Innovatie/MORIS/external/orcid"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideORCIDClient),
)

func provideORCIDClient(i do.Injector) (*exorcid.Client, error) {
	opts := env.ORCIDOptionsFromEnv()
	return exorcid.NewClient(http.DefaultClient, opts), nil
}
