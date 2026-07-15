package test

import (
	_ "embed"
	"testing"

	"github.com/golang-acexy/cloud-generator/generatorcloud"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed schema.sql
var schemaSQL string

func TestGenerateCloudSimpleDemo(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:root@(127.0.0.1:13306)/test?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"))
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Exec(schemaSQL).Error; err != nil {
		t.Fatal(err)
	}

	generator := generatorcloud.NewGen(db, "/Users/acexy/Repository/github/golang-acexy/cloud-simple-demo/internal/model", []generatorcloud.TableConfig{
		{
			TableName: "demo_department",
			ModelName: "Department",
			Router: &generatorcloud.RouterConfig{
				BaseRouter: &generatorcloud.BaseRouter{
					RelativeModelPath: []string{"..", "handler", "rest", "adm"},
					GroupPath:         "adm/department",
				},
			},
		},
		{
			TableName: "demo_employee",
			ModelName: "Employee",
			Router: &generatorcloud.RouterConfig{
				BaseRouter: &generatorcloud.BaseRouter{
					RelativeModelPath: []string{"..", "handler", "rest", "adm"},
					GroupPath:         "adm/employee",
				},
				BaseRouterWithAuthority: &generatorcloud.BaseRouterWithAuthority{
					BaseRouter: generatorcloud.BaseRouter{
						RelativeModelPath: []string{"..", "handler", "rest", "usr"},
						GroupPath:         "usr/employee",
					},
					AuthorityStructField: "UserID",
					AuthorityColumn:      "user_id",
					AuthorityFetchCode:   "biz.UsrAuthorityFetch",
				},
			},
		},
	})
	generator.SetIncludeModelPkgPath("github.com/golang-acexy/cloud-simple-demo/internal/model")
	generator.SetRepoRelativeModelPath([]string{"..", "service", "repo"})
	generator.SetServiceRelativeModelPath([]string{"..", "service", "biz"})
	generator.SetServiceBase(&generatorcloud.ServiceBase{
		DefaultOrderBy: "id desc",
		MaxQuerySize:   100,
	})

	if err = generator.Create(); err != nil {
		t.Fatal(err)
	}
}
