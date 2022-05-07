package base

import (
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/files"
)

var db *cdb.CDB

const bucketName = "jevelina"

func Init(dsn string) error {
	var err error
	if db, err = cdb.Init(dsn, bucketName); err != nil {
		return err
	}
	return nil
}

func GetDB() *cdb.CDB {
	return db
}

func Files() *files.FilesService {
	return db.Files()
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
