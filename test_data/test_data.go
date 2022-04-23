package testdata

import (
	"fmt"
	"log"

	"github.com/hromov/cdb/misc"
	"github.com/hromov/cdb/models"
	"github.com/hromov/jevelina/base"
)

var base_roles = []models.Role{
	{Role: "Admin"},
	{Role: "User"},
}
var r1 = uint8(1)
var r2 = uint8(2)
var test_users = []models.User{
	{
		Name:   "User 1",
		Email:  "user_1@gmail.com",
		RoleID: &r1,
	},
	{
		Name:   "User 2",
		Email:  "user_2@gmail.com",
		RoleID: &r2,
	},
}

var test_sources = []models.Source{
	{
		Name: "Source 1",
	},
	{
		Name: "Source 2",
	},
}

func create_test_roles() {
	db := base.GetDB()
	for _, r := range base_roles {
		db.Create(&r)
	}
}

func create_test_users() {
	db := base.GetDB()
	for _, u := range test_users {
		db.Create(&u)
	}
}

func create_test_sources() {
	db := base.GetDB()
	for _, u := range test_sources {
		db.Create(&u)
	}
}

func Fill() {
	m := &misc.Misc{DB: base.GetDB().DB}
	roles, err := m.Roles()
	if err != nil || len(roles) == 0 {
		create_test_roles()
	} else {
		log.Println(roles)
	}

	if users, err := m.Users(); err != nil || len(users) == 0 {
		create_test_users()
	} else {
		// log.Println(users)
		// users[0].Role = roles[0]
		// users[0].RoleID = roles[0].ID
		// Test(&users[0])
		fmt.Printf("%v\n", users[0])
		fmt.Printf("%v\n", users[1])

		// log.Println(users[0].Role)
	}

	sources, err := m.Sources()
	if err != nil || len(sources) == 0 {
		create_test_sources()
	} else {
		log.Println(sources)
	}
}

func Test(user *models.User) error {
	db := base.GetDB()
	r_user, err := db.Update(user)
	if err != nil {
		return err
	}
	log.Println(r_user)
	return nil
}
