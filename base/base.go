package base

import (
	"github.com/hromov/cdb"
	"github.com/hromov/cdb/contacts"
	"github.com/hromov/cdb/leads"
	"github.com/hromov/cdb/misc"
)

const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

var db *cdb.CDB

func Init() error {
	var err error
	if db, err = cdb.Init(dsn); err != nil {
		return err
	}
	return nil
}

func GetDB() *cdb.CDB {
	return db
}

func Contacts() *contacts.Contacts {
	return db.Contacts()
}

func Leads() *leads.Leads {
	return db.Leads()
}

func Misc() *misc.Misc {
	return db.Misc()
}
