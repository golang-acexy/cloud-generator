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

// RouterGen generates missing router files; existing files are skipped without content validation.
type RouterGen struct {
	gen    *Generator
	config map[string]*RouterConfig
}

// RouterData holds router template inputs, including trusted injected Go code.
type RouterData struct {
	// ModelStructName is the generated Go model name.
	ModelStructName string
	// IDType is the model ID type expression; accepted metadata does not guarantee all templates compile.
	IDType string
	// WithAuthority selects authority-aware router construction.
	WithAuthority bool
	// ParamName is the lower-initial name used for generated package variables.
	ParamName string
	// PkgName is the generated router package name.
	PkgName string
	// ModelPKG is the model import path.
	ModelPKG string
	// BizPKG is the business service import path.
	BizPKG string
	// AuthorityFetchCode is inserted verbatim as a Go expression; it is not validated for type compatibility.
	AuthorityFetchCode string
	// AuthorityStructField names the writable save DTO identity field; its existence and type are not checked.
	AuthorityStructField string
	// AuthorityColumn names the persistence ownership column; it is not checked against table metadata.
	AuthorityColumn string
	// GroupPath is inserted into the generated router group string; use a valid trusted Go string value.
	GroupPath string
	// DisableBaseHandler omits default CRUD registration, leaving the router available for custom handlers.
	DisableBaseHandler bool
}

// NewRouterGen borrows a generator and router map; entries must be non-nil and ID metadata must already exist.
func NewRouterGen(gen *Generator, config map[string]*RouterConfig) *RouterGen {
	return &RouterGen{
		gen:    gen,
		config: config,
	}
}

// Create writes missing router files after shallow path/authority checks; partial output is retained on failure.
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
