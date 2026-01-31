package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	personsvc "github.com/SURF-Innovatie/MORIS/internal/app/person"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	projectmembershiprepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project/membership"
	userrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/user"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideUserService),
)

func provideUserService(i do.Injector) (user.Service, error) {
	userRepo := do.MustInvoke[*userrepo.EntRepo](i)
	personSvc := do.MustInvoke[personsvc.Service](i)
	eventSvc := do.MustInvoke[event.Service](i)
	membership := do.MustInvoke[*projectmembershiprepo.EntRepo](i)
	userCache := do.MustInvoke[cache.UserCache](i)
	return user.NewService(userRepo, personSvc, eventSvc, membership, userCache), nil
}
