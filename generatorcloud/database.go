package generatorcloud

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Database struct {
}

func (Database) Create() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./",
		Mode:    gen.WithoutContext,
	})
	gormdb, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:13306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(gormdb)
	g.FieldWithTypeTag = false
	g.GenerateModel("demo_teacher", gen.FieldTypeReg("^(create_time|update_time)$", "gormstarter.Timestamp"))
	g.Execute()
}
