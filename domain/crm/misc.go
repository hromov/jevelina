package crm

import (
	"time"
)

type Source struct {
	ID        uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Name      string
}
