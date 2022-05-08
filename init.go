package main

import (
	"fmt"

	"github.com/hromov/amoimport"
	"github.com/hromov/jevelina/auth"
	"gorm.io/gorm"
)

const leads = "_import/amocrm_export_leads_2022-04-20.csv"
const contacts = "_import/amocrm_export_contacts_2022-04-20.csv"
const rowsToImport = 3000

//Init DataBase - push base roles, admins, import "n" rows from
func Init(db *gorm.DB) error {
	if _, err := auth.CreateInitRoles(db); err != nil {
		return fmt.Errorf("Can't create base roles error: %s", err.Error())
	}

	if _, err := auth.CreateInitUsers(db); err != nil {
		return fmt.Errorf("Can't create init users error: %s", err.Error())
	}

	if err := amoimport.Import(db, leads, contacts, rowsToImport); err != nil {
		return fmt.Errorf("Can't import error: %s", err.Error())
	}
	return nil
}
