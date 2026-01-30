package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	raidsink "github.com/SURF-Innovatie/MORIS/internal/adapter/sinks/raid"
	csvsource "github.com/SURF-Innovatie/MORIS/internal/adapter/sources/csv"
	organisationrbac "github.com/SURF-Innovatie/MORIS/internal/app/organisation/rbac"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	adapterhandler "github.com/SURF-Innovatie/MORIS/internal/handler/adapter"
	"github.com/SURF-Innovatie/MORIS/internal/infra/adapters/event_policy"
	"github.com/samber/do/v2"
)

func provideOrgClosureAdapter(i do.Injector) (*event_policy.OrgClosureAdapter, error) {
	orgRbac := do.MustInvoke[organisationrbac.Service](i)
	return event_policy.NewOrgClosureAdapter(orgRbac), nil
}

func provideRecipientAdapter(i do.Injector) (*event_policy.RecipientAdapter, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return event_policy.NewRecipientAdapter(cli), nil
}

func provideAdapterRegistry(i do.Injector) (*adapter.Registry, error) {
	registry := adapter.NewRegistry()
	registry.RegisterSource(csvsource.NewCSVSource("/tmp/import.csv"))

	raidClient := do.MustInvoke[*raid.Client](i)
	registry.RegisterSink(raidsink.NewRAiDSink(raidClient))

	return registry, nil
}

func provideAdapterHandler(i do.Injector) (*adapterhandler.Handler, error) {
	registry := do.MustInvoke[*adapter.Registry](i)
	projSvc := do.MustInvoke[queries.Service](i)
	return adapterhandler.NewHandler(registry, projSvc), nil
}
