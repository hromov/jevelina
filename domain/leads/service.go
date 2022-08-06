package leads

import (
	"context"
	"time"
)

type Repository interface {
	GetLead(context.Context, uint64) (Lead, error)
	GetLeads(context.Context, Filter) (LeadsResponse, error)
	CreateLead(context.Context, LeadRequest) (Lead, error)
	UpdateLead(context.Context, Lead) error
	DeleteLead(context.Context, uint64) error
}

//go:generate mockery --name Service --filename LeadsService.go --structname LeadsService --output ../../mocks
type Service interface {
	Get(context.Context, uint64) (Lead, error)
	List(context.Context, Filter) (LeadsResponse, error)
	Create(context.Context, LeadRequest) (Lead, error)
	Update(context.Context, LeadRequest) error
	Delete(context.Context, uint64) error
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

func (s *service) Update(ctx context.Context, lead Lead) error {
	if !lead.Step.Active && lead.ClosedAt.IsZero() {
		lead.ClosedAt = time.Now()
	}
	if lead.Step.Active && !lead.ClosedAt.IsZero() {
		lead.ClosedAt = time.Time{}
	}
	return s.r.UpdateLead(ctx, Lead)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.r.DeleteLead(ctx, id)
}

func (s *service) Create(ctx context.Context, Lead LeadRequest) (Lead, error) {
	return s.r.CreateLead(ctx, Lead)
}
