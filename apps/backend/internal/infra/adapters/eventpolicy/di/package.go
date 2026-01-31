package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/infra/adapters/eventpolicy"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideRecipientAdapter),
)

func provideRecipientAdapter(i do.Injector) (*eventpolicy.RecipientAdapter, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return eventpolicy.NewRecipientAdapter(cli), nil
}
