package model

import (
	"fmt"
	"github.com/acexy/golang-toolkit/util/json"
	"github.com/golang-acexy/starter-gorm/gormstarter"
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	teacher := Teacher{
		ID:         123,
		CreateTime: gormstarter.Timestamp{Time: time.Now()},
	}
	d := teacher.ToDTO()

	fmt.Println(json.ToJson(d))

	teachers := []*Teacher{
		{
			ID:         123,
			CreateTime: gormstarter.Timestamp{Time: time.Now()},
		},
		{
			ID:         1321,
			CreateTime: gormstarter.Timestamp{Time: time.Now()},
		},
	}

	ds := TeacherSlice(teachers).ToDTOs()
	fmt.Println(json.ToJson(ds))
}
