package finance

import (
	"fmt"

	"github.com/hromov/jevelina/storage/mysql/dao/events"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Finance struct {
	*gorm.DB
	Events *events.EventService
}

func (f *Finance) CreateWallet(item *models.Wallet) (*models.Wallet, error) {
	if err := f.DB.Omit(clause.Associations).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (f *Finance) ChangeWalletName(ID uint16, name string) error {
	var wallet *models.Wallet
	f.DB.First(&wallet, ID)
	if wallet == nil {
		return fmt.Errorf("Can't find wallet with ID = %d", ID)
	}
	wallet.Name = name
	return f.DB.Omit(clause.Associations).Save(wallet).Error
}

func (f *Finance) ChangeWalletState(ID uint16, closed bool) error {
	var wallet *models.Wallet
	f.DB.First(&wallet, ID)
	if wallet == nil {
		return fmt.Errorf("Can't find wallet with ID = %d", ID)
	}
	wallet.Closed = closed
	return f.DB.Omit(clause.Associations).Save(wallet).Error
}

func (f *Finance) DeleteWallet(ID uint16) error {
	return f.DB.Delete(&models.Wallet{ID: ID}).Error
}

func (f *Finance) ListWallets(filter *models.ListFilter) (items []models.Wallet, err error) {
	if result := f.DB.Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return
}
