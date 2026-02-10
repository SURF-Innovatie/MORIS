package bulkimport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/tx"
	internalevents "github.com/SURF-Innovatie/MORIS/internal/domain/project/events"
	"github.com/google/uuid"
)

type Service interface {
	// BulkImport adds all DOI-derived products to a single existing project.
	BulkImport(ctx context.Context, actorUserID uuid.UUID, actorPersonID uuid.UUID, projectID uuid.UUID, entries []Entry) (*Result, error)
}

type service struct {
	doiSvc            doi.Service
	productSvc        product.Service
	projectCommandSvc command.Service
	tx                tx.Manager
}

func NewService(
	doiSvc doi.Service,
	productSvc product.Service,
	projectCommandSvc command.Service,
	txManager tx.Manager,
) Service {
	return &service{
		doiSvc:            doiSvc,
		productSvc:        productSvc,
		projectCommandSvc: projectCommandSvc,
		tx:                txManager,
	}
}

func (s *service) BulkImport(
	ctx context.Context,
	actorUserID uuid.UUID,
	actorPersonID uuid.UUID,
	projectID uuid.UUID,
	entries []Entry,
) (*Result, error) {
	if projectID == uuid.Nil {
		return nil, fmt.Errorf("projectID is required")
	}

	res := &Result{
		ProjectID: projectID,
		Items:     make([]ItemResult, 0, len(entries)),
	}

	// actorUserID currently not used here; project permission checks happen in ExecuteEvent via current user in ctx.
	_ = actorUserID

	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		createdIDs := make([]uuid.UUID, 0, len(entries))

		for _, e := range entries {
			item := ItemResult{DOI: e.DOI}

			if e.DOI == "" {
				item.DOI = e.DOI
				item.Error = "empty doi"
				res.Errors = append(res.Errors, EntryError{DOI: item.DOI, Error: item.Error})
				res.Items = append(res.Items, item)
				continue
			}

			work, err := s.doiSvc.Resolve(ctx, e.DOI)
			if err != nil {
				item.Error = fmt.Sprintf("resolve doi: %v", err)
				res.Errors = append(res.Errors, EntryError{DOI: e.DOI, Error: item.Error})
				res.Items = append(res.Items, item)
				continue
			}
			item.Work = work

			createdProd, err := s.productSvc.CreateFromWork(ctx, actorPersonID, work)
			if err != nil {
				item.Error = fmt.Sprintf("create product: %v", err)
				res.Errors = append(res.Errors, EntryError{DOI: e.DOI, Error: item.Error})
				res.Items = append(res.Items, item)
				continue
			}

			item.ProductID = createdProd.Id
			res.CreatedProducts = append(res.CreatedProducts, createdProd.Id)
			createdIDs = append(createdIDs, createdProd.Id)

			res.Items = append(res.Items, item)
		}

		if len(createdIDs) > 0 {
			if err := s.addProductsViaBulkEvent(ctx, projectID, createdIDs); err != nil {
				return fmt.Errorf("bulk add products to project: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) addProductsViaBulkEvent(ctx context.Context, projectID uuid.UUID, productIDs []uuid.UUID) error {
	input := internalevents.ProductsBulkImportedInput{ProductIDs: productIDs}
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = s.projectCommandSvc.ExecuteEvent(ctx, command.ExecuteEventRequest{
		ProjectID: projectID,
		Type:      internalevents.ProductsBulkImportedType,
		Input:     b,
	})
	return err
}
