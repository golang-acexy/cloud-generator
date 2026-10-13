package generatorcloud

import (
	"fmt"
	"io"
	"text/template"

	"github.com/acexy/gen"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/acexy/golang-toolkit/util/str"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TableConfig selects a table and generated Go model; names must be unique and valid for the target package.
type TableConfig struct {
	// TableName is the existing database table to inspect.
	TableName string
	// ModelName is the generated exported Go model name and configuration lookup key.
	ModelName string

	// DTOExcluded adds exact Go field names to the applicable DTO exclusion sets.
	DTOExcluded ModelDTOExcluded

	// DisableService skips service generation but does not remove router references to that service.
	DisableService bool

	// Router optionally requests router generation for this model.
	Router *RouterConfig
}

// ModelBase configures shared DTO exclusions and time-range names; Create normalizes its fields in place.
type ModelBase struct {
	// DTOExcluded adds exact Go field names to the applicable DTO exclusion sets.
	DTOExcluded ModelDTOExcluded
	// DefaultTimeRangeField is normalized to snake case and must appear in AllowedTimeRangeFields.
	DefaultTimeRangeField string
	// AllowedTimeRangeFields lists allowed time columns; existence in every model is not checked.
	AllowedTimeRangeFields []string
}

// ServiceBase sets defaults for newly generated services; existing service files are retained.
type ServiceBase struct {
	// DefaultOrderBy supplies trusted SQL ordering; empty uses id desc.
	DefaultOrderBy string
	// MaxQuerySize caps ordinary generated list queries; nonpositive values use 500 and do not cap pages.
	MaxQuerySize int
}

// ModelDTOExcluded lists exact Go field names excluded from each DTO; global and table exclusions are combined.
type ModelDTOExcluded struct {
	// SaveDTOExcludedFields excludes save fields in addition to built-in ID and audit fields.
	SaveDTOExcludedFields []string
	// QueryDTOExcludedFields excludes fields from typed query DTOs.
	QueryDTOExcludedFields []string
	// ModifyDTOExcludedFields excludes update fields in addition to built-in ID and audit fields.
	ModifyDTOExcludedFields []string
	// DTOExcludedFields excludes response DTO fields only.
	DTOExcludedFields []string
}

// BaseRouter configures a plain router output path and group; path inputs are trusted and not confined to the output root.
type BaseRouter struct {

	// RelativeModelPath gives path components relative to the model output; its final component is the package name.
	RelativeModelPath []string
	// GroupPath is inserted into the generated router group string; use a valid trusted Go string value.
	GroupPath string
	// FilePrefix adds a directory component before the filename; it is not a filename prefix.
	FilePrefix string
}

// BaseRouterWithAuthority configures identity-aware routes; supplied code and column names must be trusted and type-compatible.
type BaseRouterWithAuthority struct {
	// BaseRouter supplies output path and group settings.
	BaseRouter
	// AuthorityFetchCode is inserted verbatim as a Go expression; it is not validated for type compatibility.
	AuthorityFetchCode string
	// AuthorityStructField names the writable save DTO identity field; its existence and type are not checked.
	AuthorityStructField string
	// AuthorityColumn names the persistence ownership column; it is not checked against table metadata.
	AuthorityColumn string
	// DisableBaseHandler omits default CRUD registration, leaving the router available for custom handlers.
	DisableBaseHandler bool
}

// Generator owns mutable generation state; do not configure or run it concurrently.
// It borrows the DB and configuration, writes files incrementally, and does not roll back partial output.
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

	tableConfigsMap map[string]TableConfig

	modelIDTypes map[string]string
}

// RouterConfig requests plain and/or authority routers; colliding output paths keep the first existing file.
type RouterConfig struct {

	// BaseRouter contains the plain router settings.
	BaseRouter *BaseRouter

	// BaseRouterWithAuthority optionally requests an authority-aware router.
	BaseRouterWithAuthority *BaseRouterWithAuthority
}

// NewGen creates a generator borrowing the DB and table configuration; it neither opens nor closes the DB.
// Configure the model import path and time ranges before Create; duplicate model names overwrite lookup entries.
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

// SetModelBase borrows the settings pointer; nil resets it to empty settings that fail Create time-range validation.
func (g *Generator) SetModelBase(base *ModelBase) {
	if base == nil {
		g.modelBase = &ModelBase{}
		return
	}
	g.modelBase = base
}

// SetIncludeModelPkgPath sets the target model import path; nonempty values are not checked against the filesystem.
func (g *Generator) SetIncludeModelPkgPath(modelPkg string) {
	g.modelPkg = modelPkg
}

// SetRepoRelativeModelPath borrows output path components; empty selects repo and the final component supplies the package name.
func (g *Generator) SetRepoRelativeModelPath(path []string) {
	g.repoRelativeModelPath = path
}

// SetServiceRelativeModelPath borrows output path components; empty selects biz and the final component supplies the package name.
func (g *Generator) SetServiceRelativeModelPath(path []string) {
	g.serviceRelativeModelPath = path
}

// SetServiceBase borrows service defaults; nil uses id desc and a list limit of 500.
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

func (g *Generator) render(tmpl string, wr io.Writer, data any) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}

// DBTypeData holds model names and database expressions for template rendering.
type DBTypeData struct {
	// StructName is the generated Go model name.
	StructName string
	// DBType is the Go expression selecting the GORM database dialect.
	DBType string
}

// Create validates required configuration and generates output incrementally.
// Models and repositories may be overwritten; service/router files are retained. Errors or dependency panics can leave partial output.
// Time fields are required even when services are disabled; this does not compile the generated application.
func (g *Generator) Create() error {
	if g.modelPkg == "" {
		return ErrModelPackageRequired
	}
	if err := g.validateTimeRangeConfig(); err != nil {
		return err
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

func (g *Generator) validateTimeRangeConfig() error {
	defaultField := str.CamelToSnake(g.modelBase.DefaultTimeRangeField)
	if defaultField == "" || len(g.modelBase.AllowedTimeRangeFields) == 0 {
		return ErrTimeRangeConfigRequired
	}
	allowedFields := coll.SliceCollect(g.modelBase.AllowedTimeRangeFields, func(field string) string {
		return str.CamelToSnake(field)
	})
	if !coll.SliceContains(allowedFields, defaultField) {
		return ErrDefaultTimeRangeFieldNotAllowed
	}
	g.modelBase.DefaultTimeRangeField = defaultField
	g.modelBase.AllowedTimeRangeFields = allowedFields
	return nil
}
