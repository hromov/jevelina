package finance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (f *Finance) CreateTransfer(ctx context.Context, t finances.Transfer) (finances.Transfer, error) {
	transfer := models.TransferFromDomain(t)
	if err := f.db.WithContext(ctx).Create(&transfer).Error; err != nil {
		return finances.Transfer{}, err
	}
	return transfer.ToDomain(), nil
}

func (f *Finance) UpdateTransfer(ctx context.Context, t finances.Transfer) error {
	transfer := models.TransferFromDomain(t)
	return f.db.WithContext(ctx).Omit(clause.Associations).Model(&models.Transfer{}).Where("id", t.ID).Updates(&transfer).Error
}

func (f *Finance) CompleteTransfer(ctx context.Context, id uint64, userID uint64) error {

	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var t *models.Transfer
		tx.First(&t, id)
		if t == nil {
			return fmt.Errorf("Can't find transfer with ID = %d", id)
		}

		if !t.DeletedAt.Time.IsZero() {
			return errors.New("Can't complete deleted target")
		}

		if t.Completed {
			return errors.New("Transfer already completed")
		}

		if t.From != nil {
			var from *models.Wallet
			if err := tx.First(&from, t.From).Error; err != nil {
				return err
			}
			from.Balance -= t.Amount
			if err := tx.Save(from).Error; err != nil {
				return err
			}
		}

		if t.To != nil {
			var to *models.Wallet
			if err := tx.First(&to, t.To).Error; err != nil {
				return err
			}
			to.Balance += t.Amount
			if err := tx.Save(to).Error; err != nil {
				return err
			}
		}

		t.Completed = true
		t.CompletedBy = userID
		n := time.Now()
		t.CompletedAt = &n
		return tx.Save(t).Error
	})
}

func (f *Finance) DeleteTransfer(ctx context.Context, id uint64, userID uint64) error {
	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var t *models.Transfer
		tx.First(&t, id)
		if t == nil {
			return fmt.Errorf("Can't find transfer with ID = %d", id)
		}
		if t.Completed {
			if t.From != nil {
				var from *models.Wallet
				if err := tx.First(&from, t.From).Error; err != nil {
					return err
				}
				from.Balance += t.Amount
				if err := tx.Save(from).Error; err != nil {
					return err
				}
			}

			if t.To != nil {
				var to *models.Wallet
				if err := tx.First(&to, t.To).Error; err != nil {
					return err
				}
				to.Balance -= t.Amount
				if err := tx.Save(to).Error; err != nil {
					return err
				}
			}
			t.DeletedBy = userID
			if err := tx.Save(t).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&models.Transfer{ID: id}).Error
	})
}

func (f *Finance) Transfers(ctx context.Context, ff finances.Filter) (finances.TransfersResponse, error) {
	tr := finances.TransfersResponse{}
	dbTransfers := []models.Transfer{}
	filter := models.ListFilterFromFin(ff)
	q := f.db.WithContext(ctx).Preload("Files").Limit(filter.Limit).Offset(filter.Offset)
	// if IDs providen - return here and it has to be used as parent's ID, because we don't know transfers IDs other way
	if len(filter.IDs) > 0 {
		search := ""
		for i, step := range filter.IDs {
			search += fmt.Sprintf("parent_id = %d", step)
			if i < (len(filter.IDs) - 1) {
				search += " OR "
			}
		}
		if err := q.Where(search).Find(&dbTransfers).Count(&tr.Total).Error; err != nil {
			return finances.TransfersResponse{}, err
		}
	} else {
		//Category is text field for now, so let's use query. TODO: if changed to ID...
		if filter.Query != "" {
			q = q.Where("category LIKE ?", "%"+filter.Query+"%")
		}
		if filter.ParentID != 0 {
			q = q.Where("parent_id = ?", filter.ParentID)
		}
		if filter.From != 0 {
			q = q.Where("from = ?", filter.From)
		}
		if filter.To != 0 {
			q = q.Where("to = ?", filter.To)
		}
		if filter.Wallet != 0 {
			q = q.Where(f.db.Where("`from` = ?", filter.Wallet).Or("`to` = ?", filter.Wallet))
		}
		q = q.Where(filter.DateCondition())
		//TODO: check if it gives all uncompleted at first place
		q.Order("completed asc").Order("completed_at desc").Order("created_at desc")
		if result := q.Find(&dbTransfers).Count(&tr.Total); result.Error != nil {
			return finances.TransfersResponse{}, result.Error
		}
	}
	tr.Transfers = make([]finances.Transfer, len(dbTransfers))
	for i, t := range dbTransfers {
		tr.Transfers[i] = t.ToDomain()
	}
	return tr, nil
}

func (f *Finance) TransferCategories(ctx context.Context) ([]string, error) {
	categories := []string{}
	err := f.db.WithContext(ctx).Raw("SELECT DISTINCT(category) FROM transfers WHERE deleted_at IS NULL ORDER BY category asc").Scan(&categories).Error
	return categories, err
}

func (f *Finance) SumByCategory(ctx context.Context, ff finances.Filter) (finances.CategorisedCashflow, error) {
	filter := models.ListFilterFromFin(ff)
	incomes := []finances.CatTotal{}
	expenses := []finances.CatTotal{}
	q := f.db.WithContext(ctx).Model(&models.Transfer{}).Where(filter.DateCondition())
	if err := q.Select("category, sum(amount) as total").Where("`from` IS NULL").Group("category").Find(&incomes).Error; err != nil {
		return finances.CategorisedCashflow{}, fmt.Errorf("Can't get incomes error: %s", err.Error())
	}
	q2 := f.db.WithContext(ctx).Model(&models.Transfer{}).Where(filter.DateCondition())
	if err := q2.Select("category, sum(amount) as total").Where("`to` IS NULL").Group("category").Find(&expenses).Error; err != nil {
		return finances.CategorisedCashflow{}, fmt.Errorf("Can't get expenses error: %s", err.Error())
	}
	return finances.CategorisedCashflow{Incomes: incomes, Expenses: expenses}, nil
}

func (f *Finance) GetTransfer(ctx context.Context, id uint64) (finances.Transfer, error) {
	transfer := models.Transfer{ID: id}
	if err := f.db.WithContext(ctx).First(&transfer).Error; err != nil {
		return finances.Transfer{}, nil
	}
	return transfer.ToDomain(), nil
}
