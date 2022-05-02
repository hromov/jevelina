package misc

import (
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Misc struct {
	*gorm.DB
}

func (m *Misc) Sources() ([]models.Source, error) {
	var sources []models.Source
	if result := m.DB.Find(&sources); result.Error != nil {
		return nil, result.Error
	}
	return sources, nil
}

func (m *Misc) Source(ID uint8) (*models.Source, error) {
	var source models.Source
	if result := m.DB.First(&source, ID); result.Error != nil {
		return nil, result.Error
	}
	return &source, nil
}

func (m *Misc) Users() ([]models.User, error) {
	var users []models.User
	if result := m.DB.Joins("Role").Find(&users); result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (m *Misc) User(ID uint64) (*models.User, error) {
	var user models.User
	if result := m.DB.Joins("Role").First(&user, ID); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (m *Misc) UserExist(mail string) (bool, error) {
	var exists bool
	if err := m.DB.Model(&models.User{}).Select("count(*) > 0").Where("Email LIKE ?", mail).Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (m *Misc) UserByEmail(mail string) (*models.User, error) {
	var user models.User
	if result := m.DB.Joins("Role").Where("Email LIKE ?", mail).First(&user); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (m *Misc) Roles() ([]models.Role, error) {
	var roles []models.Role
	if result := m.DB.Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (m *Misc) Role(ID uint8) (*models.Role, error) {
	var role models.Role
	if result := m.DB.First(&role, ID); result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (m *Misc) Products() ([]models.Product, error) {
	var items []models.Product
	if result := m.DB.Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (m *Misc) Product(ID uint32) (*models.Product, error) {
	var item models.Product
	if result := m.DB.First(&item, ID); result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (m *Misc) Manufacturers() ([]models.Manufacturer, error) {
	var items []models.Manufacturer
	if result := m.DB.Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (m *Misc) Manufacturer(ID uint16) (*models.Manufacturer, error) {
	var item models.Manufacturer
	if result := m.DB.First(&item, ID); result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (m *Misc) Steps() ([]models.Step, error) {
	var items []models.Step
	if result := m.DB.Order("`order`").Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (m *Misc) Step(ID uint8) (*models.Step, error) {
	var item models.Step
	if result := m.DB.First(&item, ID); result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (m *Misc) Tags() ([]models.Tag, error) {
	var items []models.Tag
	if result := m.DB.Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (m *Misc) Tag(ID uint8) (*models.Tag, error) {
	var item models.Tag
	if result := m.DB.First(&item, ID); result.Error != nil {
		return nil, result.Error

	}
	return &item, nil
}

func (m *Misc) Tasks(filter models.ListFilter) (*models.TasksResponse, error) {
	cr := &models.TasksResponse{}
	//How to make joins work?.Joins("Contacts")
	q := m.DB.Preload(clause.Associations).Order("created_at asc").Limit(filter.Limit).Offset(filter.Offset)
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
		q = q.Where("dead_line IS NOT NULL").Where("completed = false")
	}
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
