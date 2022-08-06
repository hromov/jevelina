package leads

import (
	"context"
	"fmt"
	"log"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/storage/mysql"
	"github.com/hromov/jevelina/storage/mysql/dao/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Leads struct {
	*gorm.DB
}

func (l *Leads) List(filter models.ListFilter) (*models.LeadsResponse, error) {
	// log.Println(limit, offset, query, query == "")
	cr := &models.LeadsResponse{}
	//How to make joins work?.Joins("Contacts")
	q := l.DB.Preload(clause.Associations).Limit(filter.Limit).Offset(filter.Offset)

	// if IDs providen - return here
	if len(filter.IDs) > 0 {
		if err := q.Find(&cr.Leads, filter.IDs).Count(&cr.Total).Error; err != nil {
			return nil, err
		}
		return cr, nil
	}
	if len(filter.Steps) > 0 {
		search := ""
		for i, step := range filter.Steps {
			search += fmt.Sprintf("step_id = %d", step)
			if i < (len(filter.Steps) - 1) {
				search += " OR "
			}
		}
		q = q.Where(search)
	}

	if filter.Query != "" {
		q = q.Where("name LIKE ?", "%"+filter.Query+"%")
	}
	if filter.ContactID != 0 {
		q = q.Where("contact_id = ?", filter.ContactID)
	}
	if filter.ResponsibleID != 0 {
		q = q.Where("responsible_id = ?", filter.ResponsibleID)
	}
	if filter.Active {
		q = q.Where("closed_at IS NULL")
	}
	if filter.StepID != 0 {
		q = q.Where("step_id = ?", filter.StepID)
	}
	dateSearh := ""
	if !filter.MinDate.IsZero() {
		dateSearh += fmt.Sprintf("closed_at >= '%s'", filter.MinDate)
	}
	if !filter.MaxDate.IsZero() {
		if dateSearh != "" {
			dateSearh += " AND "
		}
		dateSearh += fmt.Sprintf("closed_at < '%s'", filter.MaxDate)
	}
	//then we have to return datet or null
	if !filter.Completed && dateSearh != "" {
		dateSearh = fmt.Sprintf("((%s) OR closed_at IS NULL)", dateSearh)
	}
	q = q.Where(dateSearh).Order("created_at desc")

	if filter.TagID != 0 {
		IDs := []uint{}
		l.DB.Raw("select lead_id from leads_tags WHERE tag_id = ?", filter.TagID).Scan(&IDs)
		q.Find(&cr.Leads, IDs)
	} else {
		q.Find(&cr.Leads)
	}
	if result := q.Count(&cr.Total); result.Error != nil {
		return nil, result.Error
	}
	return cr, nil
}

// TODO: move to some domain
func (l *Leads) CreateLead(ctx context.Context, lr leads.LeadRequest, contact contacts.Contact) (leads.Lead, error) {
	lead := models.LeadFromRequest(lr, contact)

	if step, _ := mysql.Misc().DefaultStep(); step.ID != 0 {
		lead.StepID = &step.ID
	}
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

	if err := l.DB.WithContext(ctx).Omit(clause.Associations).Create(&lead).Error; err != nil {
		return leads.Lead{}, err
	}

	return l.GetLead(ctx, lead.ID)
}

func (l *Leads) GetLead(ctx context.Context, id uint64) (leads.Lead, error) {
	var lead models.Lead
	if err := l.DB.Unscoped().Preload(clause.Associations).First(&lead, id).Error; err != nil {
		return leads.Lead{}, err
	}

	return lead.ToDomain(), nil
}

func (l *Leads) DeleteLead(ctx context.Context, id uint64) error {
	if err := l.DB.WithContext(ctx).Delete(&models.Lead{ID: id}).Error; err != nil {
		return err
	}
	if err := misc.DeleteTaskByParent(l.DB, id); err != nil {
		log.Printf("Error while deliting tasks for lead: %s", err.Error())
	}
	return nil
}

func (l *Leads) SaveLead(ctx context.Context, lead leads.Lead) (leads.Lead, error) {
	dbLead := models.LeadFromDomain(lead)
	q := l.DB.WithContext(ctx).Omit(clause.Associations)
	if lead.ID == 0 {
		if err := q.Create(&dbLead).Error; err != nil {
			log.Printf("Can't create lead. Error: %s", err.Error())
			return leads.Lead{}, err
		}
	} else {
		if err := q.Save(&dbLead).Error; err != nil {
			log.Printf("Can't update lead with ID = %d. Error: %s", lead.ID, err.Error())
			return leads.Lead{}, err
		}
	}
	return l.GetLead(ctx, dbLead.ID)
}
