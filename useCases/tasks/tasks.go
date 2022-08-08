package tasks

import (
	"time"

	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/users"
)

type Task struct {
	ID        uint64
	ParentID  uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	DeadLine  time.Time
	Completed bool

	Files       []misc.File
	Description string
	Results     string

	Responsible users.User
	Created     users.User
	Updated     users.User
}

type TasksResponse struct {
	Tasks []Task
	Total int64
}

type Filter struct {
	IDs           []uint64
	Limit         int
	Offset        int
	Query         string
	ParentID      uint64
	ResponsibleID uint64
	MinDate       time.Time
	MaxDate       time.Time
}

type TaskData struct {
	ID            uint64
	ParentID      uint64
	DeadLine      time.Time
	Description   string
	ResponsibleID uint64
	CreatedID     uint64
}
