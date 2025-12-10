package events

import (
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type ProductRemoved struct {
	Base
	Product entities.Product `json:"product"`
}

func (ProductRemoved) isEvent()     {}
func (ProductRemoved) Type() string { return ProductRemovedType }
func (e ProductRemoved) String() string {
	if e.Product.Name != "" {
		return fmt.Sprintf("Product removed: %s", e.Product.Name)
	}
	return fmt.Sprintf("Product removed: %s", e.Product.Id)
}
