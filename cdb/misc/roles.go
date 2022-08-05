package misc

import (
	"context"

	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm/clause"
)

func (m *Misc) Roles(ctx context.Context) ([]users.Role, error) {
	var roles []models.Role
	if result := m.DB.WithContext(ctx).Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	respRoles := make([]users.Role, len(roles))
	for i, r := range roles {
		respRoles[i] = r.ToDomain()
	}
	return respRoles, nil
}

func (m *Misc) Role(ctx context.Context, ID uint8) (users.Role, error) {
	var role models.Role
	if err := m.DB.WithContext(ctx).First(&role, ID).Error; err != nil {
		return users.Role{}, err
	}
	return role.ToDomain(), nil
}

func (m *Misc) CreateRole(ctx context.Context, role users.Role) (users.Role, error) {
	dbRole := models.Role{
		Priority: role.Priority,
		Role:     role.Role,
	}
	if err := m.DB.WithContext(ctx).Omit(clause.Associations).Create(&dbRole).Error; err != nil {
		return users.Role{}, err
	}
	return dbRole.ToDomain(), nil
}

func (m *Misc) UpdateRole(ctx context.Context, role users.Role) error {
	dbRole := models.Role{
		Priority: role.Priority,
		Role:     role.Role,
	}
	return m.DB.WithContext(ctx).Model(&models.Role{ID: role.ID}).Updates(&dbRole).Error
}

func (m *Misc) DeleteRole(ctx context.Context, id uint8) error {
	return m.DB.WithContext(ctx).Delete(&models.Role{ID: id}).Error
}
