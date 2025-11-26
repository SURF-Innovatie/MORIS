package events

import "fmt"

type DescriptionChanged struct {
	Base
	Description string `json:"description"`
}

func (DescriptionChanged) isEvent()     {}
func (DescriptionChanged) Type() string { return DescriptionChangedType }
func (e DescriptionChanged) String() string {
	return fmt.Sprintf("Description changed: %s", e.Description)
}
