package generatorcloud

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/acexy/gen"
	"github.com/acexy/gen/core/generate"
	"github.com/acexy/gen/core/model"
	"github.com/acexy/gen/field"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/acexy/golang-toolkit/util/str"
	"golang.org/x/tools/imports"
)

var defaultFieldOptions = []gen.ModelOpt{
	fieldTypeBySourceType("time.Time", "gormstarter.Timestamp"),
	gen.FieldTypeReg("^(create_time|update_time|created_at|updated_at|modified_at|update_at|modified_time)$", "gormstarter.Timestamp"),
	gen.FieldGORMTag("ID", func(tag field.GormTag) field.GormTag {
		tag.Append("primary_key", "<-:false")
		return tag
	}),
}

// fieldTypeBySourceType 按生成器推导出的 Go 类型替换模型字段类型。
// 数据库的 DATE、DATETIME、TIMESTAMP 会由生成器映射为 time.Time，统一改为框架时间类型。
func fieldTypeBySourceType(sourceType, targetType string) model.ModifyFieldOpt {
	return func(field *model.Field) *model.Field {
		if field.Type == sourceType {
			field.Type = targetType
		}
		return field
	}
}

var databaseGeneratedFields = []string{
	"ID",
	"CreatedAt",
	"CreateTime",
	"UpdatedAt",
	"UpdateTime",
	"ModifiedAt",
	"UpdateAt",
	"ModifiedTime",
	"DeletedAt",
}

//go:embed tmpl/file/method.gohtml
var methodTmpl string

//go:embed tmpl/file/dto.gohtml
var dtoTmpl string

//go:embed tmpl/file/repo.gohtml
var repoTmpl string

//go:embed tmpl/file/columns.gohtml
var columnsTmpl string

type ModelData struct {
	StructName string
	DBType     string
}

type RepoData struct {
	StructName string
	ParamName  string
	Pkg        string
	PKG        string
}

type ColumnFieldData struct {
	Name string
	Type string
}

type ColumnsData struct {
	StructName string
	ParamName  string
	Fields     []ColumnFieldData
}

type ModelGen struct {
	gen *Generator
}

func NewModelGen(gen *Generator) *ModelGen {
	return &ModelGen{
		gen: gen,
	}
}

func (m *ModelGen) getDBType() (string, error) {
	switch m.gen.dBType() {
	case "mysql":
		return "gormstarter.DBTypeMySQL", nil
	case "postgres":
		return "gormstarter.DBTypePostgres", nil
	default:
		return "", ErrUnsupportedDatabase
	}
}

func (m *ModelGen) loadSettings() {
	m.gen.rawGen().WithJSONTagNameStrategy(func(c string) string { return str.SnakeToCamel(c) })
	m.gen.rawGen().DisableDefaultGormTag()
	m.gen.rawGen().MustBindGormTag(map[string]map[string][]string{
		"ID": {
			"<-":         {"false"},
			"primaryKey": nil,
		},
		"CreatedAt": {
			"<-": {"false"},
		},
		"UpdatedAt": {
			"<-": {"false"},
		},
		"CreateTime": {
			"<-": {"false"},
		},
		"UpdateTime": {
			"<-": {"false"},
		},
		"ModifiedAt": {
			"<-": {"false"},
		},
		"UpdateAt": {
			"<-": {"false"},
		},
		"ModifiedTime": {
			"<-": {"false"},
		},
		"DeletedAt": {
			"<-": {"false"},
		},
	})
	coll.SliceForEachAll(m.gen.tableConfigs, func(t TableConfig) {
		m.gen.rawGen().GenerateModelAs(t.TableName, t.ModelName, defaultFieldOptions...)
	})
}

type DtoData struct {
	*generate.QueryStructMeta

	IsSExcluded     func(s string) bool
	SExcludedFields map[string]struct{}

	IsQExcluded     func(q string) bool
	QExcludedFields map[string]struct{}

	IsMExcluded     func(m string) bool
	MExcludedFields map[string]struct{}

	IsDExcluded     func(m string) bool
	DExcludedFields map[string]struct{}
}

