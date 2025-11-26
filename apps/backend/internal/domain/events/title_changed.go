package events

import "fmt"

type TitleChanged struct {
	Base
	Title string `json:"title"`
}

func (TitleChanged) isEvent()     {}
func (TitleChanged) Type() string { return TitleChangedType }
func (e TitleChanged) String() string {
	return fmt.Sprintf("Title changed: %s", e.Title)
}
