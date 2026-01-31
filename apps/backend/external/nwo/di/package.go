package di

import (
	exnwo "github.com/SURF-Innovatie/MORIS/external/nwo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideNWOClient),
)

func provideNWOClient(i do.Injector) (exnwo.Client, error) {
	cfg := &exnwo.Config{
		BaseURL: env.Global.NWO.BaseURL,
	}
	return exnwo.NewClient(cfg), nil
}
