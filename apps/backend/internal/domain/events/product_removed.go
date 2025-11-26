package events

import (
	"fmt"

	"github.com/google/uuid"
)

type ProductRemoved struct {
	Base
	ProductID uuid.UUID `json:"productID"`
}

func (ProductRemoved) isEvent()     {}
func (ProductRemoved) Type() string { return ProductRemovedType }
func (e ProductRemoved) String() string {
	return fmt.Sprintf("Person removed: %s", e.ProductID)
}
