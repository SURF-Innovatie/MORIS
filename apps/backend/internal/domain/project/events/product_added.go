package events

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/SURF-Innovatie/MORIS/internal/domain/project"
	"github.com/google/uuid"
)

const ProductAddedType = "project.product_added"

type ProductAdded struct {
	Base
	ProductID uuid.UUID `json:"product_id"`
}

func (ProductAdded) isEvent()     {}
func (ProductAdded) Type() string { return ProductAddedType }
func (e ProductAdded) String() string {
	return fmt.Sprintf("Product added: %s", e.ProductID)
}

func (e *ProductAdded) Apply(project *project.Project) {
	project.ProductIDs = append(project.ProductIDs, e.ProductID)
}

func (e *ProductAdded) RelatedIDs() RelatedIDs {
	return RelatedIDs{ProductID: &e.ProductID}
}

func (e *ProductAdded) NotificationMessage() string {
	return "A new product has been added to the project."
}

type ProductAddedInput struct {
	ProductID uuid.UUID `json:"product_id"`
}

func DecideProductAdded(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *project.Project,
	in ProductAddedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.ProductID == uuid.Nil {
		return nil, errors.New("product id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}

	if slices.Contains(cur.ProductIDs, in.ProductID) {
		return nil, fmt.Errorf("product %s already exists in project %s", in.ProductID, cur.Id)
	}

	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = ProductAddedMeta.FriendlyName

	return &ProductAdded{
		Base:      base,
		ProductID: in.ProductID,
	}, nil
}

var ProductAddedMeta = EventMeta{
	Type:         ProductAddedType,
	FriendlyName: "Product Addition",
}

func init() {
	RegisterMeta(ProductAddedMeta, func() Event {
		return &ProductAdded{
			Base: Base{FriendlyNameStr: ProductAddedMeta.FriendlyName},
		}
	})

	RegisterDecider[ProductAddedInput](ProductAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *project.Project, in ProductAddedInput, status Status) (Event, error) {
			return DecideProductAdded(projectID, actor, cur, in, status)
		})

	RegisterInputType(ProductAddedType, ProductAddedInput{})
}
