package common

import (
	"reflect"
	"strings"
)

// StructToInputSchema returns a minimal schema:
// {"field":"string","other":"integer"}
func StructToInputSchema(v any) map[string]any {
	t := reflect.TypeOf(v)
	if t == nil {
		return map[string]any{}
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return map[string]any{}
	}

	out := make(map[string]any, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			continue
		}

		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		name := tag
		if before, _, ok := strings.Cut(tag, ","); ok {
			name = before
		}
		if name == "" {
			continue
		}

		out[name] = goTypeToSimpleType(f.Type)
	}

	return out
}

func goTypeToSimpleType(t reflect.Type) any {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// time.Time
	if t.PkgPath() == "time" && t.Name() == "Time" {
		return "datetime"
	}

	// uuid.UUID (github.com/google/uuid) -> "uuid"
	if t.PkgPath() == "github.com/google/uuid" && t.Name() == "UUID" {
		return "uuid"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "integer"

	case reflect.Bool:
		return "boolean"

	case reflect.Float32, reflect.Float64:
		return "number"

	case reflect.Slice, reflect.Array:
		return map[string]any{
			"type":  "array",
			"items": goTypeToSimpleType(t.Elem()),
		}

	case reflect.Map:
		return "object"

	case reflect.Struct:
		return StructToInputSchema(reflect.New(t).Interface())

	default:
		return "object"
	}
}
