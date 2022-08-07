package models

import (
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

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

type Step struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"unique;size:32"`
	Order     uint8
	Active    bool `gorm:"index"`
}

func (s Step) ToDomain() leads.Step {
	return leads.Step{
		ID:     s.ID,
		Name:   s.Name,
		Order:  s.Order,
		Active: s.Active,
	}
}

func StepFromDomain(s leads.Step) Step {
	return Step{
		ID:     s.ID,
		Name:   s.Name,
		Order:  s.Order,
		Active: s.Active,
	}
}

func LeadFromDomain(l leads.LeadData) Lead {

	return Lead{
		ID:     l.ID,
		Name:   l.Name,
		Budget: l.Budget,
		Profit: l.Profit,

		ContactID:      OrNil(l.ContactID),
		ResponsibleID:  OrNil(l.ResponsibleID),
		CreatedID:      OrNil(l.CreatedID),
		StepID:         OrNil(l.StepID),
		ProductID:      OrNil(l.ProductID),
		ManufacturerID: OrNil(l.ManufacturerID),
		SourceID:       OrNil(l.SourceID),
		ClosedAt:       TimeOrNil(l.ClosedAt),

		Analytics: Analytics(l.Analytics),
	}
}

func LeadFromFullDomain(l leads.Lead) Lead {

	return Lead{
		ID:     l.ID,
		Name:   l.Name,
		Budget: l.Budget,
		Profit: l.Profit,

		ContactID:      OrNil(l.Contact.ID),
		ResponsibleID:  OrNil(l.Responsible.ID),
		CreatedID:      OrNil(l.Created.ID),
		StepID:         OrNil(l.Step.ID),
		ProductID:      OrNil(l.Product.ID),
		ManufacturerID: OrNil(l.Manufacturer.ID),
		SourceID:       OrNil(l.Source.ID),
		ClosedAt:       TimeOrNil(l.ClosedAt),

		Analytics: Analytics(l.Analytics),
	}
}

func LeadFromRequest(lr leads.LeadRequest, contact contacts.Contact) Lead {
	lead := Lead{
		Name:          lr.Name,
		Budget:        uint32(lr.Price),
		ResponsibleID: &contact.Responsible.ID,
		ContactID:     &contact.ID,
	}
	item := lead
	item.Analytics.CID = lr.CID
	uid := ""
	if lr.UID == "" {
		uid = contact.Analytics.UID
	}
	item.Analytics.UID = uid
	item.Analytics.TID = lr.TID
	item.Analytics.UtmID = lr.UtmID
	item.Analytics.UtmSource = lr.UtmSource
	item.Analytics.UtmMedium = lr.UtmMedium
	item.Analytics.UtmCampaign = lr.UtmCampaign
	item.Analytics.Domain = lr.Domain

	return lead
}

func (l *Lead) ToDomain() leads.Lead {
	closedAt := time.Time{}
	if l.ClosedAt != nil {
		closedAt = *l.ClosedAt
	}
	return leads.Lead{
		ID:        l.ID,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
		ClosedAt:  closedAt,
		DeletedAt: l.DeletedAt.Time,
		Name:      l.Name,
		Budget:    l.Budget,
		Profit:    l.Profit,

		Contact:     l.Contact.ToDomain(),
		Responsible: l.Responsible.ToDomain(),
		Created:     l.Created.ToDomain(),
		Step:        l.Step.ToDomain(),

		Product:      l.Product.ToDomain(),
		Manufacturer: l.Manufacturer.ToDomain(),
		Source:       l.Source.ToDomain(),
		Analytics:    misc.Analytics(l.Analytics),
	}
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
