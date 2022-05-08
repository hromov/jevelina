package auth

import (
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
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
	return cdb.Misc().Role(UserRoleID)
}

func CreateInitUsers(db *gorm.DB) ([]*models.User, error) {
	users := GetInitUsers()

	for _, user := range users {
		if err := db.Omit(clause.Associations).Create(user).Error; err != nil {
			return nil, err
		}
	}
	return users, nil
}

func CreateInitRoles(db *gorm.DB) ([]*models.Role, error) {
	roles := GetInitRoles()
	for _, role := range roles {
		if err := db.Create(role).Error; err != nil {
			return nil, err
		}
	}
	return roles, nil
}
