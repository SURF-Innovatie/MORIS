package bulkimport

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/app/doi"
	"github.com/SURF-Innovatie/MORIS/internal/app/product"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/command"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
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
	projectQuerySvc   queries.Service
	tx                tx.Manager
}

func NewService(
	doiSvc doi.Service,
	productSvc product.Service,
	projectCommandSvc command.Service,
	projectQuerySvc queries.Service,
	txManager tx.Manager,
) Service {
	return &service{
		doiSvc:            doiSvc,
		productSvc:        productSvc,
		projectCommandSvc: projectCommandSvc,
		projectQuerySvc:   projectQuerySvc,
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

	_ = actorUserID

	err := s.tx.WithTx(ctx, func(ctx context.Context) error {
		// Load project once so we can skip adding already-linked products
		proj, err := s.projectQuerySvc.GetProject(ctx, projectID)
		if err != nil {
			return fmt.Errorf("load project: %w", err)
		}
		if proj == nil {
			return fmt.Errorf("project not found: %s", projectID)
		}

		alreadyInProject := make(map[uuid.UUID]struct{}, len(proj.Products))
		for _, proj := range proj.Products {
			alreadyInProject[proj.Id] = struct{}{}
		}

		// Only add those not already present
		toAdd := make([]uuid.UUID, 0, len(entries))

		for _, e := range entries {
			item := ItemResult{DOI: e.DOI}

			if strings.TrimSpace(e.DOI) == "" {
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

			// reuse existing product if DOI exists in DB
			p, createdNew, err := s.productSvc.CreateOrGetFromWork(ctx, actorPersonID, work)
			if err != nil {
				item.Error = fmt.Sprintf("create/get product: %v", err)
				res.Errors = append(res.Errors, EntryError{DOI: e.DOI, Error: item.Error})
				res.Items = append(res.Items, item)
				continue
			}

			item.ProductID = p.Id
			if createdNew {
				res.CreatedProducts = append(res.CreatedProducts, p.Id)
			}

			// Skip adding if project already has this product id
			if _, ok := alreadyInProject[p.Id]; ok {
				res.Items = append(res.Items, item)
				continue
			}

			// Mark as to be added, and update set to avoid duplicates within the same request
			toAdd = append(toAdd, p.Id)
			alreadyInProject[p.Id] = struct{}{}

			res.Items = append(res.Items, item)
		}

		if len(toAdd) > 0 {
			if err := s.addProductsViaBulkEvent(ctx, projectID, toAdd); err != nil {
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
