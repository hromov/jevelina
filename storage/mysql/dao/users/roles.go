package users

import (
	"context"

	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func (u *Users) Roles(ctx context.Context) ([]users.Role, error) {
	var roles []models.Role
	if result := u.db.WithContext(ctx).Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	respRoles := make([]users.Role, len(roles))
	for i, r := range roles {
		respRoles[i] = r.ToDomain()
	}
	return respRoles, nil
}

func (u *Users) Role(ctx context.Context, ID uint8) (users.Role, error) {
	var role models.Role
	if err := u.db.WithContext(ctx).First(&role, ID).Error; err != nil {
		return users.Role{}, err
	}
	return role.ToDomain(), nil
}

func (u *Users) CreateRole(ctx context.Context, role users.Role) (users.Role, error) {
	dbRole := models.Role{
		Priority: role.Priority,
		Role:     role.Role,
	}
	if err := u.db.WithContext(ctx).Omit(clause.Associations).Create(&dbRole).Error; err != nil {
		return users.Role{}, err
	}
	return dbRole.ToDomain(), nil
}

func (u *Users) UpdateRole(ctx context.Context, role users.Role) error {
	dbRole := models.Role{
		Priority: role.Priority,
		Role:     role.Role,
	}
	return u.db.WithContext(ctx).Model(&models.Role{ID: role.ID}).Updates(&dbRole).Error
}

func (u *Users) DeleteRole(ctx context.Context, id uint8) error {
	return u.db.WithContext(ctx).Delete(&models.Role{ID: id}).Error
}
