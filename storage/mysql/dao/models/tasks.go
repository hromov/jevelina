package models

import (
	"time"

	"github.com/hromov/jevelina/useCases/tasks"
	"gorm.io/gorm"
)

//Task & Notice
type Task struct {
	ID        uint64 `gorm:"primaryKey"`
	ParentID  uint64 `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeadLine  *time.Time     `gorm:"index"`
	Completed bool

	//if not - notice
	TaskTypeID *uint8
	TaskType   TaskType

	//just links
	Files       []File `gorm:"foreignKey:ParentID"`
	Description string `gorm:"size:1024"`
	Results     string `gorm:"size:512"`

	ResponsibleID *uint64
	Responsible   User `gorm:"foreignKey:ResponsibleID"`
	CreatedID     *uint64
	Created       User `gorm:"foreignKey:CreatedID"`
	UpdatedID     *uint64
	Updated       User `gorm:"foreignKey:UpdatedID"`
}

type TasksResponse struct {
	Tasks []Task
	Total int64
}

type TaskType struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

func (t *Task) ToDomain() tasks.Task {
	return tasks.Task{
		ID:        t.ID,
		ParentID:  t.ParentID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt.Time,
		DeadLine:  Time(t.DeadLine),
		Completed: t.Completed,

		Files:       FilesToDomain(t.Files),
		Description: t.Description,
		Results:     t.Results,

		Responsible: t.Responsible.ToDomain(),
		Created:     t.Created.ToDomain(),
		Updated:     t.Updated.ToDomain(),
	}
}

func TaskFromTaskData(t tasks.TaskData) Task {
	return Task{
		ID:            t.ID,
		ParentID:      t.ParentID,
		DeadLine:      TimeOrNil(t.DeadLine),
		Description:   t.Description,
		ResponsibleID: OrNil(t.ResponsibleID),
		CreatedID:     OrNil(t.CreatedID),
		Results:       t.Results,
		Completed:     t.Completed,
	}
}
