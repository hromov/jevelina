package main

import (
	"fmt"

	"github.com/hromov/jevelina/auth"
	"gorm.io/gorm"
)

//Init DataBase - push base roles, admins, import "n" rows from
func InitUsers(db *gorm.DB) error {
	if _, err := auth.CreateInitRoles(db); err != nil {
		return fmt.Errorf("Can't create base roles error: %s", err.Error())
	}

	if _, err := auth.CreateInitUsers(db); err != nil {
		return fmt.Errorf("Can't create init users error: %s", err.Error())
	}
	return nil
}
