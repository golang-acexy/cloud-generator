package test

import (
	"github.com/golang-acexy/cloud-generator/generatorcloud/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestDatabase(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:13306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	g := database.NewDatabaseGen(db, "./model", map[string]string{
		"demo_teacher": "Teacher",
		"demo_student": "Student",
	})
	g.SetModelPkg("github.com/golang-acexy/cloud-generator/test/model")
	g.Create()
}
