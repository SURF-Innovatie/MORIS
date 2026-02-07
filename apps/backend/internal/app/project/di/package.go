package di

import (
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/event"
	"github.com/SURF-Innovatie/MORIS/internal/app/eventpolicy"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	organisationhierarchy "github.com/SURF-Innovatie/MORIS/internal/app/organisation/hierarchy"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/cachewarmup"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/load"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	projectrole2 "github.com/SURF-Innovatie/MORIS/internal/app/project/role"
	"github.com/SURF-Innovatie/MORIS/internal/app/user"
	"github.com/SURF-Innovatie/MORIS/internal/infra/cache"
	projectrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/project"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProjectRoleService),
	do.Lazy(provideProjectLoader),
	do.Lazy(provideProjectQueryService),
	do.Lazy(provideProjectCommandService),
	do.Lazy(provideCacheWarmupService),
)

func provideProjectRoleService(i do.Injector) (projectrole2.Service, error) {
	repo := do.MustInvoke[projectrole2.Repository](i)
	orgSvc := do.MustInvoke[organisation.Service](i)
	orgHierarchySvc := do.MustInvoke[organisationhierarchy.Service](i)
	return projectrole2.NewService(repo, orgSvc, orgHierarchySvc), nil
}

func provideProjectLoader(i do.Injector) (*load.Loader, error) {
	eventSvc := do.MustInvoke[event.Service](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	return load.New(eventSvc, pc), nil
}

func provideProjectQueryService(i do.Injector) (queries.Service, error) {
	eventSvc := do.MustInvoke[event.Service](i)
	ldr := do.MustInvoke[*load.Loader](i)
	repo := do.MustInvoke[*projectrepo.EntRepo](i)
	roleRepo := do.MustInvoke[projectrole2.Repository](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	userSvc := do.MustInvoke[user.Service](i)
	return queries.NewService(eventSvc, ldr, repo, roleRepo, curUser, userSvc), nil
}

func provideProjectCommandService(i do.Injector) (command.Service, error) {
	eventSvc := do.MustInvoke[event.Service](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	curUser := do.MustInvoke[coreauth.CurrentUserProvider](i)
	entProv := do.MustInvoke[command.EntClientProvider](i)
	roleSvc := do.MustInvoke[projectrole2.Service](i)
	evaluator := do.MustInvoke[eventpolicy.Evaluator](i)
	orgSvc := do.MustInvoke[organisation.Service](i)
	rbacSvc := do.MustInvoke[organisationrbac.Service](i)
	evtPub := do.MustInvoke[event.Publisher](i)
	repo := do.MustInvoke[*projectrepo.EntRepo](i)
	return command.NewService(eventSvc, pc, curUser, entProv, roleSvc, evaluator, orgSvc, rbacSvc, evtPub, repo), nil
}

func provideCacheWarmupService(i do.Injector) (cachewarmup.Service, error) {
	repo := do.MustInvoke[*projectrepo.EntRepo](i)
	ldr := do.MustInvoke[*load.Loader](i)
	pc := do.MustInvoke[cache.ProjectCache](i)
	return cachewarmup.NewService(repo, ldr, pc), nil
}
