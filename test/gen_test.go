package test

import (
	"github.com/golang-acexy/cloud-generator/generatorcloud"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestDatabase(t *testing.T) {
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
	g.SetDTOExcludedFields([]string{
		"ID",
		"CreateTime",
		"UpdateTime",
	}, []string{
		"CreateTime",
		"UpdateTime",
	}, []string{
		"ID",
		"CreateTime",
		"UpdateTime",
	})

	g.SetModelPkg("github.com/golang-acexy/cloud-generator/test/model")

	g.Create()
}