func changeType(typeStr string) string {
	if typeStr == "gormstarter.Timestamp" {
		return "json.Timestamp"
	}
	return typeStr
}

func modelIDType(meta *generate.QueryStructMeta) (string, error) {
	for _, modelField := range meta.Fields {
		if modelField.Name == "ID" {
			idType := changeType(modelField.Type)
			switch idType {
			case "int", "uint", "int32", "uint32", "int64", "uint64", "string":
				return idType, nil
			default:
				return "", fmt.Errorf("%w: %s", ErrUnsupportedIDType, idType)
			}
		}
	}
	// 保持与旧版生成器一致；没有标准 ID 字段的模型仍可由调用方继续扩展。
	return "int64", nil
}

func buildColumnsData(modelName string, meta *generate.QueryStructMeta) ColumnsData {
	fields := make([]ColumnFieldData, 0, len(meta.Fields))
	for _, modelField := range meta.Fields {
		if modelField.Relation != nil {
			continue
		}
		fields = append(fields, ColumnFieldData{Name: modelField.Name, Type: modelField.Type})
	}
	return ColumnsData{StructName: modelName, ParamName: str.LowFirstChar(modelName), Fields: fields}
}

func (m *ModelGen) insertColumns(outputFile string, data ColumnsData) error {
	content, err := os.ReadFile(outputFile)
	if err != nil {
		return err
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, outputFile, content, parser.ParseComments)
	if err != nil {
		return err
	}
	insertOffset := -1
	for _, declaration := range parsed.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.TYPE {
			continue
		}
		for _, spec := range general.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if ok && typeSpec.Name.Name == data.StructName {
				insertOffset = fileSet.Position(general.End()).Offset
				break
			}
		}
		if insertOffset >= 0 {
			break
		}
	}
	if insertOffset < 0 {
		return fmt.Errorf("model struct not found: %s", data.StructName)
	}
	var columns bytes.Buffer
	if err = m.gen.render(columnsTmpl, &columns, data); err != nil {
		return err
	}
	result := make([]byte, 0, len(content)+columns.Len()+2)
	result = append(result, content[:insertOffset]...)
	result = append(result, '\n', '\n')
	result = append(result, columns.Bytes()...)
	result = append(result, content[insertOffset:]...)
	return os.WriteFile(outputFile, result, 0644)
}

