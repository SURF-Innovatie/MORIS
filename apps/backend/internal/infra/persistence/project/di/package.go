package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/role"
	projectrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project"
	membershiprepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project/membership"
	projectrolerepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project/role"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProjectRepo),
	do.Lazy(provideMembershipRepo),
	do.Lazy(provideProjectRoleRepo),
)

func provideProjectRepo(i do.Injector) (*projectrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return projectrepo.NewEntRepo(cli), nil
}

func provideMembershipRepo(i do.Injector) (*membershiprepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return membershiprepo.NewEntRepo(cli), nil
}

func provideProjectRoleRepo(i do.Injector) (role.Repository, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return projectrolerepo.NewEntRepo(cli), nil
}
