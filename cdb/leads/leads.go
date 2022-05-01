package leads

import (
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
	q := l.DB.Preload(clause.Associations).Order("created_at desc").Limit(filter.Limit).Offset(filter.Offset)
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