func (m *ModelGen) modelAppend(outputFile string, modelName string, meta *generate.QueryStructMeta) error {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 移除所有gorm tag
	coll.SliceForEachAll(meta.Fields, func(field *model.Field) {
		field.GORMTag = nil
		field.Tag = coll.MapFilterCollect(field.Tag, func(k string, v string) (string, string, bool) {
			if k != "gorm" {
				return k, v, true
			}
			return "", "", false
		})
	})

	// 追加写入接口实现函数 tableName
	dbType, err := m.getDBType()
	if err != nil {
		return err
	}
	if err = m.gen.render(methodTmpl, file, ModelData{
		StructName: modelName,
		DBType:     dbType,
	}); err != nil {
		return err
	}

	coll.SliceForEachAll(meta.Fields, func(field *model.Field) {
		field.Type = changeType(field.Type)
	})

	data := DtoData{
		QueryStructMeta: meta,
	}
	config := m.gen.tableConfigsMap[modelName]
	sExcludedFields := append([]string{}, databaseGeneratedFields...)
	sExcludedFields = append(sExcludedFields, m.gen.modelBase.DTOExcluded.SaveDTOExcludedFields...)
	sExcludedFields = append(sExcludedFields, config.DTOExcluded.SaveDTOExcludedFields...)
	if len(sExcludedFields) > 0 {
		data.SExcludedFields = coll.SliceFilterToMap(sExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	qExcludedFields := append([]string{}, m.gen.modelBase.DTOExcluded.QueryDTOExcludedFields...)
	qExcludedFields = append(qExcludedFields, config.DTOExcluded.QueryDTOExcludedFields...)
	if len(qExcludedFields) > 0 {
		data.QExcludedFields = coll.SliceFilterToMap(qExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	mExcludedFields := append([]string{}, databaseGeneratedFields...)
	mExcludedFields = append(mExcludedFields, m.gen.modelBase.DTOExcluded.ModifyDTOExcludedFields...)
	mExcludedFields = append(mExcludedFields, config.DTOExcluded.ModifyDTOExcludedFields...)
	if len(mExcludedFields) > 0 {
		data.MExcludedFields = coll.SliceFilterToMap(mExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	dExcludedFields := append([]string{}, m.gen.modelBase.DTOExcluded.DTOExcludedFields...)
	dExcludedFields = append(dExcludedFields, config.DTOExcluded.DTOExcludedFields...)
	if len(dExcludedFields) > 0 {
		data.DExcludedFields = coll.SliceFilterToMap(dExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	data.IsQExcluded = func(s string) bool {
		_, ok := data.QExcludedFields[s]
		return ok
	}
	data.IsSExcluded = func(s string) bool {
		_, ok := data.SExcludedFields[s]
		return ok
	}
	data.IsMExcluded = func(s string) bool {
		_, ok := data.MExcludedFields[s]
		return ok
	}
	data.IsDExcluded = func(s string) bool {
		_, ok := data.DExcludedFields[s]
		return ok
	}
	// 追加写入DTO
	if err = m.gen.render(dtoTmpl, file, data); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}

	// 修改import
	content, err := os.ReadFile(outputFile)
	if err != nil {
		return err
	}
	result, err := imports.Process(outputFile, content, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(outputFile, result, 0644)
}

func (m *ModelGen) createRepo(outputFile string, structName string) error {
	dir := filepath.Dir(outputFile)
	var repoPath string
	var pkg string
	if len(m.gen.repoRelativeModelPath) > 0 {
		dir = filepath.Join(append([]string{dir}, m.gen.repoRelativeModelPath...)...)
		pkg = m.gen.repoRelativeModelPath[len(m.gen.repoRelativeModelPath)-1]
	} else {
		dir = filepath.Join(dir, "repo")
		pkg = "repo"
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	repoPath = filepath.Join(dir, str.CamelToSnake(str.LowFirstChar(structName))+"_repo.go")

	fmt.Println("生成repo文件", structName, repoPath)
	var buf bytes.Buffer
	if err := m.gen.render(repoTmpl, &buf, RepoData{
		StructName: structName,
		ParamName:  str.LowFirstChar(structName),
		Pkg:        m.gen.modelPkg,
		PKG:        pkg,
	}); err != nil {
		return err
	}
	result, err := imports.Process(repoPath, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(repoPath, result, 0644)
}

func (m *ModelGen) Create() error {
	m.loadSettings()
	queryGenResult := m.gen.rawGen().ExecuteWithOutInfo()
	var services []string
	routers := make(map[string]*RouterConfig)

	for modelName, outputFile := range queryGenResult.Path {
		config := m.gen.getTableConfig(modelName)
		idType, err := modelIDType(queryGenResult.Meta[modelName])
		if err != nil && (!config.DisableService || config.Router != nil) {
			return fmt.Errorf("model %s: %w", modelName, err)
		}
		m.gen.modelIDTypes[modelName] = idType
		columnsData := buildColumnsData(modelName, queryGenResult.Meta[modelName])
		if err := m.insertColumns(outputFile, columnsData); err != nil {
			return err
		}
		if err := m.modelAppend(outputFile, modelName, queryGenResult.Meta[modelName]); err != nil {
			return err
		}
		if err := m.createRepo(outputFile, modelName); err != nil {
			return err
		}
		if !config.DisableService {
			services = append(services, modelName)
		}
		if config.Router != nil {
			routers[modelName] = config.Router
		}
	}
	if len(services) > 0 {
		if err := NewServiceGen(m.gen, services).Create(); err != nil {
			return err
		}
	}
	if len(routers) > 0 {
		if err := NewRouterGen(m.gen, routers).Create(); err != nil {
			return err
		}
	}
	return nil
}
