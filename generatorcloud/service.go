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

type ServiceGen struct {
	gen    *Generator
	models []string
}

type ServiceData struct {
	StructName             string
	ParamName              string
	ModelPKG               string
	RepoPKG                string
	PKG                    string
	IDType                 string
	DefaultOrderBy         string
	MaxQuerySize           int
	DefaultTimeRangeField  string
	AllowedTimeRangeFields []string
}

func NewServiceGen(gen *Generator, models []string) *ServiceGen {
	return &ServiceGen{
		gen:    gen,
		models: models,
	}
}

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
		//判断文件是否存在
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
