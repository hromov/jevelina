package crm

import (
	"time"

	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm"
)

type Lead struct {
	ID          uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    *time.Time
	DeletedAt   gorm.DeletedAt
	Name        string
	Budget      uint32
	Profit      int32
	Contact     Contact
	Responsible users.User
	Created     users.User
	Step        Step
	Source      Source
}

type Step struct {
	ID        uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Name      string
	Order     uint8
	Active    bool
}
