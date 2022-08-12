package models

import (
	"time"

	"github.com/hromov/jevelina/domain/finances"
	"gorm.io/gorm"
)

type Transfer struct {
	ID uint64 `gorm:"primaryKey"`
	//Usualy LeadID
	ParentID  *uint64 `gorm:"index"`
	CreatedAt time.Time
	//UserID
	CreatedBy   uint64
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	DeletedBy   uint64
	Completed   bool
	CompletedAt *time.Time
	Description string
	//UserID
	CompletedBy uint64
	//Wallet
	From *uint16 `gorm:"index"`
	//Wallet
	To *uint16 `gorm:"index"`
	// Can be changed to id later, will try like this for now
	Category string
	Amount   int64
	Files    []File `gorm:"foreignKey:ParentID"`
}

func TransferFromDomain(t finances.Transfer) Transfer {
	return Transfer{
		ID:          t.ID,
		ParentID:    OrNil64(t.ParentID),
		CreatedAt:   t.CreatedAt,
		CreatedBy:   t.CreatedBy,
		Completed:   t.Completed,
		CompletedAt: TimeOrNil(t.CompletedAt),
		Description: t.Description,
		CompletedBy: t.CompletedBy,
		From:        OrNil16(t.From),
		To:          OrNil16(t.To),
		Category:    t.Category,
		Amount:      t.Amount,
	}
}

func (t *Transfer) ToDomain() finances.Transfer {
	return finances.Transfer{
		ID:          t.ID,
		ParentID:    Val64(t.ParentID),
		CreatedAt:   t.CreatedAt,
		CreatedBy:   t.CreatedBy,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt.Time,
		DeletedBy:   t.DeletedBy,
		Completed:   t.Completed,
		CompletedAt: Time(t.CompletedAt),
		Description: t.Description,
		CompletedBy: t.CompletedBy,
		From:        Val16(t.From),
		To:          Val16(t.To),
		Category:    t.Category,
		Amount:      t.Amount,
		Files:       FilesToDomain(t.Files),
	}
}

type TransfersResponse struct {
	Transfers []Transfer
	Total     int64
}
