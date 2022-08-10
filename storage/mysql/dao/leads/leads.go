package leads

import (
	"context"
	"fmt"
	"log"

	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Leads struct {
	db *gorm.DB
}

func NewLeads(db *gorm.DB) *Leads {
	if err := db.AutoMigrate(&models.Lead{}); err != nil {
		log.Println("Can't megrate leads error: ", err.Error())
	}
	if err := db.AutoMigrate(&models.Step{}); err != nil {
		log.Println("Can't megrate leads error: ", err.Error())
	}
	return &Leads{db}
}

func (l *Leads) GetLeads(ctx context.Context, filter leads.Filter) (leads.LeadsResponse, error) {
	cr := &models.LeadsResponse{}
	q := l.db.WithContext(ctx).Preload(clause.Associations).Limit(filter.Limit).Offset(filter.Offset)

	// if IDs providen - return here
	if len(filter.IDs) > 0 {
		if err := q.Find(&cr.Leads, filter.IDs).Count(&cr.Total).Error; err != nil {
			return leads.LeadsResponse{}, err
		}
	} else {
		if len(filter.Steps) > 0 {
			search := ""
			for i, step := range filter.Steps {
				//TODO: it has to be always number or find a solution to avoid fmt.Sprintf()
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
			//TODO: it has to be always Time or find a solution to avoid fmt.Sprintf()
			dateSearh += fmt.Sprintf("closed_at >= '%s'", filter.MinDate)
		}
		if !filter.MaxDate.IsZero() {
			if dateSearh != "" {
				dateSearh += " AND "
			}
			//TODO: it has to be always Time or find a solution to avoid fmt.Sprintf()
			dateSearh += fmt.Sprintf("closed_at < '%s'", filter.MaxDate)
		}
		//then we have to return datet or null
		if !filter.Completed && dateSearh != "" {
			dateSearh = fmt.Sprintf("((%s) OR closed_at IS NULL)", dateSearh)
		}
		q = q.Where(dateSearh).Order("created_at desc")

		// old version delete if works without
		// if filter.TagID != 0 {
		// 	IDs := []uint{}
		// 	l.db.Raw("select lead_id from leads_tags WHERE tag_id = ?", filter.TagID).Scan(&IDs)
		// 	q.Find(&cr.Leads, IDs)
		// } else {
		// 	q.Find(&cr.Leads)
		// }
		if err := q.Find(&cr.Leads).Count(&cr.Total).Error; err != nil {
			return leads.LeadsResponse{}, err
		}
	}

	resp := leads.LeadsResponse{
		Leads: make([]leads.Lead, len(cr.Leads)),
		Total: cr.Total,
	}
	for i, l := range cr.Leads {
		resp.Leads[i] = l.ToDomain()
	}
	return resp, nil
}

func (l *Leads) GetLead(ctx context.Context, id uint64) (leads.Lead, error) {
	var lead models.Lead
	if err := l.db.Unscoped().Preload(clause.Associations).First(&lead, id).Error; err != nil {
		return leads.Lead{}, err
	}

	return lead.ToDomain(), nil
}

func (l *Leads) DeleteLead(ctx context.Context, id uint64) error {
	if err := l.db.WithContext(ctx).Delete(&models.Lead{ID: id}).Error; err != nil {
		return err
	}
	// TODO: do with hooks ?
	// if err := misc.DeleteTaskByParent(l.db, id); err != nil {
	// 	log.Printf("Error while deliting tasks for lead: %s", err.Error())
	// }
	return nil
}

func (l *Leads) CreateLead(ctx context.Context, lead leads.LeadData) (leads.Lead, error) {
	if step, _ := l.DefaultStep(ctx); step.ID != 0 {
		lead.StepID = step.ID
	}
	dbLead := models.LeadFromDomain(lead)
	if err := l.db.WithContext(ctx).Omit(clause.Associations).Create(&dbLead).Error; err != nil {
		return leads.Lead{}, err
	}
	return l.GetLead(ctx, dbLead.ID)
}

func (l *Leads) UpdateLead(ctx context.Context, lead leads.LeadData) error {
	dbLead := models.LeadFromDomain(lead)
	return l.db.WithContext(ctx).Omit(clause.Associations).Where("id", lead.ID).Updates(&dbLead).Error
}

func (l *Leads) GetLeadsByDates(ctx context.Context, filter leads.Filter) (leads.LeadsResponse, error) {
	if filter.ByCreationDate {
		return l.GetLeadsByDates(ctx, filter)
	}

	cr := &models.LeadsResponse{}
	q := l.db.WithContext(ctx).Preload(clause.Associations).Limit(filter.Limit).Offset(filter.Offset)

	if len(filter.Steps) > 0 {
		search := ""
		for i, step := range filter.Steps {
			//TODO: it has to be always number or find a solution to avoid fmt.Sprintf()
			search += fmt.Sprintf("step_id = %d", step)
			if i < (len(filter.Steps) - 1) {
				search += " OR "
			}
		}
		q = q.Where(search)
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
		//TODO: it has to be always Time or find a solution to avoid fmt.Sprintf()
		dateSearh += fmt.Sprintf("created_at >= '%s'", filter.MinDate)
	}
	if !filter.MaxDate.IsZero() {
		if dateSearh != "" {
			dateSearh += " AND "
		}
		//TODO: it has to be always Time or find a solution to avoid fmt.Sprintf()
		dateSearh += fmt.Sprintf("created_at < '%s'", filter.MaxDate)
	}
	q = q.Where(dateSearh).Order("created_at desc")

	if err := q.Find(&cr.Leads).Count(&cr.Total).Error; err != nil {
		return leads.LeadsResponse{}, err
	}

	resp := leads.LeadsResponse{
		Leads: make([]leads.Lead, len(cr.Leads)),
		Total: cr.Total,
	}
	for i, l := range cr.Leads {
		resp.Leads[i] = l.ToDomain()
	}
	return resp, nil
}
