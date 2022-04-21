package testdata

import (
	"fmt"
	"log"

	"github.com/hromov/cdb"
)

var base_roles = []cdb.Role{
	{Role: "Admin"},
	{Role: "User"},
}
var r1 = uint8(1)
var r2 = uint8(2)
var test_users = []cdb.User{
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

var test_sources = []cdb.Source{
	{
		Name: "Source 1",
	},
	{
		Name: "Source 2",
	},
}

func create_test_roles() {
	for _, r := range base_roles {
		cdb.Create(&r)
	}
}

func create_test_users() {
	for _, u := range test_users {
		cdb.Create(&u)
	}
}

func create_test_sources() {
	for _, u := range test_sources {
		cdb.Create(&u)
	}
}

func Fill() {
	roles, err := cdb.Roles()
	if err != nil || len(roles) == 0 {
		create_test_roles()
	} else {
		log.Println(roles)
	}

	if users, err := cdb.Users(); err != nil || len(users) == 0 {
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

	sources, err := cdb.Sources()
	if err != nil || len(sources) == 0 {
		create_test_sources()
	} else {
		log.Println(sources)
	}
}

func Test(user *cdb.User) error {
	r_user, err := cdb.Update(user)
	if err != nil {
		return err
	}
	log.Println(r_user)
	return nil
}
