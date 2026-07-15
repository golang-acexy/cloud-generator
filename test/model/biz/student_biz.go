package biz

import (
	"github.com/golang-acexy/cloud-database/databasecloud"
	"github.com/golang-acexy/cloud-database/databasecloud/rds"
	"github.com/golang-acexy/cloud-generator/test/model"
	"github.com/golang-acexy/cloud-generator/test/model/repo"
	"github.com/golang-acexy/cloud-web/webcloud"
)

var studentBizService = &StudentBizService{
	repo: repo.NewStudentRepo(),
}

func NewStudentBizService() *StudentBizService {
	return studentBizService
}

type StudentBizService struct {
	repo repo.StudentRepo
}

func (v *StudentBizService) MaxQuerySize() int {
	return 500
}

func (v *StudentBizService) DefaultOrderBy() string {
	return "id desc"
}

func (v *StudentBizService) Save(save *model.StudentSDTO) (int64, error) {
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

func (v *StudentBizService) BaseQueryByID(condition map[string]any) (*model.StudentDTO, error) {
	var entity model.Student
	row, err := v.repo.QueryOneByMap(condition, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *StudentBizService) BaseQueryOne(condition map[string]any) (*model.StudentDTO, error) {
	var entity model.Student
	row, err := v.repo.QueryOneByMap(condition, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *StudentBizService) BaseQuery(condition map[string]any) ([]*model.StudentDTO, error) {
	entities := make([]*model.Student, 0)
	db, err := v.repo.TableGORMDB()
	if err != nil {
		return nil, err
	}
	err = db.Where(condition).Order(v.DefaultOrderBy()).Limit(v.MaxQuerySize()).Scan(&entities).Error
	if err != nil {
		return nil, err
	}
	return model.StudentSlice(entities).ToDTOs()
}

func (v *StudentBizService) BaseQueryPage(condition map[string]any, pager *webcloud.Pager[model.StudentDTO]) error {
	countDB, err := v.repo.TableGORMDB()
	if err != nil {
		return err
	}
	if err = countDB.Where(condition).Count(&pager.Total).Error; err != nil || pager.Total == 0 {
		return err
	}
	entities := make([]*model.Student, 0, pager.Size)
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
	records, err := model.StudentSlice(entities).ToDTOs()
	if err != nil {
		return err
	}
	pager.Records = records
	return nil
}

func (v *StudentBizService) BaseModifyByID(update, condition map[string]any) (int64, error) {
	return v.repo.ModifyByMap(update, condition)
}

func (v *StudentBizService) BaseRemoveByID(condition map[string]any) (int64, error) {
	return v.repo.RemoveByMap(condition)
}

func (v *StudentBizService) QueryByID(id int64) (*model.StudentDTO, error) {
	var entity model.Student
	row, err := v.repo.QueryByID(id, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *StudentBizService) QueryOneByCond(condition *model.StudentQDTO) (*model.StudentDTO, error) {
	if condition == nil {
		return nil, webcloud.ErrBadRequestParameters
	}
	entityCondition, err := condition.ToT()
	if err != nil {
		return nil, err
	}
	var entity model.Student
	row, err := v.repo.QueryOneByCond(entityCondition, &entity)
	if err != nil || row == 0 {
		return nil, err
	}
	return entity.ToDTO()
}

func (v *StudentBizService) QueryByCond(condition *model.StudentQDTO) ([]*model.StudentDTO, error) {
	if condition == nil {
		return nil, webcloud.ErrBadRequestParameters
	}
	entityCondition, err := condition.ToT()
	if err != nil {
		return nil, err
	}
	entities := make([]*model.Student, 0)
	if _, err = v.repo.QueryByCond(entityCondition, v.DefaultOrderBy(), &entities); err != nil {
		return nil, err
	}
	return model.StudentSlice(entities).ToDTOs()
}

func (v *StudentBizService) QueryPage(pager webcloud.PagerDTO[model.StudentQDTO]) (webcloud.Pager[model.StudentDTO], error) {
	condition, err := pager.Condition.ToT()
	if err != nil {
		return webcloud.Pager[model.StudentDTO]{}, err
	}
	page := databasecloud.Pager[model.Student]{
		Number: pager.Number,
		Size:   pager.Size,
	}
	if err = v.repo.QueryPageByCond(condition, rds.PageQuery{OrderBySQL: v.DefaultOrderBy()}, &page); err != nil {
		return webcloud.Pager[model.StudentDTO]{}, err
	}
	records, err := model.StudentSlice(page.Records).ToDTOs()
	if err != nil {
		return webcloud.Pager[model.StudentDTO]{}, err
	}
	return webcloud.Pager[model.StudentDTO]{
		Records: records,
		Total:   page.Total,
		Number:  page.Number,
		Size:    page.Size,
	}, nil
}

func (v *StudentBizService) ModifyByID(id int64, updated *model.StudentMDTO) (int64, error) {
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

func (v *StudentBizService) ModifyByIDWithoutZeroFields(id int64, updated *model.StudentMDTO) (int64, error) {
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

func (v *StudentBizService) ModifyByIDWithMap(id int64, updated map[string]any) (int64, error) {
	if len(updated) == 0 {
		return 0, webcloud.ErrBadRequestParameters
	}
	return v.repo.ModifyByIDWithMap(updated, id)
}

func (v *StudentBizService) RemoveByID(id int64) (int64, error) {
	return v.repo.RemoveByID(id)
}

func (v *StudentBizService) RemoveByCond(condition *model.StudentQDTO) (int64, error) {
	if condition == nil {
		return 0, webcloud.ErrBadRequestParameters
	}
	entity, err := condition.ToT()
	if err != nil {
		return 0, err
	}
	return v.repo.RemoveByCond(entity)
}

func (v *StudentBizService) RemoveByMap(condition map[string]any) (int64, error) {
	if len(condition) == 0 {
		return 0, webcloud.ErrBadRequestParameters
	}
	return v.repo.RemoveByMap(condition)
}
