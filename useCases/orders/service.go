package orders

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/users"
)

//go:generate mockery --name Service --filename OrdersService.go --structname OrdersService --output ../../mocks
type Service interface {
	Create(context.Context, Order) error
	CreateForUser(context.Context, Order, users.User) error
}

type service struct {
	cs contacts.Service
	ls leads.Service
	us users.Service
}

func NewService(cs contacts.Service, ls leads.Service, us users.Service) *service {
	return &service{cs, ls, us}
}

func (s *service) Create(ctx context.Context, order Order) error {
	user, err := s.GetRandomUser(ctx)
	if err != nil {
		return err
	}

	return s.CreateForUser(ctx, order, user)
}

func (s *service) CreateForUser(ctx context.Context, order Order, user users.User) error {
	contact, err := s.cs.CreateOrGet(ctx, order.ToContactRequest(user.ID))
	if err != nil {
		return err
	}

	lead, err := s.ls.Create(ctx, order.ToLeadData(contact.ID, user.ID))
	if err != nil {
		return err
	}

	task := order.ToTaskData(lead.ID, user.ID)
	err = s.ls.CreateTask(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetRandomUser(ctx context.Context) (users.User, error) {
	userList, err := s.us.List(ctx)
	if err != nil {
		return users.User{}, nil
	}

	filtered := []users.User{}
	for _, u := range userList {
		if u.Distribution > 0.0 {
			filtered = append(filtered, u)
		}
	}

	if len(filtered) == 0 {
		return users.User{}, errors.New("no good users found")
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(filtered))
	return filtered[r], nil
}
