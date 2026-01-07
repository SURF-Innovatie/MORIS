package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var (
	sourceDir  string
	outputFile string
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
	flag.StringVar(&sourceDir, "source", ".", "Source directory to parse")
	flag.StringVar(&outputFile, "output", "events.ts", "Output file path")
	flag.Parse()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourceDir, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse directory: %v", err)
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

	// 3. Extract Struct Definitions
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

					// Map Type
					tsType := "any"
					switch t := f.Type.(type) {
					case *ast.Ident:
						switch t.Name {
						case "string":
							tsType = "string"
						case "int", "int64", "float64":
							tsType = "number"
						case "bool":
							tsType = "boolean"
						default:
							tsType = t.Name // Fallback
						}
					case *ast.SelectorExpr:
						// Handle uuid.UUID, time.Time
						if x, ok := t.X.(*ast.Ident); ok {
							if x.Name == "uuid" && t.Sel.Name == "UUID" {
								tsType = "string"
							} else if x.Name == "time" && t.Sel.Name == "Time" {
								tsType = "Datetime"
							}
						}
					case *ast.ArrayType:
						// Handle arrays, e.g. []string
						if elt, ok := t.Elt.(*ast.Ident); ok {
							if elt.Name == "string" {
								tsType = "string[]"
							}
						}
					}

					fields = append(fields, Field{Name: jsonTag, Type: tsType})
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

	// 3. Generate TypeScript Code
	tmplStr := `// This file is generated by apps/backend/cmd/codegen/main.go. DO NOT EDIT.
import customInstance from "./custom-axios";

export type Datetime = string;

export enum ProjectEventType {
{{range .}}
  {{.EnumKey}} = '{{.EventType}}',
{{end}}
}

{{range .}}
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
		log.Fatalf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, structs); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Printf("Generated %s with %d events\n", outputFile, len(structs))
}
