package fin_api

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/finance"
)

func TestCategories(t *testing.T) {
	mock, err := cdb.InitMock()
	if err != nil {
		t.Fatalf("Can't init mock error: %s", err.Error())
	}

	// // create app with mocked db, request and response to test
	// app := &api{db}
	req, err := http.NewRequest("GET", "http://localhost:8080/categories", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"category"}).
		AddRow("cat 1").
		AddRow("cat 2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT(category) FROM transfers WHERE deleted_at IS NULL ORDER BY category asc")).WillReturnRows(rows)

	// now we execute our request
	CategoriesHandler(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	data := []string{
		"cat 1",
		"cat 2",
	}
	cdb.AssertJSON(w.Body.Bytes(), data, t)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCategoriesSum(t *testing.T) {
	mock, err := cdb.InitMock()
	if err != nil {
		t.Fatalf("Can't init mock error: %s", err.Error())
	}

	// // create app with mocked db, request and response to test
	// app := &api{db}
	req, err := http.NewRequest("GET", "http://localhost:8080/analytics/categories", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	columns := []string{"category", "total"}

	rows := sqlmock.NewRows(columns).
		AddRow("cat 1", 1000).
		AddRow("cat 2", 500)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT category, sum(amount) as total FROM `transfers` WHERE ")).WillReturnRows(rows)

	// now we execute our request
	CategoriesSumHandler(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	data := []finance.SumResult{
		{Category: "cat 1", Total: 1000},
		{Category: "cat 2", Total: 500},
	}
	cdb.AssertJSON(w.Body.Bytes(), data, t)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
