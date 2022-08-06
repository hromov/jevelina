package contacts

import (
	"context"
)

type Repository interface {
	ByID(context.Context, uint64) (Contact, error)
	ByPhone(context.Context, string) (Contact, error)
	Contacts(context.Context, Filter) (ContactsResponse, error)
	DeleteContact(context.Context, uint64) error
	UpdateContact(context.Context, ContactRequest) error
	CreateContact(context.Context, ContactRequest) (Contact, error)
}

//go:generate mockery --name Service --filename ContactsService.go --structname ContactsService --output ../../mocks
type Service interface {
	Get(context.Context, uint64) (Contact, error)
	GetByPhone(context.Context, string) (Contact, error)
	List(context.Context, Filter) (ContactsResponse, error)
	Delete(context.Context, uint64) error
	Update(context.Context, ContactRequest) error
	Create(context.Context, ContactRequest) (Contact, error)
	CreateOrGet(context.Context, ContactRequest) (Contact, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Get(ctx context.Context, id uint64) (Contact, error) {
	return s.r.ByID(ctx, id)
}

func (s *service) List(ctx context.Context, f Filter) (ContactsResponse, error) {
	return s.r.Contacts(ctx, f)
}

func (s *service) GetByPhone(ctx context.Context, phone string) (Contact, error) {
	return s.r.ByPhone(ctx, phone)
}

func (s *service) Update(ctx context.Context, contact ContactRequest) error {
	return s.r.UpdateContact(ctx, contact)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.r.DeleteContact(ctx, id)
}

func (s *service) Create(ctx context.Context, contact ContactRequest) (Contact, error) {
	return s.r.CreateContact(ctx, contact)
}

func (s *service) CreateOrGet(ctx context.Context, c ContactRequest) (Contact, error) {
	contact, err := s.GetByPhone(ctx, c.Phone)
	if err != nil || contact.ID == 0 {
		return s.Create(ctx, c)
	}
	return contact, nil
}
