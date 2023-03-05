// Package loader provides configuration loading functionality
package loader

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/template"

	"github.com/go-mapper/gomapper"
	"github.com/go-mapper/gomapper/internal/logger"
	"golang.org/x/tools/go/packages"
)

var (
	//go:embed templates/main.go.tmpl
	loaderGoTmpl string

	loaderTemplate *template.Template
	templateLoader sync.Once
)

func Load(path string, buildFlags []string) (*gomapper.Config, error) {
	templateLoader.Do(func() {
		loaderTemplate = template.Must(template.New("main").Parse(loaderGoTmpl))
	})

	info, err := parsePackage(path, buildFlags)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	type TemplateData struct {
		PkgInfo PackageInfo
	}

	out := new(bytes.Buffer)

	if err := loaderTemplate.Execute(out, TemplateData{
		PkgInfo: info,
	}); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	result, err := format.Source(out.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format source: %w", err)
	}

	logger.Log("new template: \n%s", out.String())

	exists, err := checkExists(".gomapper")
	if err != nil {
		return nil, fmt.Errorf("exist: %w", err)
	}

	if !exists {
		if err := os.Mkdir(".gomapper", os.ModePerm); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("mkdir: %w", err)
		}
	}

	f, err := os.Create(".gomapper/main.go")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
		_ = os.Remove(f.Name())
	}()

	_, err = f.Write(result)
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	stdout, err := run(f.Name())
	if err != nil {
		return nil, err
	}

	config := new(gomapper.Config)
	if err := json.NewDecoder(strings.NewReader(stdout)).Decode(config); err != nil {
		return nil, fmt.Errorf("parse mappings: %w", err)
	}

	return config, nil
}

type PackageInfo struct {
	PkgPath  string
	FuncName string
}

const (
	configureMapperFunc = "ConfigureMapper"
)

func parsePackage(path string, buildFlags []string) (info PackageInfo, _ error) {
	pkgs, err := packages.Load(&packages.Config{
		BuildFlags: buildFlags,
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule,
	}, path)
	if err != nil {
		return info, fmt.Errorf("load package: %w", err)
	}
	if len(pkgs) != 1 {
		return info, fmt.Errorf("can not find package")
	}
	mapperPkg := pkgs[0]

	if len(mapperPkg.Errors) != 0 {
		return info, fmt.Errorf("pkg errors: %s", mapperPkg.Errors[0].Error())
	}

	f := mapperPkg.Types.Scope().Lookup(configureMapperFunc)
	if f == nil {
		return info, fmt.Errorf("can not find %s function", configureMapperFunc)
	}

	return PackageInfo{
		PkgPath:  f.Pkg().Path(),
		FuncName: f.Name(),
	}, nil
}

// run the go application and return stdout.
func run(path string) (string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command("go", "run", path)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go run: %w: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// checkExists returns whether the given file or directory checkExists
func checkExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
