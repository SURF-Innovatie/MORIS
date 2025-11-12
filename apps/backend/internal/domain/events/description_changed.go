package events

type DescriptionChanged struct {
	Base
	Description string `json:"description"`
}

func (DescriptionChanged) isEvent()     {}
func (DescriptionChanged) Type() string { return DescriptionChangedType }
