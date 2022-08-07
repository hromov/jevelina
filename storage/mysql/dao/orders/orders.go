package orders

import (
	"context"
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/storage/mysql"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm/clause"
)

func CreateTask(c models.CreateLeadReq, lead models.Lead) error {
	task := new(models.Task)
	if c.Description != "" {
		task.Description = c.Description
	} else {
		task.Description = "Call me!"
	}
	t := time.Now()
	task.DeadLine = &t
	task.ParentID = lead.ID
	task.ResponsibleID = lead.ResponsibleID
	if err := mysql.Misc().Create(task).Error; err != nil {
		return err
	}
	return nil
}

// TODO: move to some domain
func CreateOrder(ctx context.Context, lr leads.LeadRequest, contact contacts.Contact) (leads.Lead, error) {
	lead := models.LeadFromRequest(lr, contact)

	// if step, _ := mysql.  ..DefaultStep(context.TODO()); step.ID != 0 {
	// 	lead.StepID = &step.ID
	// }
	if lr.Source != "" {
		if source, _ := mysql.Misc().SourceByName(lr.Source); source != nil {
			lead.SourceID = &source.ID
		}
	}
	if lr.Product != "" {
		if product, _ := mysql.Misc().ProductByName(lr.Product); product != nil {
			lead.ProductID = &product.ID
		}
	}
	if lr.Manufacturer != "" {
		if manuf, _ := mysql.Misc().ManufacturerByName(lr.Manufacturer); manuf != nil {
			lead.ManufacturerID = &manuf.ID
		}
	}

	if err := mysql.Leads().DB.WithContext(ctx).Omit(clause.Associations).Create(&lead).Error; err != nil {
		return leads.Lead{}, err
	}

	return mysql.Leads().GetLead(ctx, lead.ID)
}
