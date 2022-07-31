package misc

import (
	"context"

	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm/clause"
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
	if result := m.DB.WithContext(ctx).Joins("Role").Where("Email LIKE ?", mail).First(&user); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

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

func (m *Misc) CreateUser(ctx context.Context, newUser users.ChangeUser) (users.User, error) {
	dbUser := models.User{
		Name:         newUser.Name,
		Email:        newUser.Email,
		Hash:         newUser.Hash,
		Distribution: newUser.Distribution,
		RoleID:       &newUser.RoleID,
	}
	if err := m.DB.WithContext(ctx).Omit(clause.Associations).Create(&dbUser).Error; err != nil {
		return users.User{}, err
	}
	return m.User(ctx, dbUser.ID)
}

func (m *Misc) UpdateUser(ctx context.Context, user users.ChangeUser) error {
	dbUser := models.User{
		Name:         user.Name,
		Email:        user.Email,
		Hash:         user.Hash,
		Distribution: user.Distribution,
		RoleID:       &user.RoleID,
	}
	return m.DB.WithContext(ctx).Model(&models.User{ID: user.ID}).Updates(&dbUser).Error
}

func (m *Misc) DeleteUser(ctx context.Context, id uint64) error {
	return m.DB.WithContext(ctx).Delete(&models.User{ID: id}).Error
}
