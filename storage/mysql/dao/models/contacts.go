package models

import (
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

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

func (c *Contact) ToDomain() contacts.Contact {
	return contacts.Contact{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt.Time,

		Name:        c.Name,
		SecondName:  c.SecondName,
		Responsible: c.Responsible.ToDomain(),
		Created:     c.Responsible.ToDomain(),
		Phone:       c.Phone,
		SecondPhone: c.SecondPhone,
		Email:       c.Email,
		SecondEmail: c.SecondEmail,
		URL:         c.URL,

		City:    c.City,
		Address: c.Address,

		Source:   c.Source.ToDomain(),
		Position: c.Position,

		Analytics: misc.Analytics(c.Analytics),
	}
}

func ContactFromDomain(c contacts.ContactRequest) Contact {
	contact := Contact{
		IsPerson:      true,
		Name:          c.Name,
		SecondName:    c.SecondName,
		ResponsibleID: OrNil(c.ResponsibleID),
		CreatedID:     OrNil(c.CreatedID),

		Phone:       c.Phone,
		SecondPhone: c.SecondPhone,
		Email:       c.Email,
		SecondEmail: c.SecondEmail,
		URL:         c.URL,

		City:    c.City,
		Address: c.Address,

		Position: c.Position,
		SourceID: OrNil(c.SourceID),

		Analytics: Analytics(c.Analytics),
	}
	return contact
}
