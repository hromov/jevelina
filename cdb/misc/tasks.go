package misc

import (
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (m *Misc) Tasks(filter models.ListFilter) (*models.TasksResponse, error) {
	cr := &models.TasksResponse{}
	//How to make joins work?.Joins("Contacts")
	q := m.DB.Preload(clause.Associations).Limit(filter.Limit).Offset(filter.Offset)
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
		return nil, result.Error
	}
	return cr, nil
}

func (m *Misc) Task(ID uint64) (*models.Task, error) {
	var item models.Task
	if result := m.DB.Preload(clause.Associations).First(&item, ID); result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (m *Misc) TaskTypes() ([]models.TaskType, error) {
	var items []models.TaskType
	if result := m.DB.Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (m *Misc) TaskType(ID uint8) (*models.TaskType, error) {
	var item models.TaskType
	if result := m.DB.First(&item, ID); result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func DeleteTaskByParent(db *gorm.DB, parentID uint64) error {
	return db.Delete(&models.Task{}, "parent_id = ?", parentID).Error
}
