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
	GetManufacturer(context.Context, uint32) (Manufacturer, error)
	GetManufacturerByName(context.Context, string) (Manufacturer, error)
	ListManufacturers(context.Context) ([]Manufacturer, error)
	CreateManufacturer(context.Context, Manufacturer) (Manufacturer, error)
	UpdateManufacturer(context.Context, Manufacturer) error
	DeleteManufacturer(context.Context, uint32) error
	GetSource(context.Context, uint32) (Source, error)
	GetSourceByName(context.Context, string) (Source, error)
	ListSources(context.Context) ([]Source, error)
	CreateSource(context.Context, Source) (Source, error)
	UpdateSource(context.Context, Source) error
	DeleteSource(context.Context, uint32) error
}

//go:generate mockery --name Service --filename MiscService.go --structname MiscService --output ../../mocks
type Service interface {
	GetProduct(context.Context, uint32) (Product, error)
	GetProductByName(context.Context, string) (Product, error)
	ListProducts(context.Context) ([]Product, error)
	CreateProduct(context.Context, Product) (Product, error)
	UpdateProduct(context.Context, Product) error
	DeleteProduct(context.Context, uint32) error
	GetManufacturer(context.Context, uint32) (Manufacturer, error)
	GetManufacturerByName(context.Context, string) (Manufacturer, error)
	ListManufacturers(context.Context) ([]Manufacturer, error)
	CreateManufacturer(context.Context, Manufacturer) (Manufacturer, error)
	UpdateManufacturer(context.Context, Manufacturer) error
	DeleteManufacturer(context.Context, uint32) error
	GetSource(context.Context, uint32) (Source, error)
	GetSourceByName(context.Context, string) (Source, error)
	ListSources(context.Context) ([]Source, error)
	CreateSource(context.Context, Source) (Source, error)
	UpdateSource(context.Context, Source) error
	DeleteSource(context.Context, uint32) error
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

func (s *service) GetManufacturer(ctx context.Context, id uint32) (Manufacturer, error) {
	return s.r.GetManufacturer(ctx, id)
}
func (s *service) GetManufacturerByName(ctx context.Context, name string) (Manufacturer, error) {
	return s.r.GetManufacturerByName(ctx, name)
}
func (s *service) ListManufacturers(ctx context.Context) ([]Manufacturer, error) {
	return s.r.ListManufacturers(ctx)
}
func (s *service) CreateManufacturer(ctx context.Context, p Manufacturer) (Manufacturer, error) {
	return s.r.CreateManufacturer(ctx, p)
}
func (s *service) UpdateManufacturer(ctx context.Context, p Manufacturer) error {
	return s.r.UpdateManufacturer(ctx, p)
}
func (s *service) DeleteManufacturer(ctx context.Context, id uint32) error {
	return s.r.DeleteManufacturer(ctx, id)
}

func (s *service) GetSource(ctx context.Context, id uint32) (Source, error) {
	return s.r.GetSource(ctx, id)
}
func (s *service) GetSourceByName(ctx context.Context, name string) (Source, error) {
	return s.r.GetSourceByName(ctx, name)
}
func (s *service) ListSources(ctx context.Context) ([]Source, error) {
	return s.r.ListSources(ctx)
}
func (s *service) CreateSource(ctx context.Context, p Source) (Source, error) {
	return s.r.CreateSource(ctx, p)
}
func (s *service) UpdateSource(ctx context.Context, p Source) error {
	return s.r.UpdateSource(ctx, p)
}
func (s *service) DeleteSource(ctx context.Context, id uint32) error {
	return s.r.DeleteSource(ctx, id)
}
