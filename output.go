package mapper

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

	imports := make(map[string]string, len(config.mappings))
	for source, dests := range config.mappings {

		sourcePackages := strings.Split(source.Package, "/")
		sourceAlias := sourcePackages[len(sourcePackages)-1]
		sourceAlias = strings.ReplaceAll(sourceAlias, "-", "_")
		sourceAlias = strings.ReplaceAll(sourceAlias, ".", "_")
		imports[source.Package] = sourceAlias

		for dest := range dests {
			destPackages := strings.Split(dest.Package, "/")
			destAlias := destPackages[len(destPackages)-1]
			destAlias = strings.ReplaceAll(destAlias, "-", "_")
			destAlias = strings.ReplaceAll(destAlias, ".", "_")
			imports[dest.Package] = destAlias
		}
	}

	funcToBody := make(map[string]Func, len(config.mappings))
	for source, dests := range config.mappings {
		for dest, mappings := range dests {
			builder := funcBuilder{}

			// TODO: format code
			builder.Append(fmt.Sprintf("return %s{", dest.Name))

			for _, m := range mappings {
				builder.Append(fmt.Sprintf("\t%s: sour.%s", m.Destination, m.Source))
			}

			builder.Append("}")

			// todo: handle func conflicts
			funcName := fmt.Sprintf("Map%sTo%s", source.Name, dest.Name)
			funcToBody[funcName] = Func{
				Arg:    fmt.Sprintf("sour %s", source.Name),
				Return: dest.Name,
				Body:   builder.Lines(),
			}
		}
	}

	data := TemplateData{
		Editable:     false,
		Imports:      imports,
		MappingFuncs: funcToBody,
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
