package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	eventrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/event"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(ProvideEventRepo),
)

func ProvideEventRepo(i do.Injector) (*eventrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return eventrepo.NewEntRepo(cli), nil
}
