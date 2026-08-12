package generatorcloud

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acexy/gen/core/generate"
	"github.com/acexy/gen/core/model"
)

func TestTimeFieldTypeMapping(t *testing.T) {
	fieldType := fieldTypeBySourceType("time.Time", "gormstarter.Timestamp")
	// event_time 不是创建、更新时间字段，验证规则不依赖字段名称。
	modelField := fieldType(&model.Field{
		Name:       "EventTime",
		ColumnName: "event_time",
		Type:       "time.Time",
	})
	if modelField.Type != "gormstarter.Timestamp" {
		t.Fatalf("模型时间字段类型 = %q，期望为 %q", modelField.Type, "gormstarter.Timestamp")
	}
	if dtoType := changeType(modelField.Type); dtoType != "json.Timestamp" {
		t.Fatalf("DTO 时间字段类型 = %q，期望为 %q", dtoType, "json.Timestamp")
	}
}

func TestBuildColumnsDataPreservesModelTypes(t *testing.T) {
	data := buildColumnsData("Teacher", &generate.QueryStructMeta{Fields: []*model.Field{
		{Name: "ID", Type: "int64"},
		{Name: "CreatedAt", Type: "gormstarter.Timestamp"},
	}})
	if data.StructName != "Teacher" || data.ParamName != "teacher" || len(data.Fields) != 2 {
		t.Fatalf("Columns 元数据异常: %+v", data)
	}
	if data.Fields[1].Type != "gormstarter.Timestamp" {
		t.Fatalf("Columns 字段类型被错误转换: %+v", data.Fields[1])
	}
}

func TestInsertColumnsImmediatelyAfterModel(t *testing.T) {
	outputFile := filepath.Join(t.TempDir(), "teacher.go")
	source := "package model\n\ntype Teacher struct {\n\tID int64\n\tName string\n}\n\nfunc (Teacher) TableName() string { return \"teacher\" }\n"
	if err := os.WriteFile(outputFile, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	modelGen := &ModelGen{gen: &Generator{}}
	err := modelGen.insertColumns(outputFile, ColumnsData{
		StructName: "Teacher",
		ParamName:  "teacher",
		Fields: []ColumnFieldData{
			{Name: "ID", Type: "int64"},
			{Name: "Name", Type: "string"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	generated := string(content)
	if strings.Contains(generated, "DO NOT EDIT") {
		t.Fatalf("Model 文件不应包含禁止修改标记:\n%s", generated)
	}
	columnsIndex := strings.Index(generated, "type TeacherColumns struct")
	tableNameIndex := strings.Index(generated, "func (Teacher) TableName")
	if columnsIndex < 0 || tableNameIndex < 0 || columnsIndex > tableNameIndex {
		t.Fatalf("Columns 未紧跟 Model 结构体生成:\n%s", generated)
	}
	if !strings.Contains(generated, "func TeacherColumnSet() *TeacherColumns") ||
		!strings.Contains(generated, "return &entity.Name") {
		t.Fatalf("Columns 访问函数或字段 selector 缺失:\n%s", generated)
	}
}

func TestFieldTypeBySourceTypeDoesNotChangeOtherTypes(t *testing.T) {
	fieldType := fieldTypeBySourceType("time.Time", "gormstarter.Timestamp")
	modelField := fieldType(&model.Field{Type: "string"})
	if modelField.Type != "string" {
		t.Fatalf("非时间字段类型被错误修改为 %q", modelField.Type)
	}
}
