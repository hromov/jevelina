package orders

import (
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/storage/mysql"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
)

func CreateLead(c models.CreateLeadReq, contact contacts.Contact) (models.Lead, error) {
	lead := models.Lead{
		Name:          c.Name,
		Budget:        uint32(c.Price),
		ResponsibleID: &contact.Responsible.ID,
		ContactID:     &contact.ID,
	}
	if step, _ := mysql.Misc().DefaultStep(); step != nil {
		lead.StepID = &step.ID
	}
	item := lead
	item.Analytics.CID = c.CID
	uid := ""
	if c.UID == "" {
		uid = contact.Analytics.UID
	}
	item.Analytics.UID = uid
	item.Analytics.TID = c.TID
	item.Analytics.UtmID = c.UtmID
	item.Analytics.UtmSource = c.UtmSource
	item.Analytics.UtmMedium = c.UtmMedium
	item.Analytics.UtmCampaign = c.UtmCampaign
	item.Analytics.Domain = c.Domain

	if c.Source != "" {
		if source, _ := mysql.Misc().SourceByName(c.Source); source != nil {
			lead.SourceID = &source.ID
		}
	}
	if c.Product != "" {
		if product, _ := mysql.Misc().ProductByName(c.Product); product != nil {
			lead.ProductID = &product.ID
		}
	}
	if c.Manufacturer != "" {
		if manuf, _ := mysql.Misc().ManufacturerByName(c.Manufacturer); manuf != nil {
			lead.ManufacturerID = &manuf.ID
		}
	}
	if err := mysql.Leads().Create(lead).Error; err != nil {
		return models.Lead{}, err
	}
	return lead, nil
}

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
