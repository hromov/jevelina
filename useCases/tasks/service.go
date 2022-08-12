package tasks

import "context"

type Repository interface {
	GetTask(context.Context, uint64) (Task, error)
	GetTasks(context.Context, Filter) (TasksResponse, error)
	CreateTask(context.Context, TaskData) (Task, error)
	UpdateTask(context.Context, TaskData) error
	DeleteTask(context.Context, uint64) error
	DeleteTaskByParent(context.Context, uint64) error
}

//go:generate mockery --name Service --filename TasksService.go --structname TasksService --output ../../mocks
type Service interface {
	Get(context.Context, uint64) (Task, error)
	List(context.Context, Filter) (TasksResponse, error)
	Create(context.Context, TaskData) (Task, error)
	Update(context.Context, TaskData) error
	Delete(context.Context, uint64) error
	DeleteByParent(context.Context, uint64) error
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Get(ctx context.Context, id uint64) (Task, error) {
	return s.r.GetTask(ctx, id)
}
func (s *service) List(ctx context.Context, f Filter) (TasksResponse, error) {
	return s.r.GetTasks(ctx, f)
}
func (s *service) Create(ctx context.Context, t TaskData) (Task, error) {
	return s.r.CreateTask(ctx, t)
}
func (s *service) Update(ctx context.Context, t TaskData) error {
	return s.r.UpdateTask(ctx, t)
}
func (s *service) Delete(ctx context.Context, id uint64) error {
	return s.r.DeleteTask(ctx, id)
}
func (s *service) DeleteByParent(ctx context.Context, parentID uint64) error {
	return s.r.DeleteTaskByParent(ctx, parentID)
}
