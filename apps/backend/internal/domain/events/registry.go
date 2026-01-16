package events

import (
	"github.com/SURF-Innovatie/MORIS/internal/common"
	"github.com/samber/lo"
)

var inputTypes = map[string]any{}

func RegisterInputType(eventType string, input any) {
	inputTypes[eventType] = input
}

func GetRegisteredEventTypes() []string {
	return lo.Keys(inputTypes)
}

func GetInputSchema(eventType string) map[string]any {
	if in, ok := inputTypes[eventType]; ok {
		return common.StructToInputSchema(in)
	}
	return nil
}

func init() {
	RegisterMeta(CustomFieldValueSetMeta, func() Event {
		return &CustomFieldValueSet{
			Base: Base{FriendlyNameStr: CustomFieldValueSetMeta.FriendlyName},
		}
	})
	RegisterDecider(CustomFieldValueSetType, DecideCustomFieldValueSet)
	RegisterInputType(CustomFieldValueSetType, CustomFieldValueSetInput{})
}
