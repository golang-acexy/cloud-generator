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

//go:embed tmpl/file/router.gohtml
var routerImpl string

type RouterGen struct {
	gen    *Generator
	config map[string]*RouterConfig
}

type RouterData struct {
	ModelStructName      string
	IDType               string
	WithAuthority        bool
	ParamName            string
	PkgName              string
	ModelPKG             string
	BizPKG               string
	AuthorityFetchCode   string
	AuthorityStructField string
	AuthorityColumn      string
	GroupPath            string
	DisableBaseHandler   bool
}

func NewRouterGen(gen *Generator, config map[string]*RouterConfig) *RouterGen {
	return &RouterGen{
		gen:    gen,
		config: config,
	}
}

func (s *RouterGen) Create() error {
	bizPKG := path.Join(s.gen.modelPkg, "biz")
	if len(s.gen.serviceRelativeModelPath) > 0 {
		bizPKG = path.Join(append([]string{s.gen.modelPkg}, s.gen.serviceRelativeModelPath...)...)
	}
	for structName, config := range s.config {
		if config.BaseRouter != nil {
			if len(config.BaseRouter.RelativeModelPath) == 0 {
				return fmt.Errorf("%w: %s", ErrInvalidRouterPath, structName)
			}
			dir := s.gen.baseOutput
			dir = filepath.Join(append([]string{dir}, config.BaseRouter.RelativeModelPath...)...)
			pkg := config.BaseRouter.RelativeModelPath[len(config.BaseRouter.RelativeModelPath)-1]
			filePath := filepath.Join(dir, config.BaseRouter.FilePrefix, str.CamelToSnake(str.LowFirstChar(structName))+"_router.go")
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return err
			}
			// 判断文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println(structName, "已有base router文件 略过生成")
			} else {
				fmt.Println("生成base router文件", structName, filePath)
				var buf bytes.Buffer
				if err = s.gen.render(routerImpl, &buf, RouterData{
					ModelStructName: structName,
					IDType:          s.gen.modelIDTypes[structName],
					ParamName:       str.LowFirstChar(structName),
					PkgName:         pkg,
					ModelPKG:        s.gen.modelPkg,
					BizPKG:          bizPKG,
					GroupPath:       config.BaseRouter.GroupPath,
				}); err != nil {
					return err
				}
				result, err := imports.Process(filePath, buf.Bytes(), nil)
				if err != nil {
					return err
				}
				if err = os.WriteFile(filePath, result, 0644); err != nil {
					return err
				}
			}
		}
		if config.BaseRouterWithAuthority != nil {
			if len(config.BaseRouterWithAuthority.RelativeModelPath) == 0 {
				return fmt.Errorf("%w: %s", ErrInvalidRouterPath, structName)
			}
			if config.BaseRouterWithAuthority.AuthorityFetchCode == "" ||
				config.BaseRouterWithAuthority.AuthorityStructField == "" ||
				config.BaseRouterWithAuthority.AuthorityColumn == "" {
				return fmt.Errorf("%w: %s", ErrInvalidAuthorityConfig, structName)
			}
			dir := s.gen.baseOutput
			dir = filepath.Join(append([]string{dir}, config.BaseRouterWithAuthority.RelativeModelPath...)...)
			pkg := config.BaseRouterWithAuthority.RelativeModelPath[len(config.BaseRouterWithAuthority.RelativeModelPath)-1]
			filePath := filepath.Join(dir, config.BaseRouterWithAuthority.FilePrefix, str.CamelToSnake(str.LowFirstChar(structName))+"_router.go")
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return err
			}
			// 判断文件是否存在
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println(structName, "已有权限路由文件 略过生成")
			} else {
				fmt.Println("生成权限路由文件", structName, filePath)
				var buf bytes.Buffer
				if err = s.gen.render(routerImpl, &buf, RouterData{
					ModelStructName:      structName,
					IDType:               s.gen.modelIDTypes[structName],
					WithAuthority:        true,
					ParamName:            str.LowFirstChar(structName),
					PkgName:              pkg,
					ModelPKG:             s.gen.modelPkg,
					BizPKG:               bizPKG,
					GroupPath:            config.BaseRouterWithAuthority.GroupPath,
					AuthorityStructField: config.BaseRouterWithAuthority.AuthorityStructField,
					AuthorityColumn:      config.BaseRouterWithAuthority.AuthorityColumn,
					AuthorityFetchCode:   config.BaseRouterWithAuthority.AuthorityFetchCode,
					DisableBaseHandler:   config.BaseRouterWithAuthority.DisableBaseHandler,
				}); err != nil {
					return err
				}
				result, err := imports.Process(filePath, buf.Bytes(), nil)
				if err != nil {
					return err
				}
				if err = os.WriteFile(filePath, result, 0644); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
