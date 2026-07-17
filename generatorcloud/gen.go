package generatorcloud

import (
	"fmt"
	"io"
	"text/template"

	"github.com/acexy/gen"
	"github.com/acexy/golang-toolkit/util/coll"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TableConfig struct {
	TableName string
	ModelName string
	// model 单独设置
	DTOExcluded ModelDTOExcluded

	// service 单独设置
	DisableService bool

	Router *RouterConfig
}

type ModelBase struct {
	DTOExcluded ModelDTOExcluded
}

type ServiceBase struct {
	DefaultOrderBy string
	MaxQuerySize   int
}

type ModelDTOExcluded struct {
	SaveDTOExcludedFields   []string
	QueryDTOExcludedFields  []string
	ModifyDTOExcludedFields []string
	DTOExcludedFields       []string
}

type BaseRouter struct {
	// model 相对于 modelBase 的路径
	RelativeModelPath []string
	GroupPath         string
	FilePrefix        string
}

type BaseRouterWithAuthority struct {
	BaseRouter
	AuthorityFetchCode   string
	AuthorityStructField string
	AuthorityColumn      string
	DisableBaseHandler   bool
}

type Generator struct {
	gen        *gen.Generator
	db         *gorm.DB
	modelPkg   string
	baseOutput string

	modelBase   *ModelBase
	serviceBase *ServiceBase

	repoRelativeModelPath    []string
	serviceRelativeModelPath []string

	tableConfigs []TableConfig
	// key为modelName
	tableConfigsMap map[string]TableConfig
	// key 为模型名，值为模型主键的 Go 类型。
	modelIDTypes map[string]string
}

type RouterConfig struct {
	// 不带权限控制的基础路由
	BaseRouter *BaseRouter
	// 带权限控制的基础路由
	BaseRouterWithAuthority *BaseRouterWithAuthority
}

func NewGen(db *gorm.DB, baseRootPath string, tableConfigs []TableConfig) *Generator {
	d := &Generator{
		tableConfigs: tableConfigs,
		tableConfigsMap: coll.SliceFilterToMap(tableConfigs, func(tableConfig TableConfig) (string, TableConfig, bool) {
			return tableConfig.ModelName, tableConfig, true
		}),
		baseOutput:   baseRootPath,
		db:           db,
		modelBase:    &ModelBase{},
		modelIDTypes: make(map[string]string),
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: baseRootPath,
		Mode:    gen.WithoutContext,
	})
	g.UseDB(db)
	d.gen = g
	return d
}

func (g *Generator) getTableConfig(modelName string) TableConfig {
	return g.tableConfigsMap[modelName]
}

// SetModelBase 设置model的基础信息
func (g *Generator) SetModelBase(base *ModelBase) {
	if base == nil {
		g.modelBase = &ModelBase{}
		return
	}
	g.modelBase = base
}

// SetIncludeModelPkgPath 设置model包名 用于使用该package
func (g *Generator) SetIncludeModelPkgPath(modelPkg string) {
	g.modelPkg = modelPkg
}

// SetRepoRelativeModelPath 设置 repo 相对于 modelBase 的路径 例如 上一级 []string{".."}
func (g *Generator) SetRepoRelativeModelPath(path []string) {
	g.repoRelativeModelPath = path
}

// SetServiceRelativeModelPath 添加 service 相对于 modelBase 的路径 例如 上一级 []string{".."}
func (g *Generator) SetServiceRelativeModelPath(path []string) {
	g.serviceRelativeModelPath = path
}

// SetServiceBase 设置 service 基础信息
func (g *Generator) SetServiceBase(base *ServiceBase) {
	g.serviceBase = base
}

func (g *Generator) rawGen() *gen.Generator {
	return g.gen
}

func (g *Generator) dBType() string {
	if g.db == nil {
		return "unknown"
	}
	switch g.db.Dialector.(type) {
	case *mysql.Dialector:
		return "mysql"
	case *postgres.Dialector:
		return "postgres"
	default:
		return "unknown"
	}
}

func (g *Generator) render(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}

type DBTypeData struct {
	StructName string
	DBType     string
}

func (g *Generator) Create() error {
	if g.modelPkg == "" {
		return ErrModelPackageRequired
	}
	if g.dBType() == "unknown" {
		var dialector any
		if g.db != nil {
			dialector = g.db.Dialector
		}
		return fmt.Errorf("%w: %T", ErrUnsupportedDatabase, dialector)
	}
	return NewModelGen(g).Create()
}
