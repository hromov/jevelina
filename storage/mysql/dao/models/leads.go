package models

import (
	"time"

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

		ContactID:      OrNil64(l.ContactID),
		ResponsibleID:  OrNil64(l.ResponsibleID),
		CreatedID:      OrNil64(l.CreatedID),
		StepID:         OrNil8(l.StepID),
		ProductID:      OrNil32(l.ProductID),
		ManufacturerID: OrNil16(l.ManufacturerID),
		SourceID:       OrNil8(l.SourceID),
		ClosedAt:       TimeOrNil(l.ClosedAt),

		Analytics: Analytics(l.Analytics),
	}
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
