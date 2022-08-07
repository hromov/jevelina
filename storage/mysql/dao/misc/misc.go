package misc

import (
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
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

func (m *Misc) SourceByName(name string) (*models.Source, error) {
	var source models.Source
	if result := m.DB.Where("name LIKE ?", name).First(&source); result.Error != nil {
		return nil, result.Error
	}
	return &source, nil
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

func (m *Misc) ProductByName(name string) (*models.Product, error) {
	var item models.Product
	if result := m.DB.Where("name LIKE ?", name).First(&item); result.Error != nil {
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

func (m *Misc) ManufacturerByName(name string) (*models.Manufacturer, error) {
	var item models.Manufacturer
	if result := m.DB.Where("name LIKE ?", name).First(&item); result.Error != nil {
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
