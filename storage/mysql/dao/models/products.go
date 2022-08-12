package models

import (
	"time"

	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:64;unique"`
}

func (p *Product) ToDomain() misc.Product {
	return misc.Product{
		ID:   p.ID,
		Name: p.Name,
	}
}

func ProductFromDomain(p misc.Product) Product {
	return Product{
		ID:   p.ID,
		Name: p.Name,
	}
}
