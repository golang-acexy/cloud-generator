package model

import (
	"github.com/acexy/golang-toolkit/util/json"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"github.com/jinzhu/copier"
)

const TableNameTeacher = "demo_teacher"

// Teacher 教师表
type Teacher struct {
	ID         int64                 `gorm:"primaryKey;<-:false" json:"id"`
	CreateTime gormstarter.Timestamp `gorm:"<-:false" json:"createTime"` // 创建时间
	UpdateTime gormstarter.Timestamp `gorm:"<-:false" json:"updateTime"` // 更新时间
	Name       string                `json:"name"`                       // 姓名
	Sex        string                `json:"sex"`                        // 性别
	Age        int32                 `json:"age"`                        // 年龄
	ClassNo    string                `json:"classNo"`                    // 班级编号
}

func (Teacher) TableName() string {
	return TableNameTeacher
}
func (Teacher) DBType() gormstarter.DBType {
	return gormstarter.DBTypeMySQL
}

// TeacherSDTO 定义保存操作允许接收的字段。
type TeacherSDTO struct {
	Name    string `json:"name"`
	Sex     string `json:"sex"`
	Age     int32  `json:"age"`
	ClassNo string `json:"classNo"`
}

// TeacherMDTO 定义修改操作允许接收的字段。
type TeacherMDTO struct {
	Name    string `json:"name"`
	Sex     string `json:"sex"`
	Age     int32  `json:"age"`
	ClassNo string `json:"classNo"`
}

// TeacherQDTO 定义查询操作允许接收的字段。
type TeacherQDTO struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Sex     string `json:"sex"`
	Age     int32  `json:"age"`
	ClassNo string `json:"classNo"`
}

// TeacherDTO 定义允许返回的字段。
type TeacherDTO struct {
	ID         int64          `json:"id"`
	CreateTime json.Timestamp `json:"createTime"`
	UpdateTime json.Timestamp `json:"updateTime"`
	Name       string         `json:"name"`
	Sex        string         `json:"sex"`
	Age        int32          `json:"age"`
	ClassNo    string         `json:"classNo"`
}

func (v Teacher) ToDTO() (*TeacherDTO, error) {
	var dto TeacherDTO
	if err := copier.Copy(&dto, v); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (v Teacher) ParseDTO(dto *TeacherDTO) error {
	return copier.Copy(dto, v)
}

type TeacherSlice []*Teacher

func (v TeacherSlice) ToDTOs() ([]*TeacherDTO, error) {
	dtos := make([]*TeacherDTO, 0, len(v))
	if err := copier.Copy(&dtos, v); err != nil {
		return nil, err
	}
	return dtos, nil
}

func (v TeacherSlice) ParseDTOs(dtos *[]*TeacherDTO) error {
	return copier.Copy(dtos, v)
}

func (v TeacherSDTO) ToT() (*Teacher, error) {
	var entity Teacher
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v TeacherSDTO) ParseT(entity *Teacher) error {
	return copier.Copy(entity, v)
}

func (v TeacherMDTO) ToT() (*Teacher, error) {
	var entity Teacher
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v TeacherMDTO) ParseT(entity *Teacher) error {
	return copier.Copy(entity, v)
}

func (v TeacherQDTO) ToT() (*Teacher, error) {
	var entity Teacher
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v TeacherQDTO) ParseT(entity *Teacher) error {
	return copier.Copy(entity, v)
}

func (v TeacherDTO) ToT() (*Teacher, error) {
	var entity Teacher
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v TeacherDTO) ParseT(entity *Teacher) error {
	return copier.Copy(entity, v)
}
