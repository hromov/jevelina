package leads

import (
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/users"
)

type Filter struct {
	IDs    []uint64
	Limit  int
	Offset int
	// LeadID        uint64
	ContactID uint64
	// TagID         uint8
	Query string
	// ParentID      uint64
	Active        bool
	StepID        uint8
	Steps         []uint8
	ResponsibleID uint64
	MinDate       time.Time
	MaxDate       time.Time
	// From          uint16
	// To            uint16
	// Wallet        uint16
	Completed bool
}

type Lead struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  time.Time
	DeletedAt time.Time
	Name      string
	Budget    uint32
	Profit    int32

	Contact     contacts.Contact
	Responsible users.User
	Created     users.User
	Step        Step

	Product      misc.Product
	Manufacturer misc.Manufacturer
	Source       misc.Source
	Analytics    misc.Analytics
}

type LeadsResponse struct {
	Leads []Lead
	Total int64
}

type Step struct {
	ID     uint8
	Name   string
	Order  uint8
	Active bool
}

type LeadRequest struct {
	Name        string
	Price       int
	Description string

	ClientName  string
	ClientEmail string
	ClientPhone string

	Source       string
	Product      string
	Manufacturer string

	UserEmail string
	UserHash  string

	CID string
	UID string
	TID string

	UtmID       string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string

	Domain string
}

type LeadData struct {
	ID       uint64
	ClosedAt time.Time
	Name     string
	Budget   uint32
	Profit   int32

	ContactID      uint64
	ResponsibleID  uint64
	CreatedID      uint64
	StepID         uint8
	ProductID      uint32
	ManufacturerID uint16
	SourceID       uint8
	//TODO: move it
	Analytics misc.Analytics
}
