package events

import "github.com/SURF-Innovatie/MORIS/internal/common"

var inputTypes = map[string]any{}

func RegisterInputType(eventType string, input any) {
	inputTypes[eventType] = input
}

func GetInputSchema(eventType string) map[string]any {
	if in, ok := inputTypes[eventType]; ok {
		return common.StructToInputSchema(in)
	}
	return nil
}
