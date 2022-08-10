package events

import "time"

type Event struct {
	ID              uint64
	CreatedAt       time.Time
	ParentID        uint64
	UserID          uint64
	EventParentType EventParentType
	Description     string
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
