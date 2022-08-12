package misc

import (
	"context"

	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func (m *Misc) GetManufacturer(ctx context.Context, id uint32) (misc.Manufacturer, error) {
	var item models.Manufacturer
	if result := m.db.WithContext(ctx).First(&item, id); result.Error != nil {
		return misc.Manufacturer{}, result.Error
	}
	return item.ToDomain(), nil
}

func (m *Misc) GetManufacturerByName(ctx context.Context, name string) (misc.Manufacturer, error) {
	var item models.Manufacturer
	if result := m.db.WithContext(ctx).Where("name LIKE ?", name).First(&item); result.Error != nil {
		return misc.Manufacturer{}, result.Error
	}
	return item.ToDomain(), nil
}

func (m *Misc) ListManufacturers(ctx context.Context) ([]misc.Manufacturer, error) {
	var items []models.Manufacturer
	if err := m.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	products := make([]misc.Manufacturer, len(items))
	for i, p := range items {
		products[i] = p.ToDomain()
	}
	return products, nil
}

func (m *Misc) CreateManufacturer(ctx context.Context, p misc.Manufacturer) (misc.Manufacturer, error) {
	dbManufacturer := models.ManufacturerFromDomain(p)
	if err := m.db.WithContext(ctx).Omit(clause.Associations).Create(&dbManufacturer).Error; err != nil {
		return misc.Manufacturer{}, err
	}
	return dbManufacturer.ToDomain(), nil
}

func (m *Misc) UpdateManufacturer(ctx context.Context, p misc.Manufacturer) error {
	return m.db.WithContext(ctx).Model(&models.Manufacturer{}).Where("id", p.ID).Update("name", p.Name).Error
}

func (m *Misc) DeleteManufacturer(ctx context.Context, id uint32) error {
	return m.db.WithContext(ctx).Delete(&models.Manufacturer{}, id).Error
}
