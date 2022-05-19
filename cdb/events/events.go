package events

import (
	"fmt"
	"time"

	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
)

type EventService struct {
	DB *gorm.DB
}

func (es *EventService) Save(newEvent models.NewEvent) error {
	event := &models.Event{
		CreatedAt:       time.Now(),
		UserID:          newEvent.UserID,
		ParentID:        newEvent.ParentID,
		EventParentType: newEvent.EventParentType,
		Description:     fmt.Sprintf("[%s] %s", newEvent.EventType.String(), newEvent.Message),
	}
	return es.DB.Create(event).Error
}

func (es *EventService) List(filter models.EventFilter) (*models.EventsResponse, error) {
	er := &models.EventsResponse{}
	q := es.DB
	if filter.EventParentType != 0 {
		q = q.Where("event_parent_type = ?", filter.EventParentType)
	}

	if filter.UserID != 0 {
		q = q.Where("user_id = ?", filter.UserID)
	}

	if filter.ParentID != 0 {
		q = q.Where("parent_id = ?", filter.ParentID)
	}

	q.Order("created_at desc").Limit(filter.Limit).Offset(filter.Offset)

	if result := q.Find(&er.Events).Count(&er.Total); result.Error != nil {
		return nil, result.Error
	}
	return er, nil
}
