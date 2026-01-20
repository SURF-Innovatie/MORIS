package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	sourceDir    string
	outputFile   string
	dtoSourceDir string
)

type Field struct {
	Name string
	Type string
}

type EventStruct struct {
	Name         string
	EnumKey      string
	FunctionName string
	Fields       []Field
	EventType    string
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	flag.StringVar(&sourceDir, "source", ".", "Source directory to parse")
	flag.StringVar(&outputFile, "output", "events.ts", "Output file path")
	flag.StringVar(&dtoSourceDir, "dto-source", "", "Directory containing DTO structs")
	flag.Parse()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourceDir, func(fi os.FileInfo) bool {
		return !strings.HasPrefix(fi.Name(), "_")
	}, parser.ParseComments)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse directory")
	}

	// 1. Collect all constants
	constMap := make(map[string]string)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.GenDecl)
				if !ok || decl.Tok != token.CONST {
					return true
				}
				for _, spec := range decl.Specs {
					vspec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					// Only simple constants with values
					if len(vspec.Names) == len(vspec.Values) {
						for i, name := range vspec.Names {
							if val, ok := vspec.Values[i].(*ast.BasicLit); ok {
								constMap[name.Name] = strings.Trim(val.Value, "\"")
							}
						}
					}
				}
				return true
			})
		}
	}

	// 2. Find all RegisterInputType calls and map EventType -> StructName
	eventMap := make(map[string]string)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				isRegister := false
				if ident, ok := call.Fun.(*ast.Ident); ok {
					if ident.Name == "RegisterInputType" {
						isRegister = true
					}
				} else if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "RegisterInputType" {
						isRegister = true
					}
				}

				if !isRegister {
					return true
				}

				if len(call.Args) != 2 {
					return true
				}

				// Arg 0: Event Type String
				eventType := ""
				switch arg := call.Args[0].(type) {
				case *ast.BasicLit:
					eventType = strings.Trim(arg.Value, "\"")
				case *ast.Ident: // e.g., ProductRemovedType
					if val, ok := constMap[arg.Name]; ok {
						eventType = val
					} else {
						// Fallback: try to see if it looks like a constant string we can't resolve easily
						fmt.Printf("Warning: could not resolve constant %s\n", arg.Name)
					}
				}

				// Arg 1: Struct Instance e.g. ProductRemovedInput{}
				structName := ""
				switch arg := call.Args[1].(type) {
				case *ast.CompositeLit:
					if ident, ok := arg.Type.(*ast.Ident); ok {
						structName = ident.Name
					}
				}

				if eventType != "" && structName != "" {
					eventMap[structName] = eventType
				} else {
					fmt.Printf("Skipping potential match: Type='%s', Struct='%s'\n", eventType, structName)
				}
				return true
			})

		}
	}

	// 3. Extract Struct Definitions for Events
	structs := make([]EventStruct, 0)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				ts, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}
				structType, ok := ts.Type.(*ast.StructType)
				if !ok {
					return true
				}

				structName := ts.Name.Name
				eventType, exists := eventMap[structName]
				if !exists {
					return true
				}

				var fields []Field
				for _, f := range structType.Fields.List {
					if len(f.Names) == 0 {
						continue
					}
					fieldName := f.Names[0].Name

					// Get JSON tag
					jsonTag := fieldName
					if f.Tag != nil {
						tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
						if val, ok := tag.Lookup("json"); ok {
							jsonTag = strings.Split(val, ",")[0]
						}
					}

					fields = append(fields, Field{Name: jsonTag, Type: mapGoTypeToTS(f.Type)})
				}

				// Generate Function Name
				funcName := structName
				if strings.HasSuffix(structName, "Input") {
					funcName = strings.TrimSuffix(structName, "Input")
				}
				enumKey := funcName
				funcName += "Event"

				structs = append(structs, EventStruct{
					Name:         structName,
					EnumKey:      enumKey,
					FunctionName: funcName,
					Fields:       fields,
					EventType:    eventType,
				})

				return true
			})
		}
	}

	// 4. Extract Base Event Fields from DTO
	var baseEventFields []Field
	if dtoSourceDir != "" {
		dtoPkgs, err := parser.ParseDir(fset, dtoSourceDir, nil, parser.ParseComments)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse DTO directory")
		} else {
			for _, pkg := range dtoPkgs {
				for _, file := range pkg.Files {
					ast.Inspect(file, func(n ast.Node) bool {
						ts, ok := n.(*ast.TypeSpec)
						if !ok || ts.Name.Name != "Event" {
							return true
						}
						structType, ok := ts.Type.(*ast.StructType)
						if !ok {
							return true
						}

						for _, f := range structType.Fields.List {
							if len(f.Names) == 0 {
								continue
							}

							fieldOptional := false

							// Get JSON tag
							jsonTag := f.Names[0].Name
							if f.Tag != nil {
								tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
								if val, ok := tag.Lookup("json"); ok {
									parts := strings.Split(val, ",")
									jsonTag = parts[0]
									if len(parts) > 1 && parts[1] == "omitempty" {
										fieldOptional = true
									}
								}
							}

							// Check if pointer type for optionality
							tsType := mapGoTypeToTS(f.Type)
							if _, isPointer := f.Type.(*ast.StarExpr); isPointer {
								fieldOptional = true
							}

							if fieldOptional {
								jsonTag += "?"
							}

							baseEventFields = append(baseEventFields, Field{Name: jsonTag, Type: tsType})
						}
						return false // Stop inspecting this file/node once found
					})
				}
			}
		}
	}

	// Fallback if no fields found (or extraction failed)
	if len(baseEventFields) == 0 {
		baseEventFields = []Field{
			{Name: "id", Type: "string"},
			{Name: "projectId", Type: "string"},
			{Name: "status", Type: "string"},
			{Name: "createdBy", Type: "string"},
			{Name: "at", Type: "Datetime"},
			{Name: "details", Type: "string"},
			{Name: "projectTitle", Type: "string"},
		}
	}

	// 5. Generate TypeScript Code
	tmplStr := `// This file is generated by apps/backend/cmd/codegen/main.go. DO NOT EDIT.
import customInstance from "./custom-axios";
import { PersonResponse, ProductResponse, ProjectRoleResponse, BudgetCategory, FundingSource } from "./generated-orval/model";

export type Datetime = string;

export enum ProjectEventType {
{{range .Structs}}
  {{.EnumKey}} = '{{.EventType}}',
{{end}}
}

// Base Event Interface
export interface BaseEvent {
{{- range .BaseFields}}
  {{.Name}}: {{.Type}};
{{- end}}
}

// Specific Event Interfaces
{{range .Structs}}
export interface {{.FunctionName}} extends BaseEvent {
  type: ProjectEventType.{{.EnumKey}};
  data: {{.Name}};
}
{{end}}

// Union Type
export type ProjectEvent = 
{{- range .Structs}}
  | {{.FunctionName}}
{{- end}};

// Input Interfaces
{{range .Structs}}
export interface {{.Name}} {
{{- range .Fields}}
  {{.Name}}: {{.Type}};
{{- end}}
}

export const create{{.FunctionName}} = (projectId: string, input: {{.Name}}) => {
  return customInstance<any>({
    url: ` + "`" + `/projects/${projectId}/events` + "`" + `,
    method: 'POST',
    data: {
      projectId,
      type: '{{.EventType}}',
      input
    }
  });
};
{{end}}
`
	tmpl, err := template.New("ts").Parse(tmplStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse template")
	}

	data := struct {
		Structs    []EventStruct
		BaseFields []Field
	}{
		Structs:    structs,
		BaseFields: baseEventFields,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute template")
	}

	// Ensure directory exists
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create directory")
	}

	if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
		log.Fatal().Err(err).Msg("Failed to write to file")
	}

	fmt.Printf("Generated %s with %d events\n", outputFile, len(structs))
}

func mapGoTypeToTS(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return mapGoTypeToTS(t.X)
	case *ast.Ident:
		switch t.Name {
		case "string":
			return "string"
		case "int", "int64", "float64":
			return "number"
		case "bool":
			return "boolean"
		default:
			return t.Name // Fallback
		}
	case *ast.SelectorExpr:
		// Handle uuid.UUID, time.Time, dto.PersonResponse
		if x, ok := t.X.(*ast.Ident); ok {
			if x.Name == "uuid" && t.Sel.Name == "UUID" {
				return "string"
			} else if x.Name == "time" && t.Sel.Name == "Time" {
				return "Datetime"
			} else if t.Sel.Name == "Status" {
				return "string"
			}
		}
		return t.Sel.Name
	case *ast.ArrayType:
		if elt, ok := t.Elt.(*ast.Ident); ok {
			if elt.Name == "string" {
				return "string[]"
			}
		}
		return "any[]"
	}
	return "any"
}
