//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type FieldEvent struct {
	Type                 string `yaml:"type"`
	Field                string `yaml:"field"`
	FieldType            string `yaml:"field_type"`
	JSONKey              string `yaml:"json_key"`
	FriendlyName         string `yaml:"friendly_name"`
	NoOpOnEmpty          bool   `yaml:"no_op_on_empty"`
	CompareFunc          string `yaml:"compare_func"` // e.g., "Equal" for time.Time
	RelatedID            string `yaml:"related_id"`   // e.g., "OrgNodeID"
	RequireNonNil        bool   `yaml:"require_non_nil"`
	NotificationTemplate string `yaml:"notification_template"`
}

type EntityRefEvent struct {
	Type                 string `yaml:"type"`
	Field                string `yaml:"field"`
	JSONKey              string `yaml:"json_key"`
	FriendlyName         string `yaml:"friendly_name"`
	Entity               string `yaml:"entity"`     // e.g., "OrganisationNode"
	RelatedID            string `yaml:"related_id"` // e.g., "OrgNodeID"
	RequireNonNil        bool   `yaml:"require_non_nil"`
	NotificationTemplate string `yaml:"notification_template"`
}

type EntityCollectionEvent struct {
	Entity                     string `yaml:"entity"`
	IDField                    string `yaml:"id_field"`
	SliceField                 string `yaml:"slice_field"`
	JSONKey                    string `yaml:"json_key"`
	AddFriendlyName            string `yaml:"add_friendly_name"`
	RemoveFriendlyName         string `yaml:"remove_friendly_name"`
	RelatedID                  string `yaml:"related_id"`
	AddNotificationTemplate    string `yaml:"add_notification_template"`
	RemoveNotificationTemplate string `yaml:"remove_notification_template"`
}

type Config struct {
	FieldEvents            []FieldEvent            `yaml:"field_events"`
	EntityRefEvents        []EntityRefEvent        `yaml:"entity_ref_events"`
	EntityCollectionEvents []EntityCollectionEvent `yaml:"entity_collection_events"`
}

func main() {
	configPath := filepath.Join("gen", "events.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse config: %v\n", err)
		os.Exit(1)
	}

	// Generate field events (including entity ref events which use same template)
	allFieldEvents := cfg.FieldEvents
	for _, e := range cfg.EntityRefEvents {
		allFieldEvents = append(allFieldEvents, FieldEvent{
			Type:                 e.Type,
			Field:                e.Field,
			FieldType:            "uuid.UUID",
			JSONKey:              e.JSONKey,
			FriendlyName:         e.FriendlyName,
			RelatedID:            e.RelatedID,
			RequireNonNil:        e.RequireNonNil,
			NotificationTemplate: e.NotificationTemplate,
		})
	}

	if err := generateFieldEvents(allFieldEvents); err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate field events: %v\n", err)
		os.Exit(1)
	}

	// Generate entity collection events (add/remove)
	if err := generateEntityCollectionEvents(cfg.EntityCollectionEvents); err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate entity collection events: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated events successfully")
}

func generateFieldEvents(events []FieldEvent) error {
	tmpl := template.Must(template.New("field").Funcs(template.FuncMap{
		"eventName":   eventName,
		"constName":   constName,
		"metaName":    metaName,
		"inputName":   inputName,
		"decideName":  decideName,
		"needsTime":   needsTime,
		"needsUUID":   needsUUID,
		"compareExpr": compareExpr,
		"formatExpr":  formatExpr,
		"lower":       strings.ToLower,
	}).Parse(fieldEventTemplate))

	var buf bytes.Buffer
	buf.WriteString(fieldHeader(events))

	for _, e := range events {
		if err := tmpl.Execute(&buf, e); err != nil {
			return err
		}
	}

	// Write init function
	buf.WriteString("\nfunc init() {\n")
	for _, e := range events {
		buf.WriteString(fmt.Sprintf(`	RegisterMeta(%s, func() Event {
		return &%s{Base: Base{FriendlyNameStr: %s.FriendlyName}}
	})
	RegisterDecider[%s](%s,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in %s, status Status) (Event, error) {
			return %s(projectID, actor, cur, in, status)
		})
	RegisterInputType(%s, %s{})
`, metaName(e.Type), eventName(e.Type), metaName(e.Type),
			inputName(e.Type), constName(e.Type), inputName(e.Type), decideName(e.Type),
			constName(e.Type), inputName(e.Type)))
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Write unformatted for debugging
		os.WriteFile("field_events_gen.go", buf.Bytes(), 0644)
		return fmt.Errorf("format error: %w", err)
	}

	return os.WriteFile("field_events_gen.go", formatted, 0644)
}

