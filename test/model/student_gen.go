package model

import (
	"github.com/acexy/golang-toolkit/util/json"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"github.com/jinzhu/copier"
)

const TableNameStudent = "demo_student"

// Student mapped from table <demo_student>
type Student struct {
	ID         int64                 `gorm:"primaryKey;<-:false" json:"id"`
	CreateTime gormstarter.Timestamp `gorm:"<-:false" json:"createTime"`
	UpdateTime gormstarter.Timestamp `gorm:"<-:false" json:"updateTime"`
	Name       string                `json:"name"`
	Sex        string                `json:"sex"`
	Age        int32                 `json:"age"`
	TeacherID  int64                 `json:"teacherId"`
}

func (Student) TableName() string {
	return TableNameStudent
}
func (Student) DBType() gormstarter.DBType {
	return gormstarter.DBTypeMySQL
}

// StudentSDTO 定义保存操作允许接收的字段。
type StudentSDTO struct {
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Age       int32  `json:"age"`
	TeacherID int64  `json:"teacherId"`
}

// StudentMDTO 定义修改操作允许接收的字段。
type StudentMDTO struct {
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Age       int32  `json:"age"`
	TeacherID int64  `json:"teacherId"`
}

// StudentQDTO 定义查询操作允许接收的字段。
type StudentQDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Age       int32  `json:"age"`
	TeacherID int64  `json:"teacherId"`
}

// StudentDTO 定义允许返回的字段。
type StudentDTO struct {
	ID         int64          `json:"id"`
	CreateTime json.Timestamp `json:"createTime"`
	UpdateTime json.Timestamp `json:"updateTime"`
	Name       string         `json:"name"`
	Sex        string         `json:"sex"`
	Age        int32          `json:"age"`
	TeacherID  int64          `json:"teacherId"`
}

func (v Student) ToDTO() (*StudentDTO, error) {
	var dto StudentDTO
	if err := copier.Copy(&dto, v); err != nil {
		return nil, err
	}
	return &dto, nil
}

func (v Student) ParseDTO(dto *StudentDTO) error {
	return copier.Copy(dto, v)
}

type StudentSlice []*Student

func (v StudentSlice) ToDTOs() ([]*StudentDTO, error) {
	dtos := make([]*StudentDTO, 0, len(v))
	if err := copier.Copy(&dtos, v); err != nil {
		return nil, err
	}
	return dtos, nil
}

func (v StudentSlice) ParseDTOs(dtos *[]*StudentDTO) error {
	return copier.Copy(dtos, v)
}

func (v StudentSDTO) ToT() (*Student, error) {
	var entity Student
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v StudentSDTO) ParseT(entity *Student) error {
	return copier.Copy(entity, v)
}

func (v StudentMDTO) ToT() (*Student, error) {
	var entity Student
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v StudentMDTO) ParseT(entity *Student) error {
	return copier.Copy(entity, v)
}

func (v StudentQDTO) ToT() (*Student, error) {
	var entity Student
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v StudentQDTO) ParseT(entity *Student) error {
	return copier.Copy(entity, v)
}

func (v StudentDTO) ToT() (*Student, error) {
	var entity Student
	if err := copier.Copy(&entity, v); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (v StudentDTO) ParseT(entity *Student) error {
	return copier.Copy(entity, v)
}
