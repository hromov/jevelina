package models

import (
	"time"

	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

type Source struct {
	ID        uint8 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32;unique"`
}

func (s *Source) toDomain() misc.Source {
	return misc.Source{
		ID:   s.ID,
		Name: s.Name,
	}
}
