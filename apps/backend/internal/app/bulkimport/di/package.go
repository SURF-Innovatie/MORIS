package di

import (
	"github.com/SURF-Innovatie/MORIS/internal/app/bulkimport"
	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/enttx"
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(provideBulkImportService),
)

func provideBulkImportService(i do.Injector) (bulkimport.Service, error) {
	doiSvc := do.MustInvoke[doi.Service](i)
	productSvc := do.MustInvoke[product.Service](i)
	projectCommandSvc := do.MustInvoke[command.Service](i)
	txManager := do.MustInvoke[*enttx.Manager](i)

	return bulkimport.NewService(doiSvc, productSvc, projectCommandSvc, txManager), nil
}
