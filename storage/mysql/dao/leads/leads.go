package leads

import (
	"context"
	"fmt"
	"log"

	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Leads struct {
	*gorm.DB
}

func (l *Leads) GetLeads(ctx context.Context, filter leads.Filter) (leads.LeadsResponse, error) {
	cr := &models.LeadsResponse{}
	q := l.DB.WithContext(ctx).Preload(clause.Associations).Limit(filter.Limit).Offset(filter.Offset)

	// if IDs providen - return here
	if len(filter.IDs) > 0 {
		if err := q.Find(&cr.Leads, filter.IDs).Count(&cr.Total).Error; err != nil {
			return leads.LeadsResponse{}, err
		}
	} else {
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

		// old version delete if works without
		// if filter.TagID != 0 {
		// 	IDs := []uint{}
		// 	l.DB.Raw("select lead_id from leads_tags WHERE tag_id = ?", filter.TagID).Scan(&IDs)
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

func (l *Leads) CreateLead(ctx context.Context, lead leads.LeadData) (leads.Lead, error) {
	if step, _ := l.DefaultStep(ctx); step.ID != 0 {
		lead.StepID = step.ID
	}
	dbLead := models.LeadFromDomain(lead)
	if err := l.DB.WithContext(ctx).Omit(clause.Associations).Create(&dbLead).Error; err != nil {
		return leads.Lead{}, err
	}
	return l.GetLead(ctx, dbLead.ID)
}

func (l *Leads) UpdateLead(ctx context.Context, lead leads.LeadData) error {
	dbLead := models.LeadFromDomain(lead)
	return l.DB.WithContext(ctx).Omit(clause.Associations).Where("id", lead.ID).Updates(&dbLead).Error
}

func (l *Leads) Steps(ctx context.Context) ([]models.Step, error) {
	var items []models.Step
	if result := l.DB.WithContext(ctx).Order("`order`").Find(&items); result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (l *Leads) Step(ctx context.Context, id uint8) (*models.Step, error) {
	var item models.Step
	if result := l.DB.WithContext(ctx).First(&item, id); result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (l *Leads) DefaultStep(ctx context.Context) (models.Step, error) {
	var item models.Step
	if err := l.DB.WithContext(ctx).Where("`order` = 0").First(&item).Error; err != nil {
		return models.Step{}, err
	}
	return item, nil
}

// TODO: move or return Task (make full crud)
func (l *Leads) CreateTask(ctx context.Context, t leads.TaskData) error {
	dbTask := models.TaskFromTaskData(t)
	if err := l.DB.WithContext(ctx).Omit(clause.Associations).Create(&dbTask).Error; err != nil {
		return err
	}
	return nil
}
