package finance

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
)

type TransferTest struct {
	name   string
	filter models.ListFilter
	query  string
}

func TestTransfers(t *testing.T) {
	s := &Suite{}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	columns := []string{"id", "created_at", "updated_at", "deleted_at", "description", "amount", "completed", "category"}

	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), time.Now(), gorm.DeletedAt{}, "t 1", 0, false, "cat 1").
		AddRow(2, time.Now(), time.Now(), gorm.DeletedAt{}, "t 2", 10000, true, "cat 2").
		AddRow(3, time.Now(), time.Now(), gorm.DeletedAt{}, "t 3", 5000, true, "cat 3").
		AddRow(4, time.Now(), time.Now(), gorm.DeletedAt{}, "t 4", 2000, true, "cat 3").
		AddRow(5, time.Now(), time.Now(), gorm.DeletedAt{}, "t 4", 2000, true, "cat 3")
	//TODO: is there a chance to check real values?
	count := sqlmock.NewRows([]string{"count"}).AddRow(150)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transfers` ORDER BY completed asc,completed_at desc,created_at desc")).WillReturnRows(rows)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `transfers`")).WillReturnRows(count)

	s.finance.Transfers(models.ListFilter{})
	// t.Logf("total = %d\n", tResponse.Total)
	// for i, transfer := range tResponse.Transfers {
	// 	t.Logf("%d = %+v\n", i, transfer)
	// }

	if err := s.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransfersWFilter(t *testing.T) {
	s := &Suite{}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	columns := []string{"id", "created_at", "updated_at", "deleted_at", "description", "amount", "completed", "category"}

	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), time.Now(), gorm.DeletedAt{}, "t 1", 0, false, "cat 1").
		AddRow(2, time.Now(), time.Now(), gorm.DeletedAt{}, "t 2", 10000, true, "cat 2").
		AddRow(3, time.Now(), time.Now(), gorm.DeletedAt{}, "t 3", 5000, true, "cat 3").
		AddRow(4, time.Now(), time.Now(), gorm.DeletedAt{}, "t 4", 2000, true, "cat 3").
		AddRow(5, time.Now(), time.Now(), gorm.DeletedAt{}, "t 4", 2000, true, "cat 3")
	count := sqlmock.NewRows([]string{"count"}).AddRow(3)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transfers` ORDER BY completed asc,completed_at desc,created_at desc LIMIT 3")).WillReturnRows(rows)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `transfers` LIMIT 3")).WillReturnRows(count)

	s.finance.Transfers(models.ListFilter{Limit: 3})
	// t.Logf("%+v, %s", transfers, err)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
