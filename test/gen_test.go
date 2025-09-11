package test

import (
	"testing"

	"github.com/golang-acexy/cloud-generator/generatorcloud"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGen(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:13306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	g := generatorcloud.NewGen(db, "./model", []generatorcloud.TableConfig{
		{
			TableName: "demo_teacher",
			ModelName: "Teacher",
		},
		{
			TableName: "demo_student",
			ModelName: "Student",
		},
	})
	g.SetIncludeModelPkgPath("github.com/golang-acexy/cloud-generator/test/model")
	g.SetModelBase(&generatorcloud.ModelBase{
		DTOExcluded: generatorcloud.ModelDTOExcluded{
			SaveDTOExcludedFields: []string{
				"ID",
				"CreateTime",
				"UpdateTime",
			},
			QueryDTOExcludedFields: []string{
				"CreateTime",
				"UpdateTime",
			},
			ModifyDTOExcludedFields: []string{
				"ID",
				"CreateTime",
				"UpdateTime",
			},
		},
	})
	g.Create()
}
