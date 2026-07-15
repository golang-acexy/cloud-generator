package repo

import (
	"github.com/golang-acexy/cloud-database/databasecloud/rds"
	"github.com/golang-acexy/cloud-generator/test/model"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"gorm.io/gorm"
)

type StudentMapper struct {
	gormstarter.BaseMapper[model.Student]
}

func (m StudentMapper) WithTxMapper(tx *gorm.DB) StudentMapper {
	return StudentMapper{
		BaseMapper: m.GetBaseMapperWithTx(tx),
	}
}

type StudentRepo struct {
	rds.Repository[StudentRepo, StudentMapper, model.Student]
}

var studentRepo = newStudentRepo()

func newStudentRepo() StudentRepo {
	repository := rds.NewRepository(
		StudentMapper{},
		func(repository rds.Repository[StudentRepo, StudentMapper, model.Student]) StudentRepo {
			return StudentRepo{Repository: repository}
		},
	)
	return StudentRepo{Repository: repository}
}

func NewStudentRepo() StudentRepo {
	return studentRepo
}

// 在此处扩展自定义 Mapper 能力。

// 在此处扩展自定义 Repository 能力。
