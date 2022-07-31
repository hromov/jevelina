package misc

import (
	"context"

	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/domain/users"
)

func (m *Misc) Users(ctx context.Context) ([]users.User, error) {
	var dbUsers []models.User
	if err := m.DB.WithContext(ctx).Joins("Role").Find(&dbUsers).Error; err != nil {
		return nil, err
	}
	respUsers := make([]users.User, len(dbUsers))
	for i, u := range dbUsers {
		respUsers[i] = u.ToDomain()
	}
	return respUsers, nil
}

func (m *Misc) User(ctx context.Context, ID uint64) (users.User, error) {
	var user models.User
	if result := m.DB.WithContext(ctx).Joins("Role").First(&user, ID); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

func (m *Misc) UserExist(ctx context.Context, mail string) (bool, error) {
	var exists bool
	if err := m.DB.WithContext(ctx).Model(&models.User{}).Select("count(*) > 0").Where("Email LIKE ?", mail).Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (m *Misc) UserByEmail(ctx context.Context, mail string) (users.User, error) {
	var user models.User
	if result := m.DB.Joins("Role").Where("Email LIKE ?", mail).First(&user); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

func (m *Misc) Roles() ([]users.Role, error) {
	var roles []models.Role
	if result := m.DB.Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	respRoles := make([]users.Role, len(roles))
	for i, r := range roles {
		respRoles[i] = r.ToDomain()
	}
	return respRoles, nil
}

func (m *Misc) Role(ID uint8) (*models.Role, error) {
	var role models.Role
	if result := m.DB.First(&role, ID); result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}
