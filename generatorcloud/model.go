package generatorcloud

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/acexy/gen"
	"github.com/acexy/gen/core/generate"
	"github.com/acexy/gen/core/model"
	"github.com/acexy/gen/field"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/acexy/golang-toolkit/util/str"
	"golang.org/x/tools/imports"
	"os"
	"path/filepath"
	"regexp"
)

var defaultFieldOptions = []gen.ModelOpt{
	gen.FieldTypeReg("^(create_time|update_time|created_at|modified_at|update_at|modified_time)$", "gormstarter.Timestamp"),
	gen.FieldGORMTag("ID", func(tag field.GormTag) field.GormTag {
		tag.Append("primary_key", "<-:false")
		return tag
	}),
}

//go:embed tmpl/file/method.gohtml
var methodTmpl string

//go:embed tmpl/file/dto.gohtml
var dtoTmpl string

//go:embed tmpl/file/repo.gohtml
var repoTmpl string

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

type ModelGen struct {
	gen *Generator
}

func NewModelGen(gen *Generator) *ModelGen {
	return &ModelGen{
		gen: gen,
	}
}

func (m *ModelGen) getDBType() string {
	switch m.gen.dBType() {
	case "mysql":
		return "gormstarter.DBTypeMySQL"
	case "postgres":
		return "gormstarter.DBTypePostgres"
	default:
		return "Unknown"
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
	})
	coll.SliceForeachAll(m.gen.tableConfigs, func(t TableConfig) {
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
	matches := regexp.MustCompile(`^gormstarter\.BaseModel\[(.+)]$`).FindStringSubmatch(typeStr)
	if len(matches) == 2 {
		return matches[1]
	} else if typeStr == "gormstarter.Timestamp" {
		return "json.Timestamp"
	}
	return typeStr
}

func (m *ModelGen) modelAppend(outputFile string, modelName string, meta *generate.QueryStructMeta) {
	file, _ := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// 移除所有gorm tag
	coll.SliceForeachAll(meta.Fields, func(field *model.Field) {
		field.GORMTag = nil
		field.Tag = coll.MapFilterCollect(field.Tag, func(k string, v string) (string, string, bool) {
			if k != "gorm" {
				return k, v, true
			}
			return "", "", false
		})
	})

	// 追加写入接口实现函数 tableName
	_ = m.gen.render(methodTmpl, file, ModelData{
		StructName: modelName,
		DBType:     m.getDBType(),
	})

	coll.SliceForeachAll(meta.Fields, func(field *model.Field) {
		field.Type = changeType(field.Type)
	})

	data := DtoData{
		QueryStructMeta: meta,
	}
	config := m.gen.tableConfigsMap[modelName]
	sExcludedFields := append(m.gen.modelBase.DTOExcluded.SaveDTOExcludedFields, config.DTOExcluded.SaveDTOExcludedFields...)
	if len(sExcludedFields) > 0 {
		data.SExcludedFields = coll.SliceFilterToMap(sExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	qExcludedFields := append(m.gen.modelBase.DTOExcluded.QueryDTOExcludedFields, config.DTOExcluded.QueryDTOExcludedFields...)
	if len(qExcludedFields) > 0 {
		data.QExcludedFields = coll.SliceFilterToMap(qExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	mExcludedFields := append(m.gen.modelBase.DTOExcluded.ModifyDTOExcludedFields, config.DTOExcluded.ModifyDTOExcludedFields...)
	if len(mExcludedFields) > 0 {
		data.MExcludedFields = coll.SliceFilterToMap(mExcludedFields, func(field string) (string, struct{}, bool) {
			return field, struct{}{}, true
		})
	}
	dExcludedFields := append(m.gen.modelBase.DTOExcluded.DTOExcludedFields, config.DTOExcluded.DTOExcludedFields...)
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
	_ = m.gen.render(dtoTmpl, file, data)

	// 修改import
	content, _ := os.ReadFile(outputFile)
	result, _ := imports.Process(outputFile, content, nil)
	_ = os.WriteFile(outputFile, result, 0644)
}

func (m *ModelGen) createRepo(outputFile string, structName string) {
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
	_ = os.MkdirAll(dir, os.ModePerm)
	repoPath = filepath.Join(dir, str.CamelToSnake(str.LowFirstChar(structName))+"_repo.go")

	// 判断文件是否存在
	if _, err := os.Stat(repoPath); err == nil {
		fmt.Println(structName, "已有repo文件 略过生成")
		return
	} else {
		fmt.Println("生成repo文件", structName, repoPath)
	}
	var buf bytes.Buffer
	_ = m.gen.render(repoTmpl, &buf, RepoData{
		StructName: structName,
		ParamName:  str.LowFirstChar(structName),
		Pkg:        m.gen.modelPkg,
		PKG:        pkg,
	})
	result, _ := imports.Process(repoPath, buf.Bytes(), nil)
	_ = os.WriteFile(repoPath, result, os.ModePerm)
}

func (m *ModelGen) Create() {
	m.loadSettings()
	queryGenResult := m.gen.rawGen().ExecuteWithOutInfo()
	var services []string
	routers := make(map[string]*RouterConfig)

	coll.MapForeachAll(queryGenResult.Path, func(modelName string, outputFile string) {
		m.modelAppend(outputFile, modelName, queryGenResult.Meta[modelName])
		m.createRepo(outputFile, modelName)
		config := m.gen.getTableConfig(modelName)
		if !config.DisableService {
			services = append(services, modelName)
		}
		if config.Router != nil {
			routers[modelName] = config.Router
		}
	})
	if len(services) > 0 {
		NewServiceGen(m.gen, services).Create()
	}
	if len(routers) > 0 {
		NewRouterGen(m.gen, routers).Create()
	}

}
