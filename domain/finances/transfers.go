package finances

import (
	"time"

	"github.com/hromov/jevelina/domain/misc/files"
)

type Transfer struct {
	ID          uint64
	ParentID    uint64
	CreatedAt   time.Time
	CreatedBy   uint64
	UpdatedAt   time.Time
	DeletedAt   time.Time
	DeletedBy   uint64
	Completed   bool
	CompletedAt time.Time
	Description string
	CompletedBy uint64
	From        uint16
	To          uint16
	Category    string
	Amount      int64
	Files       []files.File
}

type TransfersResponse struct {
	Transfers []Transfer
	Total     int64
}
