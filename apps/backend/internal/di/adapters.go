package di

import (
	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	raidsink "github.com/SURF-Innovatie/MORIS/internal/adapter/sinks/raid"
	csvsource "github.com/SURF-Innovatie/MORIS/internal/adapter/sources/csv"
	"github.com/SURF-Innovatie/MORIS/internal/app/organisation"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	adapterhandler "github.com/SURF-Innovatie/MORIS/internal/handler/adapter"
	eventpolicyrepo "github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventpolicy"
	"github.com/samber/do/v2"
)

func provideOrgClosureAdapter(i do.Injector) (*eventpolicyrepo.OrgClosureAdapter, error) {
	orgSvc := do.MustInvoke[organisation.Service](i)
	return eventpolicyrepo.NewOrgClosureAdapter(orgSvc), nil
}

func provideRecipientAdapter(i do.Injector) (*eventpolicyrepo.RecipientAdapter, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return eventpolicyrepo.NewRecipientAdapter(cli), nil
}

func provideNotificationAdapter(i do.Injector) (*eventpolicyrepo.NotificationAdapter, error) {
	cli := do.MustInvoke[*ent.Client](i)
	return eventpolicyrepo.NewNotificationAdapter(cli), nil
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
