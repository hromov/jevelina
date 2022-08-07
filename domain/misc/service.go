package misc

import (
	"context"
)

type Repository interface {
	GetProduct(context.Context, uint32) (Product, error)
	GetProductByName(context.Context, string) (Product, error)
	ListProducts(context.Context) ([]Product, error)
	CreateProduct(context.Context, Product) (Product, error)
	UpdateProduct(context.Context, Product) error
	DeleteProduct(context.Context, uint32) error
}

//go:generate mockery --name Service --filename MiscService.go --structname MiscService --output ../../mocks
type Service interface {
	GetProduct(context.Context, uint32) (Product, error)
	GetProductByName(context.Context, string) (Product, error)
	ListProducts(context.Context) ([]Product, error)
	CreateProduct(context.Context, Product) (Product, error)
	UpdateProduct(context.Context, Product) error
	DeleteProduct(context.Context, uint32) error
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) GetProduct(ctx context.Context, id uint32) (Product, error) {
	return s.r.GetProduct(ctx, id)
}
func (s *service) GetProductByName(ctx context.Context, name string) (Product, error) {
	return s.r.GetProductByName(ctx, name)
}
func (s *service) ListProducts(ctx context.Context) ([]Product, error) {
	return s.r.ListProducts(ctx)
}
func (s *service) CreateProduct(ctx context.Context, p Product) (Product, error) {
	return s.r.CreateProduct(ctx, p)
}
func (s *service) UpdateProduct(ctx context.Context, p Product) error {
	return s.r.UpdateProduct(ctx, p)
}
func (s *service) DeleteProduct(ctx context.Context, id uint32) error {
	return s.r.DeleteProduct(ctx, id)
}
