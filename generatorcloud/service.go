package generatorcloud

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/acexy/golang-toolkit/util/str"
	"golang.org/x/tools/imports"
)

//go:embed tmpl/file/biz.gohtml
var serviceTmpl string

// ServiceGen generates missing service files using previously collected model metadata.
type ServiceGen struct {
	gen    *Generator
	models []string
}

// ServiceData holds business service template inputs; expressions and paths must be trusted.
type ServiceData struct {
	// StructName is the generated Go model name.
	StructName string
	// ParamName is the lower-initial name used for generated package variables.
	ParamName string
	// ModelPKG is the model import path.
	ModelPKG string
	// RepoPKG is the repository import path.
	RepoPKG string
	// PKG is the output package name.
	PKG string
	// IDType is the model ID type expression; accepted metadata does not guarantee all templates compile.
	IDType string
	// DefaultOrderBy supplies trusted SQL ordering; empty uses id desc.
	DefaultOrderBy string
	// MaxQuerySize caps ordinary generated list queries; nonpositive values use 500 and do not cap pages.
	MaxQuerySize int
	// DefaultTimeRangeField is normalized to snake case and must appear in AllowedTimeRangeFields.
	DefaultTimeRangeField string
	// AllowedTimeRangeFields lists allowed time columns; existence in every model is not checked.
	AllowedTimeRangeFields []string
}

// NewServiceGen borrows a generator and model list; ID metadata must already have been collected.
func NewServiceGen(gen *Generator, models []string) *ServiceGen {
	return &ServiceGen{
		gen:    gen,
		models: models,
	}
}

// Create writes missing service files and skips existing paths; partial output is retained on failure.
func (s *ServiceGen) Create() error {
	defaultOrderBy := "id desc"
	maxQuerySize := 500
	if s.gen.serviceBase != nil {
		if s.gen.serviceBase.DefaultOrderBy != "" {
			defaultOrderBy = s.gen.serviceBase.DefaultOrderBy
		}
		if s.gen.serviceBase.MaxQuerySize > 0 {
			maxQuerySize = s.gen.serviceBase.MaxQuerySize
		}
	}
	for _, model := range s.models {
		repoPKG := path.Join(s.gen.modelPkg, "repo")
		if len(s.gen.repoRelativeModelPath) > 0 {
			repoPKG = path.Join(append([]string{s.gen.modelPkg}, s.gen.repoRelativeModelPath...)...)
		}
		dir := s.gen.baseOutput
		var servicePath string
		var pkg string
		if len(s.gen.serviceRelativeModelPath) > 0 {
			dir = filepath.Join(append([]string{dir}, s.gen.serviceRelativeModelPath...)...)
			pkg = s.gen.serviceRelativeModelPath[len(s.gen.serviceRelativeModelPath)-1]
		} else {
			pkg = "biz"
			dir = filepath.Join(dir, "biz")
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		servicePath = filepath.Join(dir, str.CamelToSnake(str.LowFirstChar(model))+"_biz.go")

		if _, err := os.Stat(servicePath); err == nil {
			fmt.Println(model, "已有biz文件 略过生成")
			continue
		} else {
			fmt.Println("生成biz文件", model, servicePath)
		}
		var buf bytes.Buffer
		if err := s.gen.render(serviceTmpl, &buf, ServiceData{
			StructName:             model,
			ParamName:              str.LowFirstChar(model),
			ModelPKG:               s.gen.modelPkg,
			RepoPKG:                repoPKG,
			PKG:                    pkg,
			IDType:                 s.gen.modelIDTypes[model],
			DefaultOrderBy:         defaultOrderBy,
			MaxQuerySize:           maxQuerySize,
			DefaultTimeRangeField:  s.gen.modelBase.DefaultTimeRangeField,
			AllowedTimeRangeFields: s.gen.modelBase.AllowedTimeRangeFields,
		}); err != nil {
			return err
		}
		result, err := imports.Process(servicePath, buf.Bytes(), nil)
		if err != nil {
			return err
		}
		if err = os.WriteFile(servicePath, result, 0644); err != nil {
			return err
		}
	}
	return nil
}
