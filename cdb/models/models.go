package models

import (
	"time"

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
}

type Lead struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time     `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:64"`
	Budget    uint32
	Profit    int32

	//implement
	ContactID *uint64
	Contact   Contact `gorm:"foreignKey:ContactID"`

	ResponsibleID *uint64
	Responsible   User `gorm:"foreignKey:ResponsibleID"`
	CreatedID     *uint64
	Created       User `gorm:"foreignKey:CreatedID"`
	StepID        *uint8
	Step          Step
	//implement
	ProductID *uint32
	Product   Product
	//implement
	ManufacturerID *uint16
	Manufacturer   Manufacturer
	SourceID       *uint8
	Source         Source
	//google analytics
	Tags []Tag `gorm:"many2many:leads_tags;"`
	// Tasks []Task

	Analytics Analytics `gorm:"embedded;embeddedPrefix:analytics_"`
}

type LeadsResponse struct {
	Leads []Lead
	Total int64
}

type Contact struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//or company
	IsPerson   bool
	Name       string `gorm:"size:32"`
	SecondName string `gorm:"size:32"`
	//implement
	ResponsibleID *uint64
	Responsible   User `gorm:"foreignKey:ResponsibleID"`
	CreatedID     *uint64
	Created       User `gorm:"foreignKey:CreatedID"`

	Tags []Tag `gorm:"many2many:contacts_tags;"`
	// Tasks       []Task
	Phone       string `gorm:"size:32"`
	SecondPhone string `gorm:"size:32"`
	Email       string `gorm:"size:128"`
	SecondEmail string `gorm:"size:128"`
	URL         string `gorm:"size:128"`

	City    string `gorm:"size:128"`
	Address string `gorm:"size:256"`

	SourceID *uint8
	Source   Source `gorm:"foreignKey:SourceID"`
	Position string `gorm:"size:128"`

	Analytics Analytics `gorm:"embedded;embeddedPrefix:analytics_"`
}

type ContactsResponse struct {
	Contacts []Contact
	Total    int64
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

type User struct {
	ID           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"size:32"`
	Email        string         `gorm:"size:128; unique"`
	Hash         string         `gorm:"size:128; unique"`
	Distribution float32        `gorm:"type:decimal(2,2);"`
	// Events    []Event
	RoleID *uint8
	Role   Role `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Role struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Priority  uint8
	Role      string `gorm:"unique;size:32"`
}

type Step struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"unique;size:32"`
	//1st, 2nd etc
	Order  uint8
	Active bool `gorm:"index"`
}

type Event struct {
	ID          uint64 `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ParentID    uint64
	Description string `gorm:"size:256"`
}

type Tag struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

//Task & Notice
type Task struct {
	ID        uint64 `gorm:"primaryKey"`
	ParentID  uint64 `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeadLine  *time.Time     `gorm:"index"`
	Completed bool

	//if not - notice
	TaskTypeID *uint8
	TaskType   TaskType

	//just links
	Files       []File `gorm:"foreignKey:ParentID"`
	Description string `gorm:"size:1024"`
	Results     string `gorm:"size:512"`

	ResponsibleID *uint64
	Responsible   User `gorm:"foreignKey:ResponsibleID"`
	CreatedID     *uint64
	Created       User `gorm:"foreignKey:CreatedID"`
	UpdatedID     *uint64
	Updated       User `gorm:"foreignKey:UpdatedID"`
}

type TasksResponse struct {
	Tasks []Task
	Total int64
}

type TaskType struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

type Source struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

type Product struct {
	ID        uint32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:64;unique"`
}

type Manufacturer struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

type CreateLeadReq struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description,omitempty"`

	ClientName  string `json:"clientname"`
	ClientEmail string `json:"clientemail,omitempty"`
	ClientPhone string `json:"clientphone,omitempty"`

	Source       string `json:"source,omitempty"`
	Product      string `json:"product,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`

	UserEmail string `json:"user_email,omitempty"`
	UserHash  string `json:"user_hash,omitempty"`

	CID string `gorm:"size:64"`
	UID string `gorm:"size:64"`
	TID string `gorm:"size:64"`

	UtmID       string `gorm:"size:64"`
	UtmSource   string `gorm:"size:64"`
	UtmMedium   string `gorm:"size:64"`
	UtmCampaign string `gorm:"size:64"`

	Domain string `gorm:"size:128"`
}

type Wallet struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
	Balance   int64
	Closed    bool
}

type Transfer struct {
	ID uint64 `gorm:"primaryKey"`
	//Usualy LeadID
	ParentID  uint64 `gorm:"index"`
	CreatedAt time.Time
	//UserID
	CreatedBy   uint64
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Completed   bool
	CompletedAt time.Time
	//UserID
	CompletedBy uint64
	//Wallet
	From uint16 `gorm:"index"`
	//Wallet
	To uint16 `gorm:"index"`
	// Can be changed to id later, will try like this for now
	Category string
	Amount   int64
	Files    []File `gorm:"foreignKey:ParentID"`
}

type File struct {
	ID        uint64 `gorm:"primaryKey"`
	ParentID  uint64 `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
}

// From      *datastore.Key   `json:"from"`
// To        *datastore.Key   `json:"to"`
// Category  string           `json:"category"`
// Amount    float32          `json:"amount" datastore:",noindex"`
// CreatedAt time.Time        `json:"created_at,string"`
// CreatedBy string           `json:"created_by"`
// Status    TransferStatus   `json:"status"`
// Wallets   []*datastore.Key `json:"wallets"`
// Lead      string           `json:"lead"`
// Ancestor  *datastore.Key   `json:"ancestor" datastore:"-"`
// Files     []File           `json:"files,omitempty" datastore:",noindex"`
