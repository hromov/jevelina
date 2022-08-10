package events

import (
	"context"
	"fmt"
)

type Repository interface {
	SaveEvent(ctx context.Context, event Event) error
	GetEvents(ctx context.Context, filter EventFilter) (EventsResponse, error)
}

//go:generate mockery --name Service --filename EventsService.go --structname EventsService --output ../../mocks
type Service interface {
	Save(ctx context.Context, event NewEvent) error
	List(ctx context.Context, filter EventFilter) (EventsResponse, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Save(ctx context.Context, e NewEvent) error {
	event := Event{
		UserID:          e.UserID,
		ParentID:        e.ParentID,
		EventParentType: e.EventParentType,
		Description:     fmt.Sprintf("[%s] %s", e.EventType.String(), e.Message),
	}
	return s.r.SaveEvent(ctx, event)
}

func (s *service) List(ctx context.Context, filter EventFilter) (EventsResponse, error) {
	return s.r.GetEvents(ctx, filter)
}
