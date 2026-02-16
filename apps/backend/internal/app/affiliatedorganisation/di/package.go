package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/affiliatedorganisation"
	"github.com/samber/do/v2"
)

var Package = do.Lazy(func(i do.Injector) (int, error) {
	do.Provide(i, func(i do.Injector) (affiliatedorganisation.Repository, error) {
		// This should be provided by the infrastructure layer, but we define the interface requirement here
		// The actual implementation is typically bound in the infra/di package or root di
		return do.Invoke[affiliatedorganisation.Repository](i)
	})

	do.Provide(i, func(i do.Injector) (affiliatedorganisation.Service, error) {
		repo := do.MustInvoke[affiliatedorganisation.Repository](i)
		return affiliatedorganisation.NewService(repo), nil
	})

	return 0, nil
})
