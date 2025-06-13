package generatorcloud

import (
	"github.com/acexy/gen"
	"github.com/acexy/golang-toolkit/util/coll"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"text/template"
)

type TableConfig struct {
	TableName string
	ModelName string
	// 保存DTO时，忽略的结构体字段
	SaveDTOExcludedFields []string
	// 查询DTO时，忽略的结构体字段
	QueryDTOExcludedFields []string
	// 修改DTO时，忽略的结构体字段
	ModifyDTOExcludedFields []string
	// DTO时，忽略的结构体字段
	DTOExcludedFields []string

	RouterConfig *RouterConfig
}

type RouterConfig struct {
}

type Generator struct {
	gen      *gen.Generator
	db       *gorm.DB
	modelPkg string

	// 保存DTO时，忽略的结构体字段
	modelBaseSaveDTOExcludedFields []string
	// 查询DTO时，忽略的结构体字段
	modelBaseQueryDTOExcludedFields []string
	// 修改DTO时，忽略的结构体字段
	modelBaseModifyDTOExcludedFields []string

	repoRelativeModelPath []string
	tableConfigs          []TableConfig
	// key为modelName
	tableConfigsMap map[string]TableConfig
}

func NewGen(db *gorm.DB, outPath string, tableConfigs []TableConfig) *Generator {
	d := &Generator{
		tableConfigs: tableConfigs,
		tableConfigsMap: coll.SliceFilterToMap(tableConfigs, func(tableConfig TableConfig) (string, TableConfig, bool) {
			return tableConfig.ModelName, tableConfig, true
		}),
		db: db,
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithoutContext,
	})
	g.UseDB(db)
	d.gen = g
	return d
}

func NewGenWithConfig(db *gorm.DB, tableConfigs []TableConfig, config gen.Config) *Generator {
	d := &Generator{
		tableConfigs: tableConfigs,
		tableConfigsMap: coll.SliceFilterToMap(tableConfigs, func(tableConfig TableConfig) (string, TableConfig, bool) {
			return tableConfig.ModelName, tableConfig, true
		}),
		db: db,
	}
	g := gen.NewGenerator(config)
	g.UseDB(db)
	d.gen = g
	return d
}

func (d *Generator) SetDTOExcludedFields(s, q, m []string) {
	d.modelBaseSaveDTOExcludedFields = s
	d.modelBaseQueryDTOExcludedFields = q
	d.modelBaseModifyDTOExcludedFields = m
}

func (d *Generator) SetModelPkg(modelPkg string) {
	d.modelPkg = modelPkg
}

// SetRepoRelativeModelPath 设置 repo 相对于 model 的路径 例如 上一级 []string{".."}
func (d *Generator) SetRepoRelativeModelPath(path []string) {
	d.repoRelativeModelPath = path
}

func (d *Generator) rawGen() *gen.Generator {
	return d.gen
}

func (d *Generator) dBType() string {
	switch d.db.Dialector.(type) {
	case *mysql.Dialector:
		return "mysql"
	case *postgres.Dialector:
		return "postgres"
	default:
		return "unknown"
	}
}

func (d *Generator) render(tmpl string, wr io.Writer, data interface{}) error {
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

func (d *Generator) Create() {
	NewModelGen(d).Create()
}
