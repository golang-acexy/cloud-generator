package repo

import (
	"github.com/golang-acexy/cloud-database/databasecloud/rds"
	"github.com/golang-acexy/cloud-generator/test/model"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"gorm.io/gorm"
)

type TeacherMapper struct {
	gormstarter.BaseMapper[model.Teacher]
}

func (m TeacherMapper) WithTxMapper(tx *gorm.DB) TeacherMapper {
	return TeacherMapper{
		BaseMapper: m.GetBaseMapperWithTx(tx),
	}
}

type TeacherRepo struct {
	rds.Repository[TeacherRepo, TeacherMapper, model.Teacher]
}

var teacherRepo = newTeacherRepo()

func newTeacherRepo() TeacherRepo {
	repository := rds.NewRepository(
		TeacherMapper{},
		func(repository rds.Repository[TeacherRepo, TeacherMapper, model.Teacher]) TeacherRepo {
			return TeacherRepo{Repository: repository}
		},
	)
	return TeacherRepo{Repository: repository}
}

func NewTeacherRepo() TeacherRepo {
	return teacherRepo
}

// 在此处扩展自定义 Mapper 能力。

// 在此处扩展自定义 Repository 能力。
