package models

import (
	"time"

	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

type Manufacturer struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

func (p *Manufacturer) ToDomain() misc.Manufacturer {
	return misc.Manufacturer{
		ID:   p.ID,
		Name: p.Name,
	}
}

func ManufacturerFromDomain(p misc.Manufacturer) Manufacturer {
	return Manufacturer{
		ID:   p.ID,
		Name: p.Name,
	}
}
