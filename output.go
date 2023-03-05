package gomapper

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"
)

//go:embed output.go.tmpl
var outputTemplateRaw string

func Output(writer io.Writer, config *Config) {
	tmpl := template.Must(template.New("output").Parse(outputTemplateRaw))

	type Func struct {
		Arg    string
		Return string
		Body   []string
	}
	type TemplateData struct {
		Editable bool
		Imports  map[string]string

		MappingFuncs map[string]Func // func name -> func data
	}

	type FromTo struct {
		From, To string
	}
	imports := make(map[string]string, len(config.Mappings))
	// destination type -> source type -> source+destination fields
	mappingByTypeName := make(map[string]map[string][]FromTo)
	for _, mapping := range config.Mappings {
		if _, ok := imports[mapping.Source.Package]; !ok {
			imports[mapping.Source.Package] = formatImport(mapping.Source.Package)
			imports[mapping.Destination.Package] = formatImport(mapping.Destination.Package)
		}

		if mappingByTypeName[mapping.Destination.Name] == nil {
			mappingByTypeName[mapping.Destination.Name] = make(map[string][]FromTo)
		}

		fieldMappings := mappingByTypeName[mapping.Destination.Name]
		fieldMappings[mapping.Source.Name] = append(fieldMappings[mapping.Source.Name], FromTo{
			From: mapping.Source.Field,
			To:   mapping.Destination.Field,
		})
	}

	funcBodyByFuncName := make(map[string]Func, len(config.Mappings))
	for source, dests := range mappingByTypeName {
		for dest, mappingFields := range dests {
			builder := funcBuilder{}

			builder.Append(fmt.Sprintf("return %s{", dest))

			for _, fromTo := range mappingFields {
				builder.Append(fmt.Sprintf("\t%s: sour.%s,", fromTo.From, fromTo.To))
			}

			builder.Append("}")

			// todo: handle func conflicts
			funcName := fmt.Sprintf("Map%sTo%s", source, dest)
			funcBodyByFuncName[funcName] = Func{
				Arg:    fmt.Sprintf("sour %s", source),
				Return: dest,
				Body:   builder.Lines(),
			}
		}
	}

	data := TemplateData{
		Editable:     false,
		Imports:      imports,
		MappingFuncs: funcBodyByFuncName,
	}

	if err := tmpl.Execute(writer, data); err != nil {
		panic(err)
	}
}

type funcBuilder struct {
	lines []string
}

func (b *funcBuilder) Append(s string) {
	b.lines = append(b.lines, s)
}

func (b *funcBuilder) NewLine() {
	b.lines = append(b.lines, "")
}

func (b *funcBuilder) Lines() []string {
	return b.lines
}

func formatImport(in string) string {
	sourcePackages := strings.Split(in, "/")
	out := sourcePackages[len(sourcePackages)-1]
	out = strings.ReplaceAll(out, "-", "_")
	out = strings.ReplaceAll(out, ".", "_")
	return out
}
