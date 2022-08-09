package finances

import "time"

type Filter struct {
	IDs    []uint64
	Limit  int
	Offset int
	// LeadID        uint64
	// ContactID     uint64
	// TagID         uint8
	Query    string
	ParentID uint64
	// Active        bool
	// StepID        uint8
	// Steps         []uint8
	// ResponsibleID uint64
	MinDate   time.Time
	MaxDate   time.Time
	From      uint16
	To        uint16
	Wallet    uint16
	Completed bool
}
