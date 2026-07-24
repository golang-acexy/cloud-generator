package generatorcloud

import (
	"testing"

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

func TestFieldTypeBySourceTypeDoesNotChangeOtherTypes(t *testing.T) {
	fieldType := fieldTypeBySourceType("time.Time", "gormstarter.Timestamp")
	modelField := fieldType(&model.Field{Type: "string"})
	if modelField.Type != "string" {
		t.Fatalf("非时间字段类型被错误修改为 %q", modelField.Type)
	}
}
