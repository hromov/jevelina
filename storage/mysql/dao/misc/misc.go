package misc

import (
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
)

type Misc struct {
	*gorm.DB
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
