package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

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

func init() {
	RegisterMeta(EventMeta{
		Type:         ProductRemovedType,
		FriendlyName: "Product Removal",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &ProductRemoved{} })
}
