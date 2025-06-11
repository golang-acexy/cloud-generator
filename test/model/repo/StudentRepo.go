package repo

import (
	"github.com/golang-acexy/cloud-database/databasecloud"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"github.com/golang-acexy/cloud-generator/test/model"
)

var studentRepo = StudentRepo{
	GormRepository: databasecloud.GormRepository[
		gormstarter.IBaseMapper[gormstarter.BaseMapper[model.Student], model.Student],
		gormstarter.BaseMapper[model.Student],
		model.Student,
	]{
		Mapper: StudentMapper{},
	},
}

func (r StudentRepo) RawMapper() StudentMapper {
	return r.RawIMapper().(StudentMapper)
}
func NewStudentRepo() StudentRepo {
	return studentRepo
}

type StudentMapper struct {
	gormstarter.BaseMapper[model.Student]
}
type StudentRepo struct {
	databasecloud.GormRepository[gormstarter.IBaseMapper[gormstarter.BaseMapper[model.Student],
		model.Student], gormstarter.BaseMapper[model.Student], model.Student]
}

// 在此处拓展自定义Mapper的业务功能

//func (m StudentMapper) MyMapperMethod() {
//	
//}

// ---------------------------------------

// 在此处拓展自定义Repo的业务功能

//func (r StudentRepo) MyRepoMethod()  {
//	
//}