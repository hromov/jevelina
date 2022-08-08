package users

import "context"

type Repository interface {
	GetUsers(context.Context) ([]User, error)
	User(context.Context, uint64) (User, error)
	UserByEmail(ctx context.Context, mail string) (User, error)
	UserExist(ctx context.Context, mail string) (bool, error)
	CreateUser(context.Context, ChangeUser) (User, error)
	UpdateUser(context.Context, ChangeUser) error
	DeleteUser(context.Context, uint64) error
	CreateRole(context.Context, Role) (Role, error)
	UpdateRole(context.Context, Role) error
	DeleteRole(ctx context.Context, id uint8) error
	Roles(context.Context) ([]Role, error)
	Role(context.Context, uint8) (Role, error)
}

//go:generate mockery --name Service --filename UsersService.go --structname UsersService --output ../../mocks
type Service interface {
	List(context.Context) ([]User, error)
	Get(context.Context, uint64) (User, error)
	GetByEmail(ctx context.Context, mail string) (User, error)
	UserExist(ctx context.Context, mail string) (bool, error)
	Create(context.Context, ChangeUser) (User, error)
	Update(context.Context, ChangeUser) error
	Delete(context.Context, uint64) error
	CreateRole(context.Context, Role) (Role, error)
	UpdateRole(context.Context, Role) error
	DeleteRole(ctx context.Context, id uint8) error
	ListRoles(context.Context) ([]Role, error)
	GetRole(context.Context, uint8) (Role, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) List(ctx context.Context) ([]User, error) {
	return s.r.GetUsers(ctx)
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

func (s *service) UpdateRole(ctx context.Context, role Role) error {
	return s.r.UpdateRole(ctx, role)
}

func (s *service) DeleteRole(ctx context.Context, id uint8) error {
	return s.r.DeleteRole(ctx, id)
}

func (s *service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.r.Roles(ctx)
}

func (s *service) GetRole(ctx context.Context, id uint8) (Role, error) {
	return s.r.Role(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (User, error) {
	return s.r.UserByEmail(ctx, email)
}

func (s *service) UserExist(ctx context.Context, mail string) (bool, error) {
	return s.r.UserExist(ctx, mail)
}
