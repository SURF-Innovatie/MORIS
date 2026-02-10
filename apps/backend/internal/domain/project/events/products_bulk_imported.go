package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const ProductsBulkImportedType = "project.products_bulk_imported"

type ProductsBulkImported struct {
	Base
	ProductIDs []uuid.UUID `json:"product_ids"`
}

func (ProductsBulkImported) isEvent()     {}
func (ProductsBulkImported) Type() string { return ProductsBulkImportedType }
func (e ProductsBulkImported) String() string {
	return fmt.Sprintf("Products bulk imported: %d", len(e.ProductIDs))
}

func (e *ProductsBulkImported) Apply(p *project.Project) {
	// Ensure we don't add duplicates if the event is replayed or input had duplicates.
	existing := make(map[uuid.UUID]struct{}, len(p.ProductIDs))
	for _, id := range p.ProductIDs {
		existing[id] = struct{}{}
	}
	for _, id := range e.ProductIDs {
		if id == uuid.Nil {
			continue
		}
		if _, ok := existing[id]; ok {
			continue
		}
		p.ProductIDs = append(p.ProductIDs, id)
		existing[id] = struct{}{}
	}
}

func (e *ProductsBulkImported) RelatedIDs() RelatedIDs {
	return RelatedIDs{}
}

func (e *ProductsBulkImported) NotificationMessage() string {
	return fmt.Sprintf("%d products have been added to the project.", len(e.ProductIDs))
}

type ProductsBulkImportedInput struct {
	ProductIDs []uuid.UUID `json:"product_ids"`
}

func DecideProductsBulkImported(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in ProductsBulkImportedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	if len(in.ProductIDs) == 0 {
		return nil, errors.New("product_ids is required")
	}

	// Validate and dedupe input; also reject ids already on the project.
	seen := map[uuid.UUID]struct{}{}
	unique := make([]uuid.UUID, 0, len(in.ProductIDs))

	for _, id := range in.ProductIDs {
		if id == uuid.Nil {
			return nil, errors.New("product_ids contains empty uuid")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		for _, existing := range cur.ProductIDs {
			if existing == id {
				return nil, fmt.Errorf("product %s already exists in project %s", id, cur.Id)
			}
		}

		unique = append(unique, id)
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = ProductsBulkImportedMeta.FriendlyName

	return &ProductsBulkImported{
		Base:       base,
		ProductIDs: unique,
	}, nil
}

var ProductsBulkImportedMeta = EventMeta{
	Type:         ProductsBulkImportedType,
	FriendlyName: "Bulk Product Import",
}

func init() {
	RegisterMeta(ProductsBulkImportedMeta, func() Event {
		return &ProductsBulkImported{
			Base: Base{FriendlyNameStr: ProductsBulkImportedMeta.FriendlyName},
		}
	})

	RegisterDecider[ProductsBulkImportedInput](ProductsBulkImportedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *project.Project, in ProductsBulkImportedInput, status Status) (Event, error) {
			return DecideProductsBulkImported(projectID, actor, cur, in, status)
		})

	RegisterInputType(ProductsBulkImportedType, ProductsBulkImportedInput{})
}
