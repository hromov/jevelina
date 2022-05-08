package auth

import (
	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm/clause"
)

var AdminRoleID = uint8(1)
var UserRoleID = uint8(2)

func GetInitUsers() []*models.User {
	users := []*models.User{
		{
			Name:   "Admin User",
			Email:  "melifarowow@gmail.com",
			Hash:   "melifarowow@gmail.com",
			RoleID: &AdminRoleID,
		},
		{
			Name:   "Random User",
			Email:  "random@random.org",
			Hash:   "random@random.org",
			RoleID: &UserRoleID,
		},
	}
	return users
}

func GetInitRoles() []*models.Role {
	roles := []*models.Role{
		{ID: AdminRoleID, Role: "Admin"},
		{ID: UserRoleID, Role: "User"},
	}
	return roles
}

func GetBaseRole() (*models.Role, error) {
	return base.GetDB().Misc().Role(UserRoleID)
}

func CreateInitUsers() ([]*models.User, error) {
	users := GetInitUsers()

	for _, user := range users {
		if err := base.GetDB().DB.Omit(clause.Associations).Create(user).Error; err != nil {
			return nil, err
		}
	}
	return users, nil
}

func CreateInitRoles() ([]*models.Role, error) {
	roles := GetInitRoles()
	for _, role := range roles {
		if err := base.GetDB().DB.Create(role).Error; err != nil {
			return nil, err
		}
	}
	return roles, nil
}
