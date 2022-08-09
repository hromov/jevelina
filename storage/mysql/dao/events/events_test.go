package events

// type Suite struct {
// 	sqlDB  *sql.DB
// 	gormDB *gorm.DB
// 	mock   sqlmock.Sqlmock
// 	es     *EventService
// }

// func (s *Suite) Init() (err error) {

// 	s.sqlDB, s.mock, err = sqlmock.New()
// 	if err != nil {
// 		return fmt.Errorf("Failed to open mock sql db, got error: %v", err)
// 	}

// 	if s.sqlDB == nil {
// 		return fmt.Errorf("mock db is null")
// 	}

// 	if s.mock == nil {
// 		return fmt.Errorf("sqlmock is null")
// 	}

// 	s.gormDB, err = gorm.Open(mysql.New(mysql.Config{
// 		Conn:                      s.sqlDB,
// 		SkipInitializeWithVersion: true,
// 	}), &gorm.Config{})

// 	if err != nil {
// 		return err
// 	}

// 	s.es = &EventService{DB: s.gormDB}

// 	return nil
// }

// func (s *Suite) Close() {
// 	s.sqlDB.Close()
// }

// type ListTest struct {
// 	name   string
// 	query  string
// 	args   []driver.Value
// 	filter models.EventFilter
// }

// var listTests = []ListTest{
// 	{
// 		name:   "no user",
// 		query:  regexp.QuoteMeta("SELECT * FROM `events` WHERE event_parent_type = ?"),
// 		args:   []driver.Value{models.TransferEvent},
// 		filter: models.EventFilter{EventParentType: models.TransferEvent},
// 	},
// 	{
// 		name:   "w user",
// 		query:  regexp.QuoteMeta("SELECT * FROM `events` WHERE event_parent_type = ? AND user_id = ?"),
// 		args:   []driver.Value{models.TransferEvent, 1},
// 		filter: models.EventFilter{EventParentType: models.TransferEvent, UserID: 1},
// 	},
// 	{
// 		name:   "by parent",
// 		query:  regexp.QuoteMeta("SELECT * FROM `events` WHERE parent_id = ?"),
// 		args:   []driver.Value{1},
// 		filter: models.EventFilter{ParentID: 1},
// 	},
// }

// func TestCreate(t *testing.T) {
// 	s := &Suite{}
// 	if err := s.Init(); err != nil {
// 		t.Fatal(err)
// 	}
// 	defer s.Close()

// 	s.mock.ExpectBegin()

// 	s.mock.ExpectExec(
// 		regexp.QuoteMeta("INSERT INTO `events` (`created_at`,`parent_id`,`user_id`,`event_parent_type`,`description`) VALUES (?,?,?,?,?)")).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	s.mock.ExpectCommit()

// 	if err := s.es.Save(models.NewEvent{ParentID: 1, UserID: 1, Message: "some message", EventType: models.Create, EventParentType: models.TransferEvent}); err != nil {
// 		t.Errorf("Failed to insert to gorm db, got error: %v", err)
// 		t.FailNow()
// 	}

// 	if err := s.mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("Failed to meet expectations, got error: %v", err)
// 	}
// }

// func TestList(t *testing.T) {
// 	s := &Suite{}
// 	if err := s.Init(); err != nil {
// 		t.Fatal(err)
// 	}
// 	defer s.Close()

// 	columns := []string{"id", "created_at", "parent_id", "user_id", "event_parent_type", "description"}

// 	rows := sqlmock.NewRows(columns).
// 		AddRow(1, time.Now(), 1, 1, 1, "Some text for 1").
// 		AddRow(2, time.Now(), 2, 1, 1, "Some text for 2")
// 	count := sqlmock.NewRows([]string{"count"}).AddRow(1)
// 	cQ := "SELECT count *"

// 	for _, test := range listTests {
// 		t.Run(test.name, func(t *testing.T) {
// 			s.mock.ExpectQuery(test.query).WithArgs(test.args...).WillReturnRows(rows)
// 			s.mock.ExpectQuery(cQ).WithArgs(test.args...).WillReturnRows(count)
// 			if _, err := s.es.List(test.filter); err != nil {
// 				require.NoError(t, err)
// 			}
// 			if err := s.mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("there were unfulfilled expectations: %s", err)
// 			}
// 		})
// 	}
// }
