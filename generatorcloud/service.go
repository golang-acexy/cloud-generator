package generatorcloud

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/acexy/golang-toolkit/util/str"
	"golang.org/x/tools/imports"
	"os"
	"path/filepath"
)

//go:embed tmpl/file/biz.gohtml
var serviceTmpl string

type ServiceGen struct {
	gen    *Generator
	models []string
}

type ServiceData struct {
	StructName string
	ParamName  string
	ModelPKG   string
	PKG        string
}
type ServiceConfig struct {
}

func NewServiceGen(gen *Generator, models []string) *ServiceGen {
	return &ServiceGen{
		gen:    gen,
		models: models,
	}
}

func (s *ServiceGen) Create() {
	coll.SliceForeachAll(s.models, func(model string) {
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
		_ = os.MkdirAll(dir, os.ModePerm)
		servicePath = filepath.Join(dir, str.CamelToSnake(str.LowFirstChar(model))+"_biz.go")
		//判断文件是否存在
		if _, err := os.Stat(servicePath); err == nil {
			fmt.Println(model, "已有biz文件 略过生成")
			return
		} else {
			fmt.Println("生成biz文件", model, servicePath)
		}
		var buf bytes.Buffer
		_ = s.gen.render(serviceTmpl, &buf, ServiceData{
			StructName: model,
			ParamName:  str.LowFirstChar(model),
			ModelPKG:   s.gen.modelPkg,
			PKG:        pkg,
		})
		result, _ := imports.Process(servicePath, buf.Bytes(), nil)
		_ = os.WriteFile(servicePath, result, os.ModePerm)
	})
}
