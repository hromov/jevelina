package orders

import (
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/misc"
)

type Order struct {
	Name        string
	Price       int
	Description string

	ClientName  string
	ClientEmail string
	ClientPhone string

	ProductID      uint32
	ManufacturerID uint16
	SourceID       uint8

	CID string
	UID string
	TID string

	UtmID       string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string

	Domain string
}

func (o *Order) ToLeadData(contactID, userID uint64) leads.LeadData {
	return leads.LeadData{
		Name:   o.Name,
		Budget: uint32(o.Price),

		ContactID:      contactID,
		ResponsibleID:  userID,
		CreatedID:      userID,
		ProductID:      o.ProductID,
		ManufacturerID: o.ManufacturerID,
		SourceID:       o.SourceID,
		Analytics: misc.Analytics{
			CID:         o.CID,
			TID:         o.TID,
			UtmID:       o.UtmID,
			UtmSource:   o.UtmSource,
			UtmMedium:   o.UtmMedium,
			UtmCampaign: o.UtmCampaign,
			Domain:      o.Domain,
		},
	}

}

func (o *Order) ToContactRequest(userID uint64) contacts.ContactRequest {
	cr := contacts.ContactRequest{
		Name:          o.ClientName,
		Email:         o.ClientEmail,
		Phone:         o.ClientPhone,
		ResponsibleID: userID,
		Analytics: misc.Analytics{
			CID:         o.CID,
			TID:         o.TID,
			UtmID:       o.UtmID,
			UtmSource:   o.UtmSource,
			UtmMedium:   o.UtmMedium,
			UtmCampaign: o.UtmCampaign,
			Domain:      o.Domain,
		},
	}
	if o.UID == "" {
		cr.Analytics.UID = o.ClientPhone
	} else {
		cr.Analytics.UID = o.UID
	}

	return cr
}

func (o *Order) ToTaskData(leadID, userID uint64) leads.TaskData {
	task := leads.TaskData{
		ParentID:      leadID,
		ResponsibleID: userID,
		CreatedID:     userID,
	}
	if o.Description != "" {
		task.Description = o.Description
	} else {
		task.Description = "Call me!"
	}
	task.DeadLine = time.Now()
	return task
}
