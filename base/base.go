package base

import (
	"github.com/hromov/jevelina/cdb"
)

var db *cdb.CDB

func Init(dsn string) error {
	var err error
	if db, err = cdb.Init(dsn); err != nil {
		return err
	}
	return nil
}

func GetDB() *cdb.CDB {
	return db
}

// func Contacts() *contacts.Contacts {
// 	return &contacts.Contacts{
// 		DB: db.DB,
// 	}
// }

// func Leads() *leads.Leads {
// 	return &leads.Leads{
// 		DB: db.DB,
// 	}
// }

// func Misc() *misc.Misc {
// 	return &misc.Misc{
// 		DB: db.DB,
// 	}
// }
