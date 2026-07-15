package biz

import (
	"github.com/golang-acexy/cloud-database/databasecloud"
	"github.com/golang-acexy/cloud-database/databasecloud/rds"
	"github.com/golang-acexy/cloud-generator/test/model"
	"github.com/golang-acexy/cloud-generator/test/model/repo"
	"github.com/golang-acexy/cloud-web/webcloud"
)

var teacherBizService = &TeacherBizService{
	repo: repo.NewTeacherRepo(),
}

func NewTeacherBizService() *TeacherBizService {
	return teacherBizService
}

type TeacherBizService struct {
	repo repo.TeacherRepo
}

func (v *TeacherBizService) MaxQuerySize() int {
	return 500
}

func (v *TeacherBizService) DefaultOrderBy() string {
	return "id desc"
}

func (v *TeacherBizService) Save(save *model.TeacherSDTO) (int64, error) {
	if save == nil {
		return 0, webcloud.ErrBadRequestParameters
	}
	entity, err := save.ToT()
	if err != nil {
		return 0, err
	}
	if _, err = v.repo.SaveWithoutZeroFields(entity); err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (v *TeacherBizService) BaseQueryByID(condition map[string]any) (*model.TeacherDTO, error) {
	var entity model.Teacher
	row, err := v.repo.QueryOneByMap(condition, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *TeacherBizService) BaseQueryOne(condition map[string]any) (*model.TeacherDTO, error) {
	var entity model.Teacher
	row, err := v.repo.QueryOneByMap(condition, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *TeacherBizService) BaseQuery(condition map[string]any) ([]*model.TeacherDTO, error) {
	entities := make([]*model.Teacher, 0)
	db, err := v.repo.TableGORMDB()
	if err != nil {
		return nil, err
	}
	err = db.Where(condition).Order(v.DefaultOrderBy()).Limit(v.MaxQuerySize()).Scan(&entities).Error
	if err != nil {
		return nil, err
	}
	return model.TeacherSlice(entities).ToDTOs()
}

func (v *TeacherBizService) BaseQueryPage(condition map[string]any, pager *webcloud.Pager[model.TeacherDTO]) error {
	countDB, err := v.repo.TableGORMDB()
	if err != nil {
		return err
	}
	if err = countDB.Where(condition).Count(&pager.Total).Error; err != nil || pager.Total == 0 {
		return err
	}
	entities := make([]*model.Teacher, 0, pager.Size)
	pageDB, err := v.repo.TableGORMDB()
	if err != nil {
		return err
	}
	if err = pageDB.Where(condition).
		Order(v.DefaultOrderBy()).
		Limit(pager.Size).
		Offset((pager.Number - 1) * pager.Size).
		Scan(&entities).Error; err != nil {
		return err
	}
	records, err := model.TeacherSlice(entities).ToDTOs()
	if err != nil {
		return err
	}
	pager.Records = records
	return nil
}

func (v *TeacherBizService) BaseModifyByID(update, condition map[string]any) (int64, error) {
	return v.repo.ModifyByMap(update, condition)
}

func (v *TeacherBizService) BaseRemoveByID(condition map[string]any) (int64, error) {
	return v.repo.RemoveByMap(condition)
}

func (v *TeacherBizService) QueryByID(id int64) (*model.TeacherDTO, error) {
	var entity model.Teacher
	row, err := v.repo.QueryByID(id, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *TeacherBizService) QueryOneByCond(condition *model.TeacherQDTO) (*model.TeacherDTO, error) {
	if condition == nil {
		return nil, webcloud.ErrBadRequestParameters
	}
	entityCondition, err := condition.ToT()
	if err != nil {
		return nil, err
	}
	var entity model.Teacher
	row, err := v.repo.QueryOneByCond(entityCondition, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *TeacherBizService) QueryByCond(condition *model.TeacherQDTO) ([]*model.TeacherDTO, error) {
	if condition == nil {
		return nil, webcloud.ErrBadRequestParameters
	}
	entityCondition, err := condition.ToT()
	if err != nil {
		return nil, err
	}
	entities := make([]*model.Teacher, 0)
	if _, err = v.repo.QueryByCond(entityCondition, v.DefaultOrderBy(), &entities); err != nil {
		return nil, err
	}
	return model.TeacherSlice(entities).ToDTOs()
}

func (v *TeacherBizService) QueryPage(pager webcloud.PagerDTO[model.TeacherQDTO]) (webcloud.Pager[model.TeacherDTO], error) {
	condition, err := pager.Condition.ToT()
	if err != nil {
		return webcloud.Pager[model.TeacherDTO]{}, err
	}
	page := databasecloud.Pager[model.Teacher]{
		Number: pager.Number,
		Size:   pager.Size,
	}
	if err = v.repo.QueryPageByCond(condition, rds.PageQuery{OrderBySQL: v.DefaultOrderBy()}, &page); err != nil {
		return webcloud.Pager[model.TeacherDTO]{}, err
	}
	records, err := model.TeacherSlice(page.Records).ToDTOs()
	if err != nil {
		return webcloud.Pager[model.TeacherDTO]{}, err
	}
	return webcloud.Pager[model.TeacherDTO]{
		Records: records,
		Total:   page.Total,
		Number:  page.Number,
		Size:    page.Size,
	}, nil
}

func (v *TeacherBizService) ModifyByID(id int64, updated *model.TeacherMDTO) (int64, error) {
	if updated == nil {
		return 0, webcloud.ErrBadRequestParameters
	}
	entity, err := updated.ToT()
	if err != nil {
		return 0, err
	}
	entity.ID = id
	return v.repo.ModifyByID(entity)
}

func (v *TeacherBizService) ModifyByIDWithoutZeroFields(id int64, updated *model.TeacherMDTO) (int64, error) {
	if updated == nil {
		return 0, webcloud.ErrBadRequestParameters
	}
	entity, err := updated.ToT()
	if err != nil {
		return 0, err
	}
	entity.ID = id
	return v.repo.ModifyByIDWithoutZeroFields(entity)
}

func (v *TeacherBizService) ModifyByIDWithMap(id int64, updated map[string]any) (int64, error) {
	if len(updated) == 0 {
		return 0, webcloud.ErrBadRequestParameters
	}
	return v.repo.ModifyByIDWithMap(updated, id)
}

func (v *TeacherBizService) RemoveByID(id int64) (int64, error) {
	return v.repo.RemoveByID(id)
}

func (v *TeacherBizService) RemoveByCond(condition *model.TeacherQDTO) (int64, error) {
	if condition == nil {
		return 0, webcloud.ErrBadRequestParameters
	}
	entity, err := condition.ToT()
	if err != nil {
		return 0, err
	}
	return v.repo.RemoveByCond(entity)
}

func (v *TeacherBizService) RemoveByMap(condition map[string]any) (int64, error) {
	if len(condition) == 0 {
		return 0, webcloud.ErrBadRequestParameters
	}
	return v.repo.RemoveByMap(condition)
}
