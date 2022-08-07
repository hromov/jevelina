package misc

import (
	"context"

	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func (m *Misc) GetProduct(ctx context.Context, id uint32) (misc.Product, error) {
	var item models.Product
	if result := m.DB.WithContext(ctx).First(&item, id); result.Error != nil {
		return misc.Product{}, result.Error
	}
	return item.ToDomain(), nil
}

func (m *Misc) GetProductByName(ctx context.Context, name string) (misc.Product, error) {
	var item models.Product
	if result := m.DB.WithContext(ctx).Where("name LIKE ?", name).First(&item); result.Error != nil {
		return misc.Product{}, result.Error
	}
	return item.ToDomain(), nil
}

func (m *Misc) ListProducts(ctx context.Context) ([]misc.Product, error) {
	var items []models.Product
	if err := m.DB.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, err
	}
	products := make([]misc.Product, len(items))
	for i, p := range items {
		products[i] = p.ToDomain()
	}
	return products, nil
}

func (m *Misc) CreateProduct(ctx context.Context, p misc.Product) (misc.Product, error) {
	dbProduct := models.ProductFromDomain(p)
	if err := m.DB.WithContext(ctx).Omit(clause.Associations).Create(&dbProduct).Error; err != nil {
		return misc.Product{}, err
	}
	return dbProduct.ToDomain(), nil
}

func (m *Misc) UpdateProduct(ctx context.Context, p misc.Product) error {
	return m.DB.WithContext(ctx).Model(&models.Product{}).Where("id", p.ID).Update("name", p.Name).Error
}

func (m *Misc) DeleteProduct(ctx context.Context, id uint32) error {
	return m.DB.WithContext(ctx).Delete(&models.Product{}, id).Error
}
