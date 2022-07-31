package users

import "context"

type Repository interface {
	Users(context.Context) ([]User, error)
}

type Service interface {
	List(context.Context) ([]User, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) List(ctx context.Context) ([]User, error) {
	return s.r.Users(ctx)
}