func generateEntityCollectionEvents(events []EntityCollectionEvent) error {
	tmpl := template.Must(template.New("entity").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).Parse(entityEventTemplate))

	var buf bytes.Buffer
	buf.WriteString(entityHeader())

	for _, e := range events {
		if err := tmpl.Execute(&buf, e); err != nil {
			return err
		}
	}

	// Write init function
	buf.WriteString("\nfunc init() {\n")
	for _, e := range events {
		buf.WriteString(fmt.Sprintf(`	RegisterMeta(%sAddedMeta, func() Event {
		return &%sAdded{Base: Base{FriendlyNameStr: %sAddedMeta.FriendlyName}}
	})
	RegisterDecider[%sAddedInput](%sAddedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in %sAddedInput, status Status) (Event, error) {
			return Decide%sAdded(projectID, actor, cur, in, status)
		})
	RegisterInputType(%sAddedType, %sAddedInput{})

	RegisterMeta(%sRemovedMeta, func() Event {
		return &%sRemoved{Base: Base{FriendlyNameStr: %sRemovedMeta.FriendlyName}}
	})
	RegisterDecider[%sRemovedInput](%sRemovedType,
		func(ctx context.Context, projectID uuid.UUID, actor uuid.UUID, cur *entities.Project, in %sRemovedInput, status Status) (Event, error) {
			return Decide%sRemoved(projectID, actor, cur, in, status)
		})
	RegisterInputType(%sRemovedType, %sRemovedInput{})
`, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity,
			e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity, e.Entity))
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		os.WriteFile("entity_events_gen.go", buf.Bytes(), 0644)
		return fmt.Errorf("format error: %w", err)
	}

	return os.WriteFile("entity_events_gen.go", formatted, 0644)
}

// Helper functions for templates
func eventName(typ string) string {
	parts := strings.Split(typ, ".")
	name := parts[len(parts)-1]
	return toPascalCase(name)
}

func constName(typ string) string {
	return eventName(typ) + "Type"
}

func metaName(typ string) string {
	return eventName(typ) + "Meta"
}

func inputName(typ string) string {
	return eventName(typ) + "Input"
}

