package leads

import (
	"context"

	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func (l *Leads) GetSteps(ctx context.Context) ([]leads.Step, error) {
	var items []models.Step
	if result := l.DB.WithContext(ctx).Order("`order`").Find(&items); result.Error != nil {
		return nil, result.Error
	}
	steps := make([]leads.Step, len(items))
	for i, s := range items {
		steps[i] = s.ToDomain()
	}
	return steps, nil
}

func (l *Leads) GetStep(ctx context.Context, id uint8) (leads.Step, error) {
	var item models.Step
	if err := l.DB.WithContext(ctx).First(&item, id).Error; err != nil {
		return leads.Step{}, err
	}
	return item.ToDomain(), nil
}

func (l *Leads) DefaultStep(ctx context.Context) (leads.Step, error) {
	var item models.Step
	if err := l.DB.WithContext(ctx).Where("`order` = 0").First(&item).Error; err != nil {
		return leads.Step{}, err
	}
	return item.ToDomain(), nil
}

func (l *Leads) CreateStep(ctx context.Context, s leads.Step) (leads.Step, error) {
	step := models.StepFromDomain(s)
	if err := l.DB.WithContext(ctx).Omit(clause.Associations).Create(&step).Error; err != nil {
		return leads.Step{}, err
	}
	return step.ToDomain(), nil
}

func (l *Leads) UpdateStep(ctx context.Context, s leads.Step) error {
	step := models.StepFromDomain(s)
	return l.DB.WithContext(ctx).Omit(clause.Associations).Where("id", s.ID).Updates(&step).Error
}

func (l *Leads) DeleteStep(ctx context.Context, id uint8) error {
	if err := l.DB.WithContext(ctx).Delete(&models.Step{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
