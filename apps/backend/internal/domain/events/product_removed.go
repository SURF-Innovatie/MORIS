package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

const ProductRemovedType = "project.product_removed"

type ProductRemoved struct {
	Base
	ProductID uuid.UUID `json:"product_id"`
}

func (ProductRemoved) isEvent()     {}
func (ProductRemoved) Type() string { return ProductRemovedType }
func (e ProductRemoved) String() string {
	return fmt.Sprintf("Product removed: %s", e.ProductID)
}

func (e *ProductRemoved) Apply(project *entities.Project) {
	shouldRemove := -1
	for i, p := range project.ProductIDs {
		if p == e.ProductID {
			shouldRemove = i
			break
		}
	}
	if shouldRemove != -1 {
		project.ProductIDs = append(project.ProductIDs[:shouldRemove], project.ProductIDs[shouldRemove+1:]...)
	}
}

func (e *ProductRemoved) RelatedIDs() RelatedIDs {
	return RelatedIDs{ProductID: &e.ProductID}
}

type ProductRemovedInput struct {
	ProductID uuid.UUID
}

func DecideProductRemoved(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in ProductRemovedInput,
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

	exist := false
	for _, x := range cur.ProductIDs {
		if x == in.ProductID {
			exist = true
			break
		}
	}
	if !exist {
		return nil, fmt.Errorf("product %s not found for project %s", in.ProductID, cur.Id)
	}

	return &ProductRemoved{
		Base:      NewBase(projectID, actor, status),
		ProductID: in.ProductID,
	}, nil
}

func init() {
	RegisterMeta(EventMeta{
		Type:         ProductRemovedType,
		FriendlyName: "Product Removal",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &ProductRemoved{} })

	RegisterDecider[ProductRemovedInput](ProductRemovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur any, in ProductRemovedInput, status Status) (Event, error) {
			p := cur.(*entities.Project)
			return DecideProductRemoved(projectID, actor, p, in, status)
		})
}
