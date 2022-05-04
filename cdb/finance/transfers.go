package finance

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (f *Finance) CreateTransfer(t *models.Transfer) (*models.Transfer, error) {
	t.Completed = false
	// t.CreatedAt = time.Now()
	if err := f.DB.Omit(clause.Associations).Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (f *Finance) UpdateTransfer(t *models.Transfer) error {
	var oldTransfer *models.Transfer
	f.DB.First(&oldTransfer, t.ID)
	if oldTransfer == nil {
		return errors.New(fmt.Sprintf("Can't find transfer with ID = %d", t.ID))
	}
	if oldTransfer.Completed {
		return errors.New("Can't change completed transfer")
	}
	return f.DB.Omit(clause.Associations).Save(t).Error
}

func (f *Finance) CompleteTransfer(ID uint64, userID uint64) error {

	return f.DB.Transaction(func(tx *gorm.DB) error {

		var t *models.Transfer
		tx.First(&t, ID)
		if t == nil {
			return errors.New(fmt.Sprintf("Can't find transfer with ID = %d", ID))
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

func (f *Finance) DeleteTransfer(ID uint64) error {
	return f.DB.Delete(&models.Transfer{ID: ID}).Error
}

func (f *Finance) Transfers(filter models.ListFilter) (*models.TransfersResponse, error) {
	log.Println(filter)
	cr := &models.TransfersResponse{}
	q := f.DB.Limit(filter.Limit).Offset(filter.Offset)
	// if IDs providen - return here and it has to be used as parent's ID, because we don't know transfers IDs other way
	if len(filter.IDs) > 0 {
		search := ""
		for i, step := range filter.IDs {
			search += fmt.Sprintf("parent_id = %d", step)
			if i < (len(filter.Steps) - 1) {
				search += " OR "
			}
		}
		if err := q.Where(search).Find(&cr.Transfers).Count(&cr.Total).Error; err != nil {
			return nil, err
		}
		return cr, nil
	}

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
		q = q.Where(f.DB.Where("`from` = ?", filter.Wallet).Or("`to` = ?", filter.Wallet))
	}
	if !filter.MinDate.IsZero() {
		q = q.Where("completed_at >= ?", filter.MinDate)
	}
	if !filter.MaxDate.IsZero() {
		q = q.Where("completed_at < ?", filter.MaxDate)
	}
	if !filter.MaxDate.IsZero() && !filter.MinDate.IsZero() {
		q = q.Or("completed_at IS NULL")
	}
	//TODO: check if it gives all uncompleted at first place
	q = q.Order("completed_at desc").Order("created_at asc")
	if result := q.Find(&cr.Transfers).Count(&cr.Total); result.Error != nil {
		return nil, result.Error
	}
	return cr, nil
}

//TODO: some grouped by category returns for analyze
