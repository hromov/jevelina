package finance

import (
	"errors"
	"fmt"
	"time"

	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (f *Finance) CreateTransfer(t *models.Transfer) (*models.Transfer, error) {
	t.Completed = false
	if t.From == t.To {
		return nil, errors.New("Can't transfer to the same wallet")
	}
	// t.CreatedAt = time.Now()
	if err := f.DB.Omit(clause.Associations).Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (f *Finance) UpdateTransfer(t *models.Transfer) error {
	var oldTransfer *models.Transfer
	if t.From == t.To {
		return errors.New("Can't transfer to the same wallet")
	}
	f.DB.Unscoped().First(&oldTransfer, t.ID)
	if oldTransfer == nil {
		return fmt.Errorf("Can't find transfer with ID = %d", t.ID)
	}
	if oldTransfer.Completed || !oldTransfer.DeletedAt.Time.IsZero() {
		return f.DB.Unscoped().Model(oldTransfer).Updates(models.Transfer{Category: t.Category, Description: t.Description}).Error
	}
	return f.DB.Omit(clause.Associations).Save(t).Error
}

func (f *Finance) CompleteTransfer(ID uint64, userID uint64) error {

	return f.DB.Transaction(func(tx *gorm.DB) error {

		var t *models.Transfer
		tx.First(&t, ID)
		if t == nil {
			return fmt.Errorf("Can't find transfer with ID = %d", ID)
		}

		if !t.DeletedAt.Time.IsZero() {
			return errors.New("Can't complete deleted target")
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

func (f *Finance) DeleteTransfer(ID uint64, userID uint64) error {
	return f.DB.Transaction(func(tx *gorm.DB) error {

		var t *models.Transfer
		tx.First(&t, ID)
		if t == nil {
			return fmt.Errorf("Can't find transfer with ID = %d", ID)
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
		return tx.Delete(&models.Transfer{ID: ID}).Error
	})
}

func (f *Finance) Transfers(filter models.ListFilter) (*models.TransfersResponse, error) {
	cr := &models.TransfersResponse{}
	q := f.DB.Unscoped().Limit(filter.Limit).Offset(filter.Offset)
	// if IDs providen - return here and it has to be used as parent's ID, because we don't know transfers IDs other way
	if len(filter.IDs) > 0 {
		search := ""
		for i, step := range filter.IDs {
			search += fmt.Sprintf("parent_id = %d", step)
			if i < (len(filter.IDs) - 1) {
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
		q.Where("category LIKE ?", "%"+filter.Query+"%")
	}
	if filter.ParentID != 0 {
		q.Where("parent_id = ?", filter.ParentID)
	}
	if filter.From != 0 {
		q.Where("from = ?", filter.From)
	}
	if filter.To != 0 {
		q.Where("to = ?", filter.To)
	}
	if filter.Wallet != 0 {
		q.Where(f.DB.Where("`from` = ?", filter.Wallet).Or("`to` = ?", filter.Wallet))
	}
	AddDateCondition(filter, q)
	//TODO: check if it gives all uncompleted at first place
	q.Order("completed asc").Order("completed_at desc").Order("created_at desc")
	if result := q.Find(&cr.Transfers).Count(&cr.Total); result.Error != nil {
		return nil, result.Error
	}
	return cr, nil
}

func (f *Finance) Categories() ([]string, error) {
	categories := make([]string, 0)
	err := f.DB.Debug().Raw("SELECT DISTINCT(category) FROM transfers WHERE deleted_at IS NULL ORDER BY category asc").Scan(&categories).Error
	return categories, err
}

type CatTotal struct {
	Category string
	Total    int
}

type CategorisedCashflow struct {
	Incomes  []CatTotal
	Expenses []CatTotal
}

func (f *Finance) SumByCategory(filter models.ListFilter) (*CategorisedCashflow, error) {
	incomes := make([]CatTotal, 0)
	expenses := make([]CatTotal, 0)
	q := f.DB.Model(&models.Transfer{})
	AddDateCondition(filter, q)
	if err := q.Select("category, sum(amount) as total").Where("`from` IS NULL").Group("category").Find(&incomes).Error; err != nil {
		return nil, fmt.Errorf("Can't get incomes error: %s", err.Error())
	}
	q2 := f.DB.Model(&models.Transfer{})
	AddDateCondition(filter, q2)
	if err := q2.Select("category, sum(amount) as total").Where("`to` IS NULL").Group("category").Find(&expenses).Error; err != nil {
		return nil, fmt.Errorf("Can't get expenses error: %s", err.Error())
	}
	return &CategorisedCashflow{Incomes: incomes, Expenses: expenses}, nil
}

func AddDateCondition(filter models.ListFilter, q *gorm.DB) {
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
	//then we have to return datet or null
	if !filter.Completed && dateSearh != "" {
		dateSearh = fmt.Sprintf("((%s) OR completed_at IS NULL)", dateSearh)
	}
	q.Where(dateSearh)
}
