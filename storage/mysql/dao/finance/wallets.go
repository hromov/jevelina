package finance

import (
	"context"
	"fmt"
	"log"

	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
)

type Finance struct {
	db *gorm.DB
}

func NewFinance(db *gorm.DB, automigrate bool) *Finance {
	if automigrate {
		if err := db.AutoMigrate(&models.Wallet{}); err != nil {
			log.Printf("migration for %s error: %s\n", "wallet", err.Error())
		}
		if err := db.AutoMigrate(&models.Transfer{}); err != nil {
			log.Printf("migration for %s error: %s\n", "transfer", err.Error())
		}
	}
	return &Finance{db}
}

func (f *Finance) CreateWallet(ctx context.Context, w finances.Wallet) (finances.Wallet, error) {
	dbWallet := models.WalletFromDomain(w)
	if err := f.db.WithContext(ctx).Create(&dbWallet).Error; err != nil {
		return finances.Wallet{}, err
	}
	return dbWallet.ToDomain(), nil
}

func (f *Finance) ChangeWalletName(ctx context.Context, id uint16, name string) error {
	wallet := models.Wallet{ID: id}
	if err := f.db.WithContext(ctx).First(&wallet).Error; err != nil {
		return fmt.Errorf("Can't find wallet with ID = %d. Error: %s", id, err.Error())
	}
	wallet.Name = name
	return f.db.WithContext(ctx).Save(wallet).Error
}

func (f *Finance) ChangeWalletState(ctx context.Context, id uint16, closed bool) error {
	wallet := models.Wallet{ID: id}
	if err := f.db.WithContext(ctx).First(&wallet).Error; err != nil {
		return fmt.Errorf("Can't find wallet with ID = %d. Error: %s", id, err.Error())
	}
	wallet.Closed = closed
	return f.db.WithContext(ctx).Save(&wallet).Error
}

func (f *Finance) DeleteWallet(ctx context.Context, id uint16) error {
	return f.db.WithContext(ctx).Delete(&models.Wallet{ID: id}).Error
}

func (f *Finance) ListWallets(ctx context.Context) ([]finances.Wallet, error) {
	dbWallets := []models.Wallet{}
	if err := f.db.WithContext(ctx).Find(&dbWallets).Error; err != nil {
		return nil, err
	}
	wallets := make([]finances.Wallet, len(dbWallets))
	for i, w := range dbWallets {
		wallets[i] = w.ToDomain()
	}
	return wallets, nil
}
