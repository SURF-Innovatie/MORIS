package events

type TitleChanged struct {
	Base
	Title string `json:"title"`
}

func (TitleChanged) isEvent()     {}
func (TitleChanged) Type() string { return TitleChangedType }
