package model

import (
	"github.com/acexy/golang-toolkit/logger"
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

// StudentSDTO 保存时结构体 限定只允许操作的字段
type StudentSDTO struct {
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Age       int32  `json:"age"`
	TeacherID int64  `json:"teacherId"`
}

// StudentMDTO 修改时结构体 限定只允许操作的字段
type StudentMDTO struct {
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Age       int32  `json:"age"`
	TeacherID int64  `json:"teacherId"`
}

// StudentQDTO 查询时结构体 限定只允许使用的字段
type StudentQDTO struct {
	ID        int64  `gorm:"primaryKey;<-:false" json:"id"`
	Name      string `json:"name"`
	Sex       string `json:"sex"`
	Age       int32  `json:"age"`
	TeacherID int64  `json:"teacherId"`
}

// StudentDTO 结果结构体 限定只允许展示的字段
type StudentDTO struct {
	ID         int64          `gorm:"primaryKey;<-:false" json:"id"`
	CreateTime json.Timestamp `gorm:"<-:false" json:"createTime"`
	UpdateTime json.Timestamp `gorm:"<-:false" json:"updateTime"`
	Name       string         `json:"name"`
	Sex        string         `json:"sex"`
	Age        int32          `json:"age"`
	TeacherID  int64          `json:"teacherId"`
}

func (v Student) ToDTO() *StudentDTO {
	var dto StudentDTO
	err := copier.Copy(&dto, v)
	if err != nil {
		logger.Logrus().Errorln("copier.Copy error:", err)
		return nil
	}
	return &dto
}

type StudentSlice []*Student

func (v StudentSlice) ToDTOs() []*StudentDTO {
	var dtos []*StudentDTO
	err := copier.Copy(&dtos, v)
	if err != nil {
		logger.Logrus().Errorln("copier.Copy error:", err)
		return nil
	}
	return dtos
}

func (v StudentSDTO) ToT() *Student {
	var t Student
	err := copier.Copy(&t, v)
	if err != nil {
		logger.Logrus().Errorln("copier.Copy error:", err)
		return nil
	}
	return &t
}

func (v StudentMDTO) ToT() *Student {
	var t Student
	err := copier.Copy(&t, v)
	if err != nil {
		logger.Logrus().Errorln("copier.Copy error:", err)
		return nil
	}
	return &t
}

func (v StudentQDTO) ToT() *Student {
	var t Student
	err := copier.Copy(&t, v)
	if err != nil {
		logger.Logrus().Errorln("copier.Copy error:", err)
		return nil
	}
	return &t
}

func (v StudentDTO) ToT() *Student {
	var t Student
	err := copier.Copy(&t, v)
	if err != nil {
		logger.Logrus().Errorln("copier.Copy error:", err)
		return nil
	}
	return &t
}
