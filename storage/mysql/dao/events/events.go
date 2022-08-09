package events

import (
	"context"
	"log"

	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"github.com/hromov/jevelina/utils/events"
	"gorm.io/gorm"
)

type Events struct {
	db *gorm.DB
}

func NewEvents(db *gorm.DB) *Events {
	if err := db.AutoMigrate(&models.Event{}); err != nil {
		log.Println("Can't auto migrate events error: ", err.Error())
	}
	return &Events{db}
}

func (e *Events) SaveEvent(ctx context.Context, newEvent events.Event) error {
	event := models.Event(newEvent)
	return e.db.WithContext(ctx).Create(&event).Error
}

func (e *Events) GetEvents(ctx context.Context, filter events.EventFilter) (events.EventsResponse, error) {
	er := events.EventsResponse{}
	q := e.db.WithContext(ctx)
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
	list := []models.Event{}
	if err := q.Find(&list).Count(&er.Total).Error; err != nil {
		return events.EventsResponse{}, err
	}
	er.Events = make([]events.Event, len(list))
	for i, event := range list {
		er.Events[i] = event.ToDomain()
	}
	return er, nil
}
