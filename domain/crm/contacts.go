package crm

import (
	"time"

	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm"
)

type Contact struct {
	ID          uint64 `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
	Name        string
	SecondName  string
	Responsible users.User
	Created     users.User
	Phone       string
	SecondPhone string
	Email       string
	SecondEmail string
	URL         string
	City        string
	Address     string
	Source      Source
	Position    string
}
