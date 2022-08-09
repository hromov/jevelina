package models

import (
	"time"

	"github.com/hromov/jevelina/domain/finances"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uint16 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:32"`
	Balance   int64
	Closed    bool
}

func WalletFromDomain(w finances.Wallet) Wallet {
	return Wallet{
		ID:      w.ID,
		Name:    w.Name,
		Balance: w.Balance,
		Closed:  w.Closed,
	}
}

func (w *Wallet) ToDomain() finances.Wallet {
	return finances.Wallet{
		ID:      w.ID,
		Name:    w.Name,
		Balance: w.Balance,
		Closed:  w.Closed,
	}
}
