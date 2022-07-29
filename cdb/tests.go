package cdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenTest() (*CDB, error) {
	db, err := gorm.Open(mysql.Open(testDSN), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database error: %s", err.Error())
	}
	return &CDB{DB: db}, nil
}

func OpenAndInitTest() (*CDB, error) {
	db, err := OpenTest()
	if err != nil {
		return nil, err
	}
	return db, db.Init()
}

func InitMock() (sqlmock.Sqlmock, error) {

	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, fmt.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	if db == nil {
		return nil, fmt.Errorf("mock db is null")
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	currentDB = &CDB{DB: gormDB}

	return mock, nil
}

func AssertJSON(actual []byte, data interface{}, t *testing.T) {
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}

	if !bytes.Equal(expected, actual) {
		t.Errorf("the expected json: %s is different from actual %s", expected, actual)
	}
}
