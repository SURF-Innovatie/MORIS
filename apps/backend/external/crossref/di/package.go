package di

import (
	excrossref "github.com/SURF-Innovatie/MORIS/external/crossref"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideCrossrefClient),
)

func provideCrossrefClient(i do.Injector) (excrossref.Client, error) {
	cfg := &excrossref.Config{
		BaseURL:   env.Global.Crossref.BaseURL,
		UserAgent: env.Global.Crossref.UserAgent,
		Mailto:    env.Global.Crossref.Mailto,
	}
	return excrossref.NewClient(cfg), nil
}
