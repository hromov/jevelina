package finance

import (
	"database/sql/driver"
	"fmt"
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
	args   []driver.Value
}

const timeForm = "Jan-02-2006"

var timeExample, _ = time.Parse(timeForm, "May-08-2022")
var parentID = 1000
var tests = []TransferTest{
	{
		name:   "no filter",
		filter: models.ListFilter{},
		query:  regexp.QuoteMeta("SELECT * FROM `transfers` ORDER BY completed asc,completed_at desc,created_at desc"),
	},
	{
		name:   "limit 3",
		filter: models.ListFilter{Limit: 3},
		query:  regexp.QuoteMeta("SELECT * FROM `transfers` ORDER BY completed asc,completed_at desc,created_at desc LIMIT 3"),
	},
	{
		name:   "date range",
		filter: models.ListFilter{MinDate: timeExample, MaxDate: timeExample},
		query: regexp.QuoteMeta(
			fmt.Sprintf("SELECT * FROM `transfers` WHERE ((completed_at >= '%s' AND completed_at < '%s') OR completed_at IS NULL", timeExample, timeExample),
		),
	},
	{
		name:   "single parent",
		filter: models.ListFilter{ParentID: 1000},
		query:  regexp.QuoteMeta("SELECT * FROM `transfers` WHERE parent_id = ? ORDER BY completed asc,completed_at desc,created_at desc"),
		args:   []driver.Value{1000},
	},
	{
		name:   "multiple parents",
		filter: models.ListFilter{IDs: []uint64{1000, 1001}},
		query:  regexp.QuoteMeta("SELECT * FROM `transfers` WHERE parent_id = 1000 OR parent_id = 1001"),
	},
}

func TestTransfers(t *testing.T) {
	s := &Suite{}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	columns := []string{"id", "created_at", "updated_at", "deleted_at", "description", "amount", "completed", "category", "parent_id", "form", "to"}

	rows := sqlmock.NewRows(columns).
		AddRow(1, time.Now(), time.Now(), gorm.DeletedAt{}, "t 1", 0, false, "cat 1", 1000, 100, nil).
		AddRow(2, time.Now(), time.Now(), gorm.DeletedAt{}, "t 2", 10000, true, "cat 2", 1001, 100, nil).
		AddRow(3, time.Now(), time.Now(), gorm.DeletedAt{}, "t 3", 5000, true, "cat 3", nil, 100, 101).
		AddRow(4, time.Now(), time.Now(), gorm.DeletedAt{}, "t 4", 2000, true, "cat 3", 1002, nil, 101).
		AddRow(5, time.Now(), time.Now(), gorm.DeletedAt{}, "t 4", 2000, true, "cat 3", 1003, nil, 100)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args != nil {
				s.mock.ExpectQuery(test.query).WithArgs(test.args...).WillReturnRows(rows)
			} else {
				s.mock.ExpectQuery(test.query).WillReturnRows(rows)
			}
			//TODO: do I need fake couner?
			//s.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `transfers`")).WillReturnRows(count)

			s.finance.Transfers(test.filter)

			if err := s.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
