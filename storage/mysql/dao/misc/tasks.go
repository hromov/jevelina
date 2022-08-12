package misc

import (
	"context"

	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"github.com/hromov/jevelina/useCases/tasks"
	"gorm.io/gorm/clause"
)

func (m *Misc) GetTasks(ctx context.Context, filter tasks.Filter) (tasks.TasksResponse, error) {
	cr := &models.TasksResponse{}
	q := m.db.WithContext(ctx).Preload(clause.Associations).Limit(filter.Limit).Offset(filter.Offset)
	if filter.Query != "" {
		q = q.Where("name LIKE ?", "%"+filter.Query+"%")
	}
	if filter.ParentID != 0 {
		q = q.Where("parent_id = ?", filter.ParentID)
	}
	if filter.ResponsibleID != 0 {
		q = q.Where("responsible_id = ?", filter.ResponsibleID)
	}
	if !filter.MinDate.IsZero() {
		q = q.Where("dead_line >= ?", filter.MinDate)
	}
	if !filter.MaxDate.IsZero() {
		q = q.Where("dead_line < ?", filter.MaxDate)
	}
	if !filter.MinDate.IsZero() || !filter.MaxDate.IsZero() {
		q = q.Where("dead_line IS NOT NULL").Where("completed = false").Order("dead_line desc")
	}
	q.Order("created_at asc")
	if result := q.Find(&cr.Tasks).Count(&cr.Total); result.Error != nil {
		return tasks.TasksResponse{}, result.Error
	}
	taskList := make([]tasks.Task, len(cr.Tasks))
	for i, t := range cr.Tasks {
		taskList[i] = t.ToDomain()
	}
	return tasks.TasksResponse{Tasks: taskList, Total: cr.Total}, nil
}

func (m *Misc) GetTask(ctx context.Context, id uint64) (tasks.Task, error) {
	var item models.Task
	if result := m.db.WithContext(ctx).Preload(clause.Associations).First(&item, id); result.Error != nil {
		return tasks.Task{}, result.Error
	}
	return m.GetTask(ctx, id)
}

func (m *Misc) DeleteTaskByParent(ctx context.Context, parentID uint64) error {
	return m.db.WithContext(ctx).Delete(&models.Task{}, "parent_id = ?", parentID).Error
}

func (m *Misc) DeleteTask(ctx context.Context, id uint64) error {
	return m.db.WithContext(ctx).Delete(&models.Task{ID: id}).Error
}

func (m *Misc) CreateTask(ctx context.Context, task tasks.TaskData) (tasks.Task, error) {
	dbTask := models.TaskFromTaskData(task)
	if err := m.db.WithContext(ctx).Omit(clause.Associations).Create(&dbTask).Error; err != nil {
		return tasks.Task{}, err
	}
	return dbTask.ToDomain(), nil
}

func (m *Misc) UpdateTask(ctx context.Context, task tasks.TaskData) error {
	dbTask := models.TaskFromTaskData(task)
	return m.db.WithContext(ctx).Model(&models.Task{}).Where("id", task.ID).Updates(&dbTask).Error
}

// TODO: don't use task types rn
// func (m *Misc) TaskTypes() ([]models.TaskType, error) {
// 	var items []models.TaskType
// 	if result := m.db.Find(&items); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return items, nil
// }

// func (m *Misc) TaskType(ID uint8) (*models.TaskType, error) {
// 	var item models.TaskType
// 	if result := m.db.First(&item, ID); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &item, nil
// }
