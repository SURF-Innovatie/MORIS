package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	eventpolicyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideEventPolicyRepo),
)

func provideEventPolicyRepo(i do.Injector) (*eventpolicyrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return eventpolicyrepo.NewEntRepository(cli), nil
}
