package leads

import (
	"fmt"

	"github.com/hromov/jevelina/cdb/models"
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
		q = q.Find(&cr.Leads, IDs)
	} else {
		q = q.Find(&cr.Leads)
	}
	if result := q.Count(&cr.Total); result.Error != nil {
		return nil, result.Error
	}
	return cr, nil
}

func (l *Leads) ByID(ID uint64) (*models.Lead, error) {
	// log.Println(limit, offset, query, query == "")
	var lead models.Lead

	if result := l.DB.Preload(clause.Associations).First(&lead, ID); result.Error != nil {
		return nil, result.Error
	}
	return &lead, nil
}
