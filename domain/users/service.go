package users

import "context"

type Repository interface {
	Users(context.Context) ([]User, error)
	User(context.Context, uint64) (User, error)
	CreateUser(context.Context, ChangeUser) (User, error)
	UpdateUser(context.Context, ChangeUser) error
	DeleteUser(context.Context, uint64) error
	CreateRole(context.Context, Role) (Role, error)
}

//go:generate mockery --name Service --filename UsersService.go --structname UsersService --output ../../mocks
type Service interface {
	List(context.Context) ([]User, error)
	Get(context.Context, uint64) (User, error)
	Create(context.Context, ChangeUser) (User, error)
	Update(context.Context, ChangeUser) error
	Delete(context.Context, uint64) error
	CreateRole(context.Context, Role) (Role, error)
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

func (s *service) Get(ctx context.Context, id uint64) (User, error) {
	return s.r.User(ctx, id)
}

func (s *service) Create(ctx context.Context, newUser ChangeUser) (User, error) {
	return s.r.CreateUser(ctx, newUser)
}

func (s *service) Update(ctx context.Context, newUser ChangeUser) error {
	return s.r.UpdateUser(ctx, newUser)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.r.DeleteUser(ctx, id)
}

func (s *service) CreateRole(ctx context.Context, role Role) (Role, error) {
	return s.r.CreateRole(ctx, role)
}
