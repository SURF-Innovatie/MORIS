package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
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

func (e *ProductAdded) Apply(project *entities.Project) {
	project.ProductIDs = append(project.ProductIDs, e.ProductID)
}

func (e *ProductAdded) RelatedIDs() RelatedIDs {
	return RelatedIDs{ProductID: &e.ProductID}
}

func (e *ProductAdded) NotificationMessage() string {
	return "A new product has been added to the project."
}

type ProductAddedInput struct {
	ProductID uuid.UUID
}

func DecideProductAdded(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
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

	for _, x := range cur.ProductIDs {
		if x == in.ProductID {
			return nil, fmt.Errorf("product %s already exists in project %s", in.ProductID, cur.Id)
		}
	}

	return &ProductAdded{
		Base:      NewBase(projectID, actor, status),
		ProductID: in.ProductID,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         ProductAddedType,
		FriendlyName: "Product Addition",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &ProductAdded{} })

	RegisterDecider[ProductAddedInput](ProductAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur any, in ProductAddedInput, status Status) (Event, error) {
			p := cur.(*entities.Project)
			return DecideProductAdded(projectID, actor, p, in, status)
		})
}
