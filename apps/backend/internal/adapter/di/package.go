package di

import (
	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/adapter"
	raidsink "github.com/SURF-Innovatie/MORIS/internal/adapter/sinks/raid"
	csvsource "github.com/SURF-Innovatie/MORIS/internal/adapter/sources/csv"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideAdapterRegistry),
)

func provideAdapterRegistry(i do.Injector) (*adapter.Registry, error) {
	registry := adapter.NewRegistry()
	registry.RegisterSource(csvsource.NewCSVSource("/tmp/import.csv"))

	raidClient := do.MustInvoke[*raid.Client](i)
	registry.RegisterSink(raidsink.NewRAiDSink(raidClient))

	return registry, nil
}
