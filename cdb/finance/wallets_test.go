package finance

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	sqlDB   *sql.DB
	gormDB  *gorm.DB
	mock    sqlmock.Sqlmock
	finance *Finance
}

func (s *Suite) Init() (err error) {

	s.sqlDB, s.mock, err = sqlmock.New()
	if err != nil {
		return fmt.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	if s.sqlDB == nil {
		return fmt.Errorf("mock db is null")
	}

	if s.mock == nil {
		return fmt.Errorf("sqlmock is null")
	}

	s.gormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      s.sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Suite) Close() {
	s.sqlDB.Close()
}

func TestCreateWallet(t *testing.T) {
	s := &Suite{}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	initWallet := &models.Wallet{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
		Name:      "test wallet",
		Balance:   0,
		Closed:    false,
		ID:        123,
	}

	walletID := initWallet.ID

	s.finance = &Finance{DB: s.gormDB}

	s.mock.ExpectBegin()

	s.mock.ExpectExec(
		regexp.QuoteMeta("INSERT INTO `wallets` (`created_at`,`updated_at`,`deleted_at`,`name`,`balance`,`closed`,`id`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(initWallet.CreatedAt, initWallet.UpdatedAt, initWallet.DeletedAt, initWallet.Name, initWallet.Balance, initWallet.Closed, initWallet.ID).
		WillReturnResult(sqlmock.NewResult(int64(walletID), 1))

	s.mock.ExpectCommit()

	if _, err := s.finance.CreateWallet(initWallet); err != nil {
		t.Errorf("Failed to insert to gorm db, got error: %v", err)
		t.FailNow()
	}

	if err := s.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}
