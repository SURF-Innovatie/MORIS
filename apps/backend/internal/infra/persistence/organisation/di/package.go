package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	organisationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	organisationhierarchyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/hierarchy"
	organisationrbacrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/rbac"
	organisationrolerepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/role"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideOrgRepo),
	do.Lazy(provideOrgRBACRepo),
	do.Lazy(provideOrgHierarchyRepo),
	do.Lazy(provideOrgRoleRepo),
)

func provideOrgRepo(i do.Injector) (*organisationrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return organisationrepo.NewEntRepo(cli), nil
}

func provideOrgRBACRepo(i do.Injector) (*organisationrbacrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return organisationrbacrepo.NewEntRepo(cli), nil
}

func provideOrgHierarchyRepo(i do.Injector) (*organisationhierarchyrepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return organisationhierarchyrepo.NewEntRepo(cli), nil
}

func provideOrgRoleRepo(i do.Injector) (*organisationrolerepo.EntRepo, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return organisationrolerepo.NewEntRepo(cli), nil
}
