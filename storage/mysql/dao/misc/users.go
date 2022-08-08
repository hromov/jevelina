package misc

import (
	"context"

	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func (m *Misc) Users(ctx context.Context) ([]users.User, error) {
	var dbUsers []models.User
	if err := m.db.WithContext(ctx).Joins("Role").Find(&dbUsers).Error; err != nil {
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
	if result := m.db.WithContext(ctx).Joins("Role").First(&user, ID); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

func (m *Misc) UserExist(ctx context.Context, mail string) (bool, error) {
	var exists bool
	if err := m.db.WithContext(ctx).Model(&models.User{}).Select("count(*) > 0").Where("Email LIKE ?", mail).Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (m *Misc) UserByEmail(ctx context.Context, mail string) (users.User, error) {
	var user models.User
	if result := m.db.WithContext(ctx).Joins("Role").Where("Email LIKE ?", mail).First(&user); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

func (m *Misc) CreateUser(ctx context.Context, user users.ChangeUser) (users.User, error) {
	dbUser := models.UserFromDomain(user)
	if err := m.db.WithContext(ctx).Omit(clause.Associations).Create(&dbUser).Error; err != nil {
		return users.User{}, err
	}
	return m.User(ctx, dbUser.ID)
}

func (m *Misc) UpdateUser(ctx context.Context, user users.ChangeUser) error {
	dbUser := models.UserFromDomain(user)
	return m.db.WithContext(ctx).Model(&models.User{ID: user.ID}).Updates(&dbUser).Error
}

func (m *Misc) DeleteUser(ctx context.Context, id uint64) error {
	return m.db.WithContext(ctx).Delete(&models.User{ID: id}).Error
}
