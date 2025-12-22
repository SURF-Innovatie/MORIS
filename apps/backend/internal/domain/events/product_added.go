package events

import (
	"context"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

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

func init() {
	RegisterMeta(EventMeta{
		Type:         ProductAddedType,
		FriendlyName: "Product Addition",
		CheckNotification: func(ctx context.Context, event Event, client *ent.Client) bool {
			return true
		},
	}, func() Event { return &ProductAdded{} })
}
