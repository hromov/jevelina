package cdb

import (
	"fmt"

	"github.com/hromov/jevelina/cdb/contacts"
	"github.com/hromov/jevelina/cdb/files"
	"github.com/hromov/jevelina/cdb/finance"
	"github.com/hromov/jevelina/cdb/leads"
	"github.com/hromov/jevelina/cdb/misc"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testDSN = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

type CDB struct {
	BucketName string
	*gorm.DB
}

var currentDB *CDB

func GetDB() *CDB {
	return currentDB
}

//test part
func Contacts() *contacts.Contacts {
	return &contacts.Contacts{DB: currentDB.DB}
}

func Leads() *leads.Leads {
	return &leads.Leads{DB: currentDB.DB}
}

func Misc() *misc.Misc {
	return &misc.Misc{DB: currentDB.DB}
}

func Finance() *finance.Finance {
	return &finance.Finance{DB: currentDB.DB}
}

func Files() *files.FilesService {
	return &files.FilesService{DB: currentDB.DB, BucketName: currentDB.BucketName}
}

func (db *CDB) SetBucket(bucketName string) {
	db.BucketName = bucketName
}

func Open(dsn string) (*CDB, error) {
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 30% fester but not so safe... let's give it a try
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database error: %s", err.Error())
	}
	return &CDB{DB: db}, nil
}

func OpenTest() (*CDB, error) {
	db, err := gorm.Open(mysql.Open(testDSN), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database error: %s", err.Error())
	}
	return &CDB{DB: db}, nil
}

func OpenAndInit(dsn string) (*CDB, error) {
	db, err := Open(dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Init()
}

func OpenAndInitTest() (*CDB, error) {
	db, err := OpenTest()
	if err != nil {
		return nil, err
	}
	return db, db.Init()
}

func (db *CDB) Init() error {

	currentDB = db
	// if table exist - do nothink, if not - create init structure with test data
	if !currentDB.DB.Migrator().HasTable("roles") {
		if err := db.AutoMigrate(&models.Role{}); err != nil {
			return err
		}
	}
	if !currentDB.DB.Migrator().HasTable("contacts") {
		if err := db.AutoMigrate(&models.Contact{}); err != nil {
			return err
		}
	}

	if !currentDB.DB.Migrator().HasTable("leads") {
		if err := db.AutoMigrate(&models.Lead{}); err != nil {
			return err
		}
	}

	if !currentDB.DB.Migrator().HasTable("tasks") {
		if err := db.AutoMigrate(&models.Task{}); err != nil {
			return err
		}
	}

	if !currentDB.DB.Migrator().HasTable("wallets") {
		if err := db.AutoMigrate(&models.Wallet{}); err != nil {
			return err
		}
	}

	if !currentDB.DB.Migrator().HasTable("transfers") {
		if err := db.AutoMigrate(&models.Transfer{}); err != nil {
			return err
		}
	}
	if !currentDB.DB.Migrator().HasTable("files") {
		if err := db.AutoMigrate(&models.File{}); err != nil {
			return err
		}
	}

	if !currentDB.DB.Migrator().HasTable("sources") {
		if err := db.AutoMigrate(&models.Source{}); err != nil {
			return err
		}
	}

	//TODO: check if works on clean -> Move to file -> Readme.MD
	user := models.User{ID: 1}
	if err := currentDB.DB.First(&user).Error; err != nil {
		if err := InitUsers(db.DB); err != nil {
			return fmt.Errorf("Can't create base roles error: %s", err.Error())
		}
	}

	return nil
}
