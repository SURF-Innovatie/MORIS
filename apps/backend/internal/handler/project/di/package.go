package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/customfield"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	commandhandler "github.com/SURF-Innovatie/MORIS/internal/handler/project/command"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideProjectHandler),
	do.Lazy(provideProjectCommandHandler),
)

func provideProjectHandler(i do.Injector) (*projecthandler.Handler, error) {
	svc := do.MustInvoke[queries.Service](i)
	cfSvc := do.MustInvoke[customfield.Service](i)
	return projecthandler.NewHandler(svc, cfSvc), nil
}

func provideProjectCommandHandler(i do.Injector) (*commandhandler.Handler, error) {
	svc := do.MustInvoke[command.Service](i)
	return commandhandler.NewHandler(svc), nil
}
