package di

import (
	"net/http"

	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideRAiDClient),
)

func provideRAiDClient(i do.Injector) (*raid.Client, error) {
	opts := raid.DefaultOptions()
	return raid.NewClient(http.DefaultClient, opts), nil
}
