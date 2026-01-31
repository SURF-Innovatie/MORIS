package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	organisationrole "github.com/SURF-Innovatie/MORIS/internal/app/organisation/role"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/role"
	organisationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/organisation"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideOrganisationHandler),
	do.Lazy(provideOrgRBACHandler),
	do.Lazy(provideOrgRoleHandler),
)

func provideOrgRBACHandler(i do.Injector) (*organisationhandler.RBACHandler, error) {
	svc := do.MustInvoke[organisationrbac.Service](i)
	return organisationhandler.NewRBACHandler(svc), nil
}

func provideOrgRoleHandler(i do.Injector) (*organisationhandler.RoleHandler, error) {
	roleSvc := do.MustInvoke[organisationrole.Service](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	return organisationhandler.NewRoleHandler(roleSvc, rbacSvc), nil
}

func provideOrganisationHandler(i do.Injector) (*organisationhandler.Handler, error) {
	orgSvc := do.MustInvoke[organisation.Service](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	roleSvc := do.MustInvoke[role.Service](i)
	cfSvc := do.MustInvoke[customfield.Service](i)
	return organisationhandler.NewHandler(orgSvc, rbacSvc, roleSvc, cfSvc), nil
}
