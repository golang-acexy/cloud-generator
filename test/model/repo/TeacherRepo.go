package repo

import (
	"github.com/golang-acexy/cloud-database/databasecloud"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"github.com/golang-acexy/cloud-generator/test/model"
)

var teacherRepo = TeacherRepo{
	GormRepository: databasecloud.GormRepository[
		gormstarter.IBaseMapper[gormstarter.BaseMapper[model.Teacher], model.Teacher],
		gormstarter.BaseMapper[model.Teacher],
		model.Teacher,
	]{
		Mapper: TeacherMapper{},
	},
}

func (r TeacherRepo) RawMapper() TeacherMapper {
	return r.RawIMapper().(TeacherMapper)
}
func NewTeacherRepo() TeacherRepo {
	return teacherRepo
}

type TeacherMapper struct {
	gormstarter.BaseMapper[model.Teacher]
}
type TeacherRepo struct {
	databasecloud.GormRepository[gormstarter.IBaseMapper[gormstarter.BaseMapper[model.Teacher],
		model.Teacher], gormstarter.BaseMapper[model.Teacher], model.Teacher]
}

// 在此处拓展自定义Mapper的业务功能

//func (m TeacherMapper) MyMapperMethod() {
//	
//}

// ---------------------------------------

// 在此处拓展自定义Repo的业务功能

//func (r TeacherRepo) MyRepoMethod()  {
//	
//}