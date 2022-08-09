package models

import (
	"fmt"
	"time"

	"github.com/hromov/jevelina/domain/finances"
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

func ListFilterFromFin(f finances.Filter) ListFilter {
	return ListFilter{
		IDs:       f.IDs,
		Limit:     f.Limit,
		Offset:    f.Offset,
		Query:     f.Query,
		ParentID:  f.ParentID,
		MinDate:   f.MinDate,
		MaxDate:   f.MaxDate,
		From:      f.From,
		To:        f.To,
		Wallet:    f.Wallet,
		Completed: f.Completed,
	}
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