func decideName(typ string) string {
	return "Decide" + eventName(typ)
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func needsTime(events []FieldEvent) bool {
	for _, e := range events {
		if e.FieldType == "time.Time" {
			return true
		}
	}
	return false
}

func needsUUID(events []FieldEvent) bool {
	for _, e := range events {
		if e.FieldType == "uuid.UUID" || e.RequireNonNil {
			return true
		}
	}
	return true // Always need for projectID/actor
}

func compareExpr(e FieldEvent) string {
	if e.CompareFunc != "" {
		return fmt.Sprintf("cur.%s.%s(in.%s)", e.Field, e.CompareFunc, e.Field)
	}
	return fmt.Sprintf("cur.%s == in.%s", e.Field, e.Field)
}

func formatExpr(e FieldEvent) string {
	if e.FieldType == "time.Time" {
		return fmt.Sprintf(`e.%s.Format("2006-01-02")`, e.Field)
	}
	return fmt.Sprintf("e.%s", e.Field)
}

func fieldHeader(events []FieldEvent) string {
	imports := []string{
		`"context"`,
		`"errors"`,
		`"fmt"`,
		`"github.com/SURF-Innovatie/MORIS/internal/domain/entities"`,
		`"github.com/google/uuid"`,
	}
	if needsTime(events) {
		imports = append(imports, `"time"`)
	}
	return fmt.Sprintf(`// Code generated by gen/generator.go. DO NOT EDIT.

package events

import (
	%s
)

`, strings.Join(imports, "\n\t"))
}

func entityHeader() string {
	return `// Code generated by gen/generator.go. DO NOT EDIT.

package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

`
}

const fieldEventTemplate = `
// --- {{eventName .Type}} ---

const {{constName .Type}} = "{{.Type}}"

type {{eventName .Type}} struct {
	Base
	{{.Field}} {{.FieldType}} ` + "`json:\"{{.JSONKey}}\"`" + `
}

func ({{eventName .Type}}) isEvent()     {}
func ({{eventName .Type}}) Type() string { return {{constName .Type}} }
func (e {{eventName .Type}}) String() string {
	return "{{.FriendlyName}}: " + fmt.Sprint({{formatExpr .}})
}

func (e *{{eventName .Type}}) Apply(project *entities.Project) {
	project.{{.Field}} = e.{{.Field}}
}

func (e *{{eventName .Type}}) NotificationMessage() string {
	return "Project {{.Field | lower}} has been updated."
}
{{if .NotificationTemplate}}
func (e *{{eventName .Type}}) NotificationTemplate() string {
	return "{{.NotificationTemplate}}"
}

func (e *{{eventName .Type}}) NotificationVariables() map[string]string {
	return map[string]string{
		"event.{{.Field}}": fmt.Sprint(e.{{.Field}}),
	}
}
{{end}}{{if .RelatedID}}
func (e *{{eventName .Type}}) RelatedIDs() RelatedIDs {
	return RelatedIDs{ {{.RelatedID}}: &e.{{.Field}} }
}
{{end}}
type {{inputName .Type}} struct {
	{{.Field}} {{.FieldType}} ` + "`json:\"{{.JSONKey}}\"`" + `
}

func {{decideName .Type}}(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in {{inputName .Type}},
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
{{if .RequireNonNil}}	if in.{{.Field}} == uuid.Nil {
		return nil, errors.New("{{.JSONKey}} is required")
	}
{{end}}{{if .NoOpOnEmpty}}	if in.{{.Field}} == "" {
		return nil, nil
	}
{{end}}	if {{compareExpr .}} {
		return nil, nil
	}
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = {{metaName .Type}}.FriendlyName

	return &{{eventName .Type}}{
		Base:  base,
		{{.Field}}: in.{{.Field}},
	}, nil
}

var {{metaName .Type}} = EventMeta{
	Type:         {{constName .Type}},
	FriendlyName: "{{.FriendlyName}}",
}
`

const entityEventTemplate = `
// --- {{.Entity}}Added / {{.Entity}}Removed ---

const {{.Entity}}AddedType = "project.{{.Entity | lower}}_added"
const {{.Entity}}RemovedType = "project.{{.Entity | lower}}_removed"

type {{.Entity}}Added struct {
	Base
	{{.IDField}} uuid.UUID ` + "`json:\"{{.JSONKey}}\"`" + `
}

func ({{.Entity}}Added) isEvent()     {}
func ({{.Entity}}Added) Type() string { return {{.Entity}}AddedType }
func (e {{.Entity}}Added) String() string {
	return fmt.Sprintf("{{.Entity}} added: %s", e.{{.IDField}})
}

func (e *{{.Entity}}Added) Apply(project *entities.Project) {
	project.{{.SliceField}} = append(project.{{.SliceField}}, e.{{.IDField}})
}

func (e *{{.Entity}}Added) RelatedIDs() RelatedIDs {
	return RelatedIDs{ {{.RelatedID}}: &e.{{.IDField}} }
}

func (e *{{.Entity}}Added) NotificationMessage() string {
	return "A new {{.Entity | lower}} has been added to the project."
}
{{if .AddNotificationTemplate}}
func (e *{{.Entity}}Added) NotificationTemplate() string {
	return "{{.AddNotificationTemplate}}"
}

func (e *{{.Entity}}Added) NotificationVariables() map[string]string {
	return map[string]string{
		"event.{{.IDField}}": e.{{.IDField}}.String(),
	}
}
{{end}}

type {{.Entity}}AddedInput struct {
	{{.IDField}} uuid.UUID ` + "`json:\"{{.JSONKey}}\"`" + `
}

func Decide{{.Entity}}Added(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in {{.Entity}}AddedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.{{.IDField}} == uuid.Nil {
		return nil, errors.New("{{.JSONKey}} is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	for _, x := range cur.{{.SliceField}} {
		if x == in.{{.IDField}} {
			return nil, fmt.Errorf("{{.Entity | lower}} %s already exists", in.{{.IDField}})
		}
	}
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = {{.Entity}}AddedMeta.FriendlyName

	return &{{.Entity}}Added{
		Base:      base,
		{{.IDField}}: in.{{.IDField}},
	}, nil
}

var {{.Entity}}AddedMeta = EventMeta{
	Type:         {{.Entity}}AddedType,
	FriendlyName: "{{.AddFriendlyName}}",
}

// --- {{.Entity}}Removed ---

type {{.Entity}}Removed struct {
	Base
	{{.IDField}} uuid.UUID ` + "`json:\"{{.JSONKey}}\"`" + `
}

func ({{.Entity}}Removed) isEvent()     {}
func ({{.Entity}}Removed) Type() string { return {{.Entity}}RemovedType }
func (e {{.Entity}}Removed) String() string {
	return fmt.Sprintf("{{.Entity}} removed: %s", e.{{.IDField}})
}

func (e *{{.Entity}}Removed) Apply(project *entities.Project) {
	for i, x := range project.{{.SliceField}} {
		if x == e.{{.IDField}} {
			project.{{.SliceField}} = append(project.{{.SliceField}}[:i], project.{{.SliceField}}[i+1:]...)
			return
		}
	}
}

func (e *{{.Entity}}Removed) RelatedIDs() RelatedIDs {
	return RelatedIDs{ {{.RelatedID}}: &e.{{.IDField}} }
}
{{if .RemoveNotificationTemplate}}
func (e *{{.Entity}}Removed) NotificationMessage() string {
	return "A {{.Entity | lower}} has been removed from the project."
}

func (e *{{.Entity}}Removed) NotificationTemplate() string {
	return "{{.RemoveNotificationTemplate}}"
}

func (e *{{.Entity}}Removed) NotificationVariables() map[string]string {
	return map[string]string{
		"event.{{.IDField}}": e.{{.IDField}}.String(),
	}
}
{{end}}

type {{.Entity}}RemovedInput struct {
	{{.IDField}} uuid.UUID ` + "`json:\"{{.JSONKey}}\"`" + `
}

func Decide{{.Entity}}Removed(
	projectID uuid.UUID,
	actor uuid.UUID,
	cur *entities.Project,
	in {{.Entity}}RemovedInput,
	status Status,
) (Event, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project id is required")
	}
	if in.{{.IDField}} == uuid.Nil {
		return nil, errors.New("{{.JSONKey}} is required")
	}
	if cur == nil {
		return nil, errors.New("current project is required")
	}
	found := false
	for _, x := range cur.{{.SliceField}} {
		if x == in.{{.IDField}} {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("{{.Entity | lower}} %s not found", in.{{.IDField}})
	}
	base := NewBase(projectID, actor, status)
	base.FriendlyNameStr = {{.Entity}}RemovedMeta.FriendlyName

	return &{{.Entity}}Removed{
		Base:      base,
		{{.IDField}}: in.{{.IDField}},
	}, nil
}

var {{.Entity}}RemovedMeta = EventMeta{
	Type:         {{.Entity}}RemovedType,
	FriendlyName: "{{.RemoveFriendlyName}}",
}
`
