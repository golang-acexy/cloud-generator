package database

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/acexy/gen"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/acexy/golang-toolkit/util/str"
	"os"
	"path/filepath"
)

var defaultFieldOptions = []gen.ModelOpt{
	gen.FieldTypeReg("^(create_time|update_time)$", "gormstarter.Timestamp"),
	gen.FieldTypeReg("^id$", "gormstarter.BaseModel[uint64]"),
}

//go:embed tmpl/file/model.gohtml
var modelTmpl string

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
}

type ModelGen struct {
	gen *DatabaseGen
}

func NewModelGen(gen *DatabaseGen) *ModelGen {
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
	if !m.gen.disabledDefaultSettings {
		m.gen.rawGen().WithJSONTagNameStrategy(func(c string) string { return str.SnakeToCamel(c) })
		m.gen.rawGen().DisableGormTag()
		coll.MapForeachAll(m.gen.tableModelName, func(tableName, modelName string) {
			m.gen.rawGen().GenerateModelAs(tableName, modelName, defaultFieldOptions...)
		})
	}
}

func (m *ModelGen) appendDBType(outputFile string, structName string) {
	file, _ := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	_ = m.gen.render(modelTmpl, file, ModelData{
		StructName: structName,
		DBType:     m.getDBType(),
	})
}

func (m *ModelGen) createRepo(outputFile string, structName string) {
	dir := filepath.Dir(outputFile)
	dir = filepath.Join(dir, "repo")
	_ = os.MkdirAll(dir, os.ModePerm)
	repoPath := filepath.Join(dir, structName+"Repo.go")
	// 判断文件是否存在
	if _, err := os.Stat(repoPath); err == nil {
		fmt.Println(structName, "已有repo文件 略过生成")
		return
	}
	var buf bytes.Buffer
	_ = m.gen.render(repoTmpl, &buf, RepoData{
		StructName: structName,
		ParamName:  str.LowFirstChar(structName),
		Pkg:        m.gen.modelPkg,
	})
	_ = os.WriteFile(repoPath, buf.Bytes(), os.ModePerm)
}

func (m *ModelGen) Create() {
	m.loadSettings()
	modelInfo := m.gen.rawGen().ExecuteWithOutInfo()
	coll.MapForeachAll(modelInfo, func(tableName string, outputFile string) {
		m.appendDBType(outputFile, m.gen.tableModelName[tableName])
		m.createRepo(outputFile, m.gen.tableModelName[tableName])
	})

	// 创建repo
}
