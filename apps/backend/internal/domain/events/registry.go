package events

import (
	"github.com/SURF-Innovatie/MORIS/internal/common"
)

var inputTypes = map[string]any{}

func RegisterInputType(eventType string, input any) {
	inputTypes[eventType] = input
}

func GetRegisteredEventTypes() []string {
	keys := make([]string, 0, len(inputTypes))
	for k := range inputTypes {
		keys = append(keys, k)
	}
	return keys
}

func GetInputSchema(eventType string) map[string]any {
	if in, ok := inputTypes[eventType]; ok {
		return common.StructToInputSchema(in)
	}
	return nil
}

func init() {
	RegisterMeta(CustomFieldValueSetMeta, func() Event { return &CustomFieldValueSet{} })
	RegisterDecider(CustomFieldValueSetType, DecideCustomFieldValueSet)
	RegisterInputType(CustomFieldValueSetType, CustomFieldValueSetInput{})
}
