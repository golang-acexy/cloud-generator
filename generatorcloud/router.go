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

//go:embed tmpl/file/router.gohtml
var routerImpl string

type RouterGen struct {
	gen    *Generator
	config map[string]*RouterConfig
}

type RouterData struct {
	ModelStructName     string
	ParamName           string
	PkgName             string
	AuthorityFetchCode  string
	DataLimitStructName string
	GroupPath           string
	DisableBaseHandler  bool
}

func NewRouterGen(gen *Generator, config map[string]*RouterConfig) *RouterGen {
	return &RouterGen{
		gen:    gen,
		config: config,
	}
}

func (s *RouterGen) Create() {
	coll.MapForeachAll(s.config, func(structName string, config *RouterConfig) {
		if config.BaseRouter != nil {
			dir := s.gen.baseOutput
			dir = filepath.Join(append([]string{dir}, config.BaseRouter.RelativeModelPath...)...)
			pkg := config.BaseRouter.RelativeModelPath[len(config.BaseRouter.RelativeModelPath)-1]
			_ = os.MkdirAll(dir, os.ModePerm)
			filePath := filepath.Join(dir, config.BaseRouter.FilePrefix, str.CamelToSnake(str.LowFirstChar(structName))+"_router.go")
			// 判断文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println(structName, "已有base router文件 略过生成")
				goto R
			} else {
				fmt.Println("生成base router文件", structName, filePath)
			}
			var buf bytes.Buffer
			_ = s.gen.render(routerImpl, &buf, RouterData{
				ModelStructName: structName,
				ParamName:       str.LowFirstChar(structName),
				PkgName:         pkg,
				GroupPath:       config.BaseRouter.GroupPath,
			})
			result, _ := imports.Process(filePath, buf.Bytes(), nil)
			_ = os.WriteFile(filePath, result, os.ModePerm)
		}
	R:
		if config.BaseRouterWithDataCheck != nil {
			dir := s.gen.baseOutput
			dir = filepath.Join(append([]string{dir}, config.BaseRouterWithDataCheck.RelativeModelPath...)...)
			pkg := config.BaseRouterWithDataCheck.RelativeModelPath[len(config.BaseRouterWithDataCheck.RelativeModelPath)-1]
			_ = os.MkdirAll(dir, os.ModePerm)
			filePath := filepath.Join(dir, config.BaseRouterWithDataCheck.FilePrefix, str.CamelToSnake(str.LowFirstChar(structName))+"_router.go")
			// 判断文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println(structName, "已有router with data limit文件 略过生成")
				return
			} else {
				fmt.Println("生成router with data limit文件", structName, filePath)
			}
			var buf bytes.Buffer
			_ = s.gen.render(routerImpl, &buf, RouterData{
				ModelStructName:     structName,
				ParamName:           str.LowFirstChar(structName),
				PkgName:             pkg,
				GroupPath:           config.BaseRouterWithDataCheck.GroupPath,
				DataLimitStructName: config.BaseRouterWithDataCheck.DataLimitStructName,
				AuthorityFetchCode:  config.BaseRouterWithDataCheck.AuthorityFetchCode,
				DisableBaseHandler:  config.BaseRouterWithDataCheck.DisableBaseHandler,
			})
			result, _ := imports.Process(filePath, buf.Bytes(), nil)
			_ = os.WriteFile(filePath, result, os.ModePerm)
		}
	})
}
