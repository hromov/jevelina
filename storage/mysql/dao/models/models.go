package models

import (
	"fmt"
	"time"

	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

type ListFilter struct {
	IDs           []uint64
	Limit         int
	Offset        int
	LeadID        uint64
	ContactID     uint64
	TagID         uint8
	Query         string
	ParentID      uint64
	Active        bool
	StepID        uint8
	Steps         []uint8
	ResponsibleID uint64
	MinDate       time.Time
	MaxDate       time.Time
	From          uint16
	To            uint16
	Wallet        uint16
	Completed     bool
}
type Analytics struct {
	CID string `gorm:"size:64"`
	UID string `gorm:"size:64"`
	TID string `gorm:"size:64"`

	UtmID       string `gorm:"size:64"`
	UtmSource   string `gorm:"size:64"`
	UtmMedium   string `gorm:"size:64"`
	UtmCampaign string `gorm:"size:64"`

	Domain string `gorm:"size:128"`
}

type Tag struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

type Manufacturer struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

func (m *Manufacturer) ToDomain() misc.Manufacturer {
	return misc.Manufacturer{
		ID:   m.ID,
		Name: m.Name,
	}
}

type Wallet struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32"`
	Balance   int64
	Closed    bool
}

type Transfer struct {
	ID uint64 `gorm:"primaryKey"`
	//Usualy LeadID
	ParentID  *uint64 `gorm:"index"`
	CreatedAt time.Time
	//UserID
	CreatedBy   uint64
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	DeletedBy   uint64
	Completed   bool
	CompletedAt *time.Time
	Description string
	//UserID
	CompletedBy uint64
	//Wallet
	From *uint16 `gorm:"index"`
	//Wallet
	To *uint16 `gorm:"index"`
	// Can be changed to id later, will try like this for now
	Category string
	Amount   int64
	Files    []File `gorm:"foreignKey:ParentID"`
}

type TransfersResponse struct {
	Transfers []Transfer
	Total     int64
}

type File struct {
	ID        uint64 `gorm:"primaryKey"`
	ParentID  uint64 `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:32"`
	URL       string `gorm:"size:128"`
}

type FileAddReq struct {
	Parent uint64
	Name   string
	Type   string
	Value  string
}

type EventParentType int16

const (
	TransferEvent EventParentType = iota + 1
	LeadEvent
	ContactEvent
)

type EventType int16

const (
	Create EventType = iota + 1
	Update
	Delete
	CategoryChange
)

func (et EventType) String() string {
	switch et {
	case Create:
		return "Create"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case CategoryChange:
		return "Category Change"
	}
	return "unknown"
}

type Event struct {
	ID              uint64 `gorm:"primaryKey"`
	CreatedAt       time.Time
	ParentID        uint64
	UserID          uint64
	EventParentType EventParentType
	Description     string `gorm:"size:512"`
}

type NewEvent struct {
	ParentID        uint64
	UserID          uint64
	Message         string
	EventType       EventType
	EventParentType EventParentType
}

type EventsResponse struct {
	Events []Event
	Total  int64
}

type EventFilter struct {
	ParentID        uint64
	UserID          uint64
	EventParentType EventParentType
	Limit           int
	Offset          int
}

// TODO: move to cdb
func (filter *ListFilter) DateCondition() string {
	dateSearh := ""
	if !filter.MinDate.IsZero() {
		dateSearh += fmt.Sprintf("completed_at >= '%s'", filter.MinDate)
	}
	if !filter.MaxDate.IsZero() {
		if dateSearh != "" {
			dateSearh += " AND "
		}
		dateSearh += fmt.Sprintf("completed_at < '%s'", filter.MaxDate)
	}
	if !filter.Completed && dateSearh != "" {
		dateSearh = fmt.Sprintf("((%s) OR completed_at IS NULL)", dateSearh)
	}
	return dateSearh
}
