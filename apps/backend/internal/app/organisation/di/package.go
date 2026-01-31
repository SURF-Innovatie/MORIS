package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation/hierarchy"
	organisationhierarchy "github.com/SURF-Innovatie/MORIS/internal/app/organisation/hierarchy"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	organisationrole "github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	organisationrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation"
	organisationhierarchyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/hierarchy"
	organisationrbacrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/rbac"
	organisationrolerepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/organisation/role"
	personrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/person"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideOrganisationService),
	do.Lazy(provideOrgRBACService),
	do.Lazy(provideOrgRoleService),
	do.Lazy(provideOrgHierarchyService),
)

func provideOrganisationService(i do.Injector) (organisation.Service, error) {
	orgRepo := do.MustInvoke[*organisationrepo.EntRepo](i)
	personSvc := do.MustInvoke[*personrepo.EntRepo](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	txManager := do.MustInvoke[*enttx.Manager](i)
	return organisation.NewService(orgRepo, personSvc, rbacSvc, txManager), nil
}

func provideOrgRBACService(i do.Injector) (organisationrbac.Service, error) {
	repo := do.MustInvoke[*organisationrbacrepo.EntRepo](i)
	return organisationrbac.NewService(repo), nil
}

func provideOrgRoleService(i do.Injector) (organisationrole.Service, error) {
	repo := do.MustInvoke[*organisationrolerepo.EntRepo](i)
	return organisationrole.NewService(repo), nil
}

func provideOrgHierarchyService(i do.Injector) (hierarchy.Service, error) {
	repo := do.MustInvoke[*organisationhierarchyrepo.EntRepo](i)
	return organisationhierarchy.NewService(repo), nil
}
