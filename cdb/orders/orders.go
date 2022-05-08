package orders

import (
	"strings"
	"time"

	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

func CreateOrGetContact(c *models.CreateLeadReq, user *models.User) (*models.Contact, error) {
	cdb.Contacts().List(models.ListFilter{Query: c.ClientPhone})
	var contact *models.Contact
	// var err error
	if c.ClientPhone != "" && len(c.ClientPhone) > 5 {
		contact, _ = cdb.Contacts().ByPhone(c.ClientPhone)
		// if err != nil {
		// 	return nil, err
		// }
	}
	if contact == nil {
		contact = &models.Contact{
			Name:  c.ClientName,
			Phone: c.ClientPhone,
			Email: c.ClientEmail,
		}
		//TODO: move to generics in 1.19, if possible
		item := contact
		item.Analytics.CID = c.CID
		//if no UID was provided - use phone
		if c.UID == "" {
			item.Analytics.UID = contact.Phone
		} else {
			item.Analytics.UID = c.UID
		}
		item.Analytics.TID = c.TID
		item.Analytics.UtmID = c.UtmID
		item.Analytics.UtmSource = c.UtmSource
		item.Analytics.UtmMedium = c.UtmMedium
		item.Analytics.UtmCampaign = c.UtmCampaign
		item.Analytics.Domain = c.Domain

		if c.Source != "" {
			if source, _ := cdb.Misc().SourceByName(c.Source); source != nil {
				contact.SourceID = &source.ID
			}
		}

		contact.ResponsibleID = &user.ID
	} else {
		//check for updated fields
		if c.ClientName != "" && strings.Compare(contact.Name, c.ClientName) != 0 {
			if contact.Name == "" {
				contact.Name = c.ClientName
			} else if contact.SecondName == "" {
				contact.SecondName = c.ClientName
			}
		}
		if c.ClientEmail != "" && strings.Compare(contact.Email, c.ClientEmail) != 0 {
			if contact.Email == "" {
				contact.Email = c.ClientEmail
			} else if contact.SecondEmail == "" {
				contact.SecondEmail = c.ClientEmail
			}
		}
	}

	if contact.ID == 0 {
		if err := cdb.Contacts().Create(contact).Error; err != nil {
			return nil, err
		}
	} else {
		if err := cdb.Contacts().Save(contact).Error; err != nil {
			return nil, err
		}
	}
	return contact, nil
}

func CreateLead(c *models.CreateLeadReq, contact *models.Contact) (*models.Lead, error) {
	lead := &models.Lead{
		Name:          c.Name,
		Budget:        uint32(c.Price),
		ResponsibleID: contact.ResponsibleID,
		ContactID:     &contact.ID,
	}
	if step, _ := cdb.Misc().DefaultStep(); step != nil {
		lead.StepID = &step.ID
	}
	item := lead
	item.Analytics.CID = c.CID
	item.Analytics.UID = c.UID
	item.Analytics.TID = c.TID
	item.Analytics.UtmID = c.UtmID
	item.Analytics.UtmSource = c.UtmSource
	item.Analytics.UtmMedium = c.UtmMedium
	item.Analytics.UtmCampaign = c.UtmCampaign
	item.Analytics.Domain = c.Domain

	if c.Source != "" {
		if source, _ := cdb.Misc().SourceByName(c.Source); source != nil {
			lead.SourceID = &source.ID
		}
	}
	if c.Product != "" {
		if product, _ := cdb.Misc().ProductByName(c.Product); product != nil {
			lead.ProductID = &product.ID
		}
	}
	if c.Manufacturer != "" {
		if manuf, _ := cdb.Misc().ManufacturerByName(c.Manufacturer); manuf != nil {
			lead.ManufacturerID = &manuf.ID
		}
	}
	if err := cdb.Leads().Create(lead).Error; err != nil {
		return nil, err
	}
	return lead, nil
}

func CreateTask(c *models.CreateLeadReq, lead *models.Lead) error {
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
	if err := cdb.Misc().Create(task).Error; err != nil {
		return err
	}
	return nil
}
