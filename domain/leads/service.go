package leads

import (
	"context"
	"time"
)

type Repository interface {
	GetLead(context.Context, uint64) (Lead, error)
	GetLeads(context.Context, Filter) (LeadsResponse, error)
	CreateLead(context.Context, LeadData) (Lead, error)
	UpdateLead(context.Context, LeadData) error
	DeleteLead(context.Context, uint64) error
	CreateTask(context.Context, TaskData) error
	GetStep(context.Context, uint8) (Step, error)
	GetSteps(context.Context) ([]Step, error)
	CreateStep(context.Context, Step) (Step, error)
	UpdateStep(context.Context, Step) error
	DeleteStep(context.Context, uint8) error
}

//go:generate mockery --name Service --filename LeadsService.go --structname LeadsService --output ../../mocks
type Service interface {
	Get(context.Context, uint64) (Lead, error)
	List(context.Context, Filter) (LeadsResponse, error)
	Create(context.Context, LeadData) (Lead, error)
	Update(context.Context, LeadData) error
	Delete(context.Context, uint64) error
	GetStep(context.Context, uint8) (Step, error)
	GetSteps(context.Context) ([]Step, error)
	CreateTask(context.Context, TaskData) error
	CreateStep(context.Context, Step) (Step, error)
	UpdateStep(context.Context, Step) error
	DeleteStep(context.Context, uint8) error
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Get(ctx context.Context, id uint64) (Lead, error) {
	return s.r.GetLead(ctx, id)
}

func (s *service) List(ctx context.Context, f Filter) (LeadsResponse, error) {
	return s.r.GetLeads(ctx, f)
}

func (s *service) Update(ctx context.Context, lead LeadData) error {
	if !lead.Step.Active && lead.ClosedAt.IsZero() {
		lead.ClosedAt = time.Now()
	}
	if lead.Step.Active && !lead.ClosedAt.IsZero() {
		lead.ClosedAt = time.Time{}
	}
	return s.r.UpdateLead(ctx, lead)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.r.DeleteLead(ctx, id)
}

func (s *service) Create(ctx context.Context, lead LeadData) (Lead, error) {
	return s.r.CreateLead(ctx, lead)
}

func (s *service) CreateTask(ctx context.Context, t TaskData) error {
	return s.r.CreateTask(ctx, t)
}
func (s *service) CreateStep(ctx context.Context, step Step) (Step, error) {
	return s.r.CreateStep(ctx, step)
}
func (s *service) UpdateStep(ctx context.Context, step Step) error {
	return s.r.UpdateStep(ctx, step)
}

func (s *service) GetStep(ctx context.Context, id uint8) (Step, error) {
	return s.r.GetStep(ctx, id)
}

func (s *service) DeleteStep(ctx context.Context, id uint8) error {
	return s.r.DeleteStep(ctx, id)
}

func (s *service) GetSteps(ctx context.Context) ([]Step, error) {
	return s.r.GetSteps(ctx)
}
