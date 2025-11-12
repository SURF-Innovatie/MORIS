package events

type PersonRemoved struct {
	Base
	Name string `json:"name"`
}

func (PersonRemoved) isEvent()     {}
func (PersonRemoved) Type() string { return PersonRemovedType }
