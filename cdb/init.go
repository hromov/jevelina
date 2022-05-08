package cdb

import (
	"errors"
	"fmt"
	"log"

	"github.com/hromov/jevelina/cdb/contacts"
	"github.com/hromov/jevelina/cdb/files"
	"github.com/hromov/jevelina/cdb/finance"
	"github.com/hromov/jevelina/cdb/leads"
	"github.com/hromov/jevelina/cdb/misc"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const TestDSN = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

var bucketName string

type CDB struct {
	*gorm.DB
}

//test part
func (db *CDB) Contacts() *contacts.Contacts {
	return &contacts.Contacts{DB: db.DB}
}

func (db *CDB) Leads() *leads.Leads {
	return &leads.Leads{DB: db.DB}
}

func (db *CDB) Misc() *misc.Misc {
	return &misc.Misc{DB: db.DB}
}

func (db *CDB) Finance() *finance.Finance {
	return &finance.Finance{DB: db.DB}
}

func (db *CDB) Files() *files.FilesService {
	return &files.FilesService{DB: db.DB, BucketName: bucketName}
}

func Init(dsn string, bucketName string) (*CDB, error) {
	if bucketName == "" {
		//TODO: shoud it return error? Check is it real?
		log.Println("FILES BUCKET NOT PROVIDED")
	}

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 30% fester but not so safe... let's give it a try
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect database error: %s", err.Error()))
	}

	// if table exist - do nothink, if not - create init structure with test data
	if !db.Migrator().HasTable("roles") {
		if err := db.AutoMigrate(&models.Role{}); err != nil {
			return nil, err
		}
	}
	if !db.Migrator().HasTable("contacts") {
		if err := db.AutoMigrate(&models.Contact{}); err != nil {
			return nil, err
		}
	}

	if !db.Migrator().HasTable("leads") {
		if err := db.AutoMigrate(&models.Lead{}); err != nil {
			return nil, err
		}
	}

	if !db.Migrator().HasTable("tasks") {
		if err := db.AutoMigrate(&models.Task{}); err != nil {
			return nil, err
		}
	}

	if !db.Migrator().HasTable("wallets") {
		if err := db.AutoMigrate(&models.Wallet{}); err != nil {
			return nil, err
		}
	}

	if !db.Migrator().HasTable("transfers") {
		if err := db.AutoMigrate(&models.Transfer{}); err != nil {
			return nil, err
		}
	}
	if !db.Migrator().HasTable("files") {
		if err := db.AutoMigrate(&models.File{}); err != nil {
			return nil, err
		}
	}
	return &CDB{DB: db}, nil
}
