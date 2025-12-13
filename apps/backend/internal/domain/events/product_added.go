package events

import (
	"fmt"

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
