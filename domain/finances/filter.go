package finances

import "time"

type Filter struct {
	IDs       []uint64
	Limit     int
	Offset    int
	Query     string
	ParentID  uint64
	MinDate   time.Time
	MaxDate   time.Time
	From      uint16
	To        uint16
	Wallet    uint16
	Completed bool
}
