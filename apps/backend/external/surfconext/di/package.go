package di

import (
	"net/http"

	exsurfconext "github.com/SURF-Innovatie/MORIS/external/surfconext"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideSurfconextClient),
)

func provideSurfconextClient(i do.Injector) (*exsurfconext.Client, error) {
	opts := env.SurfconextOptionsFromEnv()
	return exsurfconext.NewClient(http.DefaultClient, opts), nil
}
