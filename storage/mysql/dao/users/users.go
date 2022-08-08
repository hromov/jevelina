package users

import (
	"context"
	"log"

	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "user", err.Error())
	}

	if err := db.AutoMigrate(&models.Role{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "role", err.Error())
	}

	user := models.User{ID: 1}
	if err := db.First(&user).Error; err != nil {
		if err := InitUsers(db); err != nil {
			log.Printf("Can't create base roles error: %s", err.Error())
		}
	}
	return &Users{db}
}

func (u *Users) GetUsers(ctx context.Context) ([]users.User, error) {
	var dbUsers []models.User
	if err := u.db.WithContext(ctx).Joins("Role").Find(&dbUsers).Error; err != nil {
		return nil, err
	}
	respUsers := make([]users.User, len(dbUsers))
	for i, u := range dbUsers {
		respUsers[i] = u.ToDomain()
	}
	return respUsers, nil
}

func (u *Users) User(ctx context.Context, ID uint64) (users.User, error) {
	var user models.User
	if result := u.db.WithContext(ctx).Joins("Role").First(&user, ID); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

func (u *Users) UserExist(ctx context.Context, mail string) (bool, error) {
	var exists bool
	if err := u.db.WithContext(ctx).Model(&models.User{}).Select("count(*) > 0").Where("Email LIKE ?", mail).Find(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

func (u *Users) UserByEmail(ctx context.Context, mail string) (users.User, error) {
	var user models.User
	if result := u.db.WithContext(ctx).Joins("Role").Where("Email LIKE ?", mail).First(&user); result.Error != nil {
		return users.User{}, result.Error
	}
	return user.ToDomain(), nil
}

func (u *Users) CreateUser(ctx context.Context, user users.ChangeUser) (users.User, error) {
	dbUser := models.UserFromDomain(user)
	if err := u.db.WithContext(ctx).Omit(clause.Associations).Create(&dbUser).Error; err != nil {
		return users.User{}, err
	}
	return u.User(ctx, dbUser.ID)
}

func (u *Users) UpdateUser(ctx context.Context, user users.ChangeUser) error {
	dbUser := models.UserFromDomain(user)
	return u.db.WithContext(ctx).Model(&models.User{ID: user.ID}).Updates(&dbUser).Error
}

func (u *Users) DeleteUser(ctx context.Context, id uint64) error {
	return u.db.WithContext(ctx).Delete(&models.User{ID: id}).Error
}
