package models

import (
	"time"

	"github.com/hromov/jevelina/utils/events"
)

type Event struct {
	ID              uint64 `gorm:"primaryKey"`
	CreatedAt       time.Time
	ParentID        uint64
	UserID          uint64
	EventParentType events.EventParentType
	Description     string `gorm:"size:512"`
}

func (e *Event) ToDomain() events.Event {
	return events.Event(*e)
}
