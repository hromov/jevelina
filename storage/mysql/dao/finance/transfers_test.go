package finance

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type FilterTest struct {
	name    string
	filter  models.ListFilter
	queries []string
	args    []driver.Value
}

const timeForm = "Jan-02-2006"

var timeExample, _ = time.Parse(timeForm, "May-08-2022")
var filterTests = []FilterTest{
	{
		name:    "no filter",
		filter:  models.ListFilter{},
		queries: []string{regexp.QuoteMeta("SELECT * FROM `transfers` WHERE `transfers`.`deleted_at` IS NULL ORDER BY completed asc,completed_at desc,created_at desc")},
	},
	// {
	// 	name:   "date range",
	// 	filter: models.ListFilter{MinDate: timeExample, MaxDate: timeExample},
	// 	queries: []string{
	// 		"SELECT \\* FROM `transfers` WHERE .*",
	// 	},
	// },
	// {
	// 	name:    "single parent",
	// 	filter:  models.ListFilter{ParentID: 1000},
	// 	queries: []string{"SELECT \\* FROM `transfers` WHERE parent_id = ? .*"},
	// 	args:    []driver.Value{1000},
	// },
	// {
	// 	name:    "multiple parents",
	// 	filter:  models.ListFilter{IDs: []uint64{1000, 1001}},
	// 	queries: []string{"SELECT \\* FROM `transfers` WHERE \\(parent_id = 1000 OR parent_id = 1001\\) AND .*"},
	// },
}

var sumTests = []FilterTest{
	{
		name:   "sum, no filter",
		filter: models.ListFilter{},
		queries: []string{
			regexp.QuoteMeta("SELECT category, sum(amount) as total FROM `transfers` WHERE `from` IS NULL AND `transfers`.`deleted_at` IS NULL GROUP BY `category`"),
			regexp.QuoteMeta("SELECT category, sum(amount) as total FROM `transfers` WHERE `to` IS NULL AND `transfers`.`deleted_at` IS NULL GROUP BY `category`"),
		},
	},
	{
		name:   "sum, date range, completed",
		filter: models.ListFilter{MinDate: timeExample, MaxDate: timeExample, Completed: true},
		queries: []string{
			regexp.QuoteMeta(
				fmt.Sprintf("SELECT category, sum(amount) as total FROM `transfers` WHERE (completed_at >= '%s' AND completed_at < '%s') AND `from` IS NULL AND `transfers`.`deleted_at` IS NULL GROUP BY `category`", timeExample, timeExample),
			),
			regexp.QuoteMeta(
				fmt.Sprintf("SELECT category, sum(amount) as total FROM `transfers` WHERE (completed_at >= '%s' AND completed_at < '%s') AND `to` IS NULL AND `transfers`.`deleted_at` IS NULL GROUP BY `category`", timeExample, timeExample),
			)},
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

	countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
	filesRow := sqlmock.NewRows([]string{"id", "parent_id", "created_at", "updated_at", "name", "url"}).
		AddRow(1, 1, time.Now(), time.Now(), "some", "some")

	for _, test := range filterTests {
		t.Run(test.name, func(t *testing.T) {
			if test.args != nil {
				for _, query := range test.queries {
					s.mock.ExpectQuery(query).WithArgs(test.args...).WillReturnRows(rows)
				}

			} else {
				for _, query := range test.queries {
					s.mock.ExpectQuery(query).WillReturnRows(rows)
				}
			}
			s.mock.ExpectQuery("SELECT \\* FROM `files` WHERE .*").WillReturnRows(filesRow)
			s.mock.ExpectQuery("SELECT count\\(\\*\\) FROM `transfers` .*").WillReturnRows(countRow)

			_, err := s.finance.Transfers(test.filter)
			require.NoError(t, err)
			require.NoError(t, s.mock.ExpectationsWereMet())
		})
	}
}

func TestSumByCategory(t *testing.T) {
	s := &Suite{}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	columns := []string{"category", "total"}

	rows := sqlmock.NewRows(columns).
		AddRow("cat 1", 1000).
		AddRow("cat 2", 500)

	for _, test := range sumTests {
		t.Run(test.name, func(t *testing.T) {
			if test.args != nil {
				for _, query := range test.queries {
					s.mock.ExpectQuery(query).WithArgs(test.args...).WillReturnRows(rows)
				}
			} else {
				for _, query := range test.queries {
					s.mock.ExpectQuery(query).WillReturnRows(rows)
				}
			}

			_, err := s.finance.SumByCategory(test.filter)
			require.NoError(t, err)
			require.NoError(t, s.mock.ExpectationsWereMet())
		})
	}
}
