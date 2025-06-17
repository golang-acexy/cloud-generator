package test

import (
	"github.com/golang-acexy/cloud-generator/generatorcloud"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGen(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:13306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	g := generatorcloud.NewGen(db, "/Users/acexy/Repository/github/golang-acexy/cloud-simple-demo/internal/model", []generatorcloud.TableConfig{
		{
			TableName: "demo_teacher",
			ModelName: "Teacher",
			Router: &generatorcloud.RouterConfig{
				BaseRouter: &generatorcloud.BaseRouter{
					RelativeModelPath: []string{"..", "handler", "rest", "adm"},
					GroupPath:         "adm/teacher",
				},
			},
		},
		{
			TableName: "demo_student",
			ModelName: "Student",
			Router: &generatorcloud.RouterConfig{
				BaseRouter: &generatorcloud.BaseRouter{
					RelativeModelPath: []string{"..", "handler", "rest", "adm"},
					GroupPath:         "adm/student",
				},
				BaseRouterWithDataCheck: &generatorcloud.BaseRouterWithDataCheck{
					BaseRouter: generatorcloud.BaseRouter{
						RelativeModelPath: []string{"..", "handler", "rest", "usr"},
						GroupPath:         "usr/student",
					},
					DataLimitStructName: "UserID",
					AuthorityFetchCode:  "biz.UsrAuthorityFetch",
				},
			},
		},
	})
	// model基础设置
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
	g.SetIncludeModelPkgPath("github.com/golang-acexy/cloud-simple-demo/internal/model")

	g.SetRepoRelativeModelPath([]string{
		"..", "service", "repo",
	})
	g.SetServiceRelativeModelPath([]string{
		"..", "service", "biz",
	})

	g.Create()
}
