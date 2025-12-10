package events

import (
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

type ProductAdded struct {
	Base
	Product entities.Product `json:"product"`
}

func (ProductAdded) isEvent()     {}
func (ProductAdded) Type() string { return ProductAddedType }
func (e ProductAdded) String() string {
	if e.Product.Name != "" {
		return fmt.Sprintf("Product added: %s", e.Product.Name)
	}
	return fmt.Sprintf("Product added: %s", e.Product.Id)
}
