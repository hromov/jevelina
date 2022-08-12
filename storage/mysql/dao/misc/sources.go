package misc

import (
	"context"

	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func (m *Misc) GetSource(ctx context.Context, id uint32) (misc.Source, error) {
	var item models.Source
	if result := m.db.WithContext(ctx).First(&item, id); result.Error != nil {
		return misc.Source{}, result.Error
	}
	return item.ToDomain(), nil
}

func (m *Misc) GetSourceByName(ctx context.Context, name string) (misc.Source, error) {
	var item models.Source
	if result := m.db.WithContext(ctx).Where("name LIKE ?", name).First(&item); result.Error != nil {
		return misc.Source{}, result.Error
	}
	return item.ToDomain(), nil
}

func (m *Misc) ListSources(ctx context.Context) ([]misc.Source, error) {
	var items []models.Source
	if err := m.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	products := make([]misc.Source, len(items))
	for i, p := range items {
		products[i] = p.ToDomain()
	}
	return products, nil
}

func (m *Misc) CreateSource(ctx context.Context, p misc.Source) (misc.Source, error) {
	dbSource := models.SourceFromDomain(p)
	if err := m.db.WithContext(ctx).Omit(clause.Associations).Create(&dbSource).Error; err != nil {
		return misc.Source{}, err
	}
	return dbSource.ToDomain(), nil
}

func (m *Misc) UpdateSource(ctx context.Context, p misc.Source) error {
	return m.db.WithContext(ctx).Model(&models.Source{}).Where("id", p.ID).Update("name", p.Name).Error
}

func (m *Misc) DeleteSource(ctx context.Context, id uint32) error {
	return m.db.WithContext(ctx).Delete(&models.Source{}, id).Error
}
